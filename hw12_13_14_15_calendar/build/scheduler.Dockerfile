FROM golang:1.18-alpine3.16 AS builder
RUN apk add make
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build-scheduler

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /build/bin/calendar_scheduler /app/scheduler
COPY --from=builder /build/configs/config_scheduler.yaml /app/configs/config.yaml
RUN addgroup -S application && adduser -S application -G application && chmod 755 /app/scheduler
USER application
ENTRYPOINT ["/app/scheduler"]
CMD ["-config", "/app/configs/config.yaml"]