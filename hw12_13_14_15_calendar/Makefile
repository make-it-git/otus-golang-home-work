BIN := "./bin/calendar"
BIN_MIGRATE := "./bin/migrate"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build: build-calendar build-sender build-scheduler build-migrate

build-calendar:
	CGO_ENABLED=0 go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-sender:
	CGO_ENABLED=0 go build -v -o $(BIN)_sender ./cmd/sender

build-scheduler:
	CGO_ENABLED=0 go build -v -o $(BIN)_scheduler ./cmd/scheduler

build-migrate:
	CGO_ENABLED=0 go build -v -o $(BIN_MIGRATE) ./cmd/migrate.go

test:
	go test -race ./internal/...

up:
	docker compose -f deployments/docker-compose.yaml up -d

down:
	docker compose -f deployments/docker-compose.yaml down

integration-tests: up
	go test -race ./tests --tags integration

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc --go_out=internal/adapters/grpc --go-grpc_out=internal/adapters/grpc api/*.proto

.PHONY: build build-calendar build-sender build-scheduler build-migrate test up simple-integration-test install-lint-deps lint generate
