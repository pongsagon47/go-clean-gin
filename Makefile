# Makefile
.PHONY: build run test clean docker-build docker-run migrate seed dev help

# Variables
APP_NAME=go-clean-gin
DOCKER_IMAGE=$(APP_NAME):latest

# Default target
.DEFAULT_GOAL := help

## Build the application
build:
	go build -o bin/$(APP_NAME) cmd/main.go

## Run the application
run:
	go run cmd/main.go

## Run the application with hot reload
dev:
	air

## Run tests
test:
	go test -v ./...

## Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

## Format code
fmt:
	go fmt ./...

## Run golangci-lint
lint:
	golangci-lint run

## Tidy dependencies
tidy:
	go mod tidy

## Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

## Run Docker containers
docker-run:
	docker compose up -d

## Stop Docker containers
docker-stop:
	docker compose down

## View logs
docker-logs:
	docker compose logs -f

## Install dependencies
install:
	go mod download

## Database migration (if you add migration files)
migrate:
	echo "Running migrations..."

## Seed database
seed:
	echo "Seeding database..."

## Show help
help:
	@echo "Available commands:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/##//g'