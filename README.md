[![Docker Pulls](https://badgen.net/docker/pulls/zyclonite/container-auditor)](https://hub.docker.com/r/zyclonite/container-auditor)

## container-auditor

### build

`docker build -t zyclonite/container-auditor .`

### run

`docker run --name container-auditor -d -v /var/run/docker.sock:/var/run/docker.sock -p 9103:9103 zyclonite/container-auditor`
