FROM golang:1.20.1
ADD ./portCaptureServer /go/src/portCaptureServer
ADD ./portCaptureServerTranslator /go/src/portCaptureServerTranslator
WORKDIR /go/src/portCaptureServerTranslator
RUN go get -d -v portCaptureServerTranslator
# RUN go install portCaptureServerTranslator
# ENTRYPOINT ["/go/bin/portCaptureServerTranslator"]
# NOTE: No compilation as requested
CMD ["go", "run", "."]
