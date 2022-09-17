FROM golang:1.18-alpine3.16 AS builder
RUN apk add make
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build-sender

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /build/bin/calendar_sender /app/sender
COPY --from=builder /build/configs/config_sender.yaml /app/configs/config.yaml
RUN addgroup -S application && adduser -S application -G application && chmod 755 /app/sender
USER application
ENTRYPOINT ["/app/sender"]
CMD ["-config", "/app/configs/config.yaml"]