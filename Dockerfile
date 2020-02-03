FROM golang as builder
ADD . /go/senderscore_exporter/
WORKDIR /go/senderscore_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/senderscore_exporter

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/senderscore_exporter /app/senderscore_exporter
EXPOSE 9665
VOLUME /config
ENTRYPOINT /app/senderscore_exporter -config.path=/config/config.yml