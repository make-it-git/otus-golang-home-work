FROM golang:1.18-alpine3.16 AS builder
RUN apk add make git
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build-calendar

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /build/bin/calendar /app/calendar
COPY --from=builder /build/configs/config_calendar.yaml /app/configs/config.yaml
RUN addgroup -S application && adduser -S application -G application && chmod 755 /app/calendar
USER application
ENTRYPOINT ["/app/calendar"]
CMD ["-config", "/app/configs/config.yaml"]