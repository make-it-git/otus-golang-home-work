BIN := "./bin/calendar"
BIN_MIGRATE := "./bin/migrate"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build: build-calendar build-sender build-scheduler build-migrate

build-calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-sender:
	go build -v -o $(BIN)_sender ./cmd/sender

build-scheduler:
	go build -v -o $(BIN)_scheduler ./cmd/scheduler

build-migrate:
	go build -v -o $(BIN_MIGRATE) ./cmd/migrate.go

run: run-calendar run-sender run-scheduler run-migrate

run-calendar: build-calendar run-db
	$(BIN) -config ./configs/config.yaml

run-sender: build-sender run-db
	$(BIN)_sender -config ./configs/config_sender.yaml

run-scheduler: build-scheduler run-db
	$(BIN)_scheduler -config ./configs/config_scheduler.yaml

run-migrate: build-migrate run-db
	$(BIN_MIGRATE) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

run-db:
	docker compose -f deployments/docker-compose.yaml up -d

simple-integration-test: run-db
	sh ./tests/integration.sh

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc --go_out=internal/adapters/grpc --go-grpc_out=internal/adapters/grpc api/*.proto

.PHONY: build run build-img run-img version test lint
