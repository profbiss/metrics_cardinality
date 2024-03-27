FROM golang:1.22-alpine as builder
WORKDIR /build
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app

FROM alpine
ENV LISTEN_ADDR ":8080"
ENV METRICS_URL "http://localhost:58915/metrics"

WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder app /app

CMD ["/app"]
EXPOSE 8080