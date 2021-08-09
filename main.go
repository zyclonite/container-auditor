package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":9103", "The address to listen on for HTTP requests.")
var dockerHost = flag.String("docker-host", "unix:///var/run/docker.sock", "Docker host socket or tcp port.")
var (
  versionGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
    Name: "container_image_version_info",
    Help: "Image hash and name for all running containers.",
  },
	[]string{"container_id", "container_names", "image", "image_id"})
)

func recordMetrics() {
  go func() {
    for {
			ctx := context.Background()
			cli, err := client.NewClientWithOpts(client.WithHost(*dockerHost), client.WithAPIVersionNegotiation())
			if err != nil {
				log.Print(err)
			}

			containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: false})
			if err != nil {
				log.Print(err)
			}

			versionGauge.Reset()
			for _, container := range containers {
				//log.Print("ContainerID: ", container.ID, "\nNames: ", strings.Join(container.Names, ", "), "\nImage: ", container.Image, "\nImageID: ", container.ImageID)
				versionGauge.WithLabelValues(container.ID, strings.Join(container.Names, ", "), container.Image, container.ImageID).Set(1.0)
			}

      time.Sleep(10 * time.Second)
    }
  }()
}

func main() {
	flag.Parse()

	reg := prometheus.NewRegistry()
	reg.MustRegister(versionGauge)
	promhandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	recordMetrics()

	// sudo ./container-auditor -listen-address :9303 -docker-host unix:///var/run/podman/podman.sock
	r := mux.NewRouter()

	r.Path("/metrics").Handler(promhandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(*addr, r))
}
