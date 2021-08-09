FROM golang:1.16-alpine as builder

COPY . /go/src

RUN cd src \
  && go mod download \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -ldflags "-s -w"

FROM scratch
LABEL version="0.0.1"
LABEL description="Container Auditor"

COPY --from=builder /go/src/container-auditor /
COPY --from=builder /go/src/public /public/

EXPOSE 9103

ENTRYPOINT ["/container-auditor"]
