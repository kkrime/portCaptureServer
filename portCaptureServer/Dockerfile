FROM golang:1.20.1
ADD . /go/src/portCaptureServer
WORKDIR /go/src/portCaptureServer
RUN go get -d -v portCaptureServer
# RUN go install portCaptureServer
# ENTRYPOINT ["/go/bin/portCaptureServer"]
CMD ["go", "run", "."]
