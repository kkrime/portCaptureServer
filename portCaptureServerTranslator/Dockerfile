FROM golang:1.20.1-bullseye as builder
ADD ./portCaptureServer /go/src/portCaptureServer
ADD ./portCaptureServerTranslator /go/src/portCaptureServerTranslator
WORKDIR /go/src/portCaptureServerTranslator
RUN go get -d -v portCaptureServerTranslator
RUN go install portCaptureServerTranslator

FROM golang:1.20.1-bullseye 
# security
RUN addgroup --system limited
RUN adduser --system --disabled-password --ingroup limited --home /app appuser
USER appuser

COPY --from=builder /go/bin/portCaptureServerTranslator /app/portCaptureServerTranslator

ENTRYPOINT ["/app/portCaptureServerTranslator"]
