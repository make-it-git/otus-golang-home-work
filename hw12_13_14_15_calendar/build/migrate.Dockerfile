FROM golang:1.18-alpine3.16 AS builder
RUN apk add make
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build-migrate

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /build/bin/migrate /app/migrate
COPY --from=builder /build/migrations /app/migrations
COPY --from=builder /build/configs/config_calendar.yaml /app/configs/config.yaml
RUN addgroup -S application && adduser -S application -G application && chmod 755 /app/migrate
USER application
ENTRYPOINT ["/app/migrate"]
CMD ["-config", "/app/configs/config.yaml"]