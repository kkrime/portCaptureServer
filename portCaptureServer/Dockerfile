FROM golang:1.20.1-bullseye as builder
ADD . /go/src/portCaptureServer
WORKDIR /go/src/portCaptureServer
RUN go get -d -v portCaptureServer
RUN go install portCaptureServer

FROM golang:1.20.1-bullseye 
# security
RUN addgroup --system limited
RUN adduser --system --disabled-password --ingroup limited --home /app appuser
USER appuser

COPY --from=builder /go/bin/portCaptureServer /app/portCaptureServer
COPY --from=builder /go/src/portCaptureServer/config/docker_compose_config.yaml /app/config.yaml

ENTRYPOINT ["/app/portCaptureServer"]
