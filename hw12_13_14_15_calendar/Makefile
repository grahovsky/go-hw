BIN_CALENDAR = "./bin/calendar"
BIN_SCHEDULER = "./bin/scheduler"
BIN_SENDER = "./bin/sender"

DOCKER_IMG_CALENDAR = "calendar:develop"
DOCKER_IMG_SCHEDULER = "scheduler:develop"
DOCKER_IMG_SENDER = "sender:develop"
DOCKER_IMG_MIGRATIONS = "migrations:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DBUSER ?= "postgres"
DBPASS ?= "postgres"

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler

build-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

build-all: build-calendar build-scheduler build-sender

build: build-all

run-calendar: build-calendar
	$(BIN_CALENDAR) --config ./configs/calendar.yaml

run-scheduler: build-scheduler
	$(BIN_SCHEDULER) --config ./configs/scheduler.yaml

run-sender: build-sender
	$(BIN_SENDER) --config ./configs/sender.yaml

build-img-calendar:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_CALENDAR) \
		-f build/calendar.Dockerfile .

build-img-scheduler:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SCHEDULER) \
		-f build/scheduler.Dockerfile .

build-img-sender:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f build/sender.Dockerfile .

build-img-all: build-img-calendar build-img-scheduler build-img-sender

run-img-calendar: build-img-calendar
	docker run --network host $(DOCKER_IMG_CALENDAR)

run-img-scheduler: build-img-scheduler
	docker run --network host $(DOCKER_IMG_SCHEDULER)

run-img-sender: build-img-sender
	docker run --network host $(DOCKER_IMG_SENDER)

run-img-all: run-img-calendar run-img-scheduler run-img-sender

version: build-calendar
	$(BIN_CALENDAR) --version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run --timeout=90s ./...

install-migrate-deps:
	(which goose > /dev/null) || go install github.com/pressly/goose/v3/cmd/goose@latest

migrate: install-migrate-deps
	goose --dir=./migrations postgres "host=0.0.0.0 user=${DBUSER} password=${DBPASS} dbname=calendar sslmode=disable" up

generate: 
	go generate ./...

up:
	docker compose -f deployments/calendar-compose.yaml --project-directory deployments up -d --build

rebuild: 
	docker compose -f deployments/calendar-compose.yaml --project-directory deployments build --no-cache

logs:
	docker compose -f deployments/calendar-compose.yaml logs -f

down:
	docker compose -f deployments/calendar-compose.yaml --project-directory deployments down

clean:
	docker compose -f deployments/calendar-compose.yaml --project-directory deployments down --rmi all -v

.PHONY: build-all run-calendar run-img-all version test lint
