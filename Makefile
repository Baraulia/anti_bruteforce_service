#Postrges
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= backend
POSTGRES_PORT ?= 5435
POSTGRES_CONTAINER ?= postgres-ab-service
#Golang
BIN_SERVICE := "./bin/service"
SERVICE_IMG="ab-service:develop"

SERVICE_CONTAINER ?= ab-service

NETWORK_NAME ?= ab_network

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)


build:
	go build -v -o $(BIN_SERVICE) -ldflags "$(LDFLAGS)" ./cmd

run-postgres:
ifneq ($(shell docker ps -q --filter "name=$(POSTGRES_CONTAINER)"),)
	@echo "Container $(POSTGRES_CONTAINER) is already running."
else
	docker run -d --name $(POSTGRES_CONTAINER) \
    	-e POSTGRES_USER=$(POSTGRES_USER) \
    	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
    	-e POSTGRES_DB=$(POSTGRES_DB) \
    	-p $(POSTGRES_PORT):5432 \
    	-v postgres-data:/var/lib/postgresql/data \
    	postgres:latest
endif

run: build run-postgres
	$(BIN_SERVICE) -config ./configs/config.yaml

build-service-img:
	docker build \
		--build-arg LDFLAGS="$(LDFLAGS)" \
		-t $(SERVICE_IMG) \
		-f build/Dockerfile .

run-service-img: build-service-img
	docker run --name $(SERVICE_CONTAINER) \
 	-e sqlHost=$(POSTGRES_CONTAINER) -e sqlPort=5432 $(SERVICE_IMG)

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest

lint: install-lint-deps
	golangci-lint run ./...

create_network:
	@docker network inspect $(NETWORK_NAME) &>/dev/null || \
    docker network create $(NETWORK_NAME)

delete_network:
	docker network rm $(NETWORK_NAME)

up:
	docker-compose  -f ./deployments/docker-compose.yaml up
down:
	docker-compose  -f ./deployments/docker-compose.yaml down


integration-tests:
	set -e ;\
	docker-compose -f ./deployments/docker-compose.test.yaml up --build -d ;\
	test_status_code=0 ;\
   	docker-compose -f ./deployments/docker-compose.test.yaml run integration_tests go test || test_status_code=$$? ;\
	docker-compose -f ./deployments/docker-compose.test.yaml down ;\
	exit $$test_status_code ;

integration-tests-teardown:
	docker-compose -f ./deployments/docker-compose.test.yaml down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose rm -f

.PHONY: build run build-img run-img version test lint
