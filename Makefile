.PHONY: build run dev test clean docker-build docker-run help install setup

# Variables
APP_NAME=go-clean-gin
DOCKER_IMAGE=$(APP_NAME):latest
SERVER_PORT?=8080

# Default target
.DEFAULT_GOAL := help

## Install dependencies
install:
	go mod download
	go mod tidy

## Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/githubnemo/CompileDaemon@latest || echo "CompileDaemon installation failed"
	@go install github.com/air-verse/air@latest || go install github.com/cosmtrek/air@v1.49.0 || echo "Air installation failed"

## Setup project (first time)
setup: install install-tools
	@echo "Setting up project..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file. Please configure it."; \
	fi
	@mkdir -p tmp logs bin
	@echo "Setup complete! Run 'make dev' to start development."

## Check if port is available
check-port:
	@PORT=$${SERVER_PORT:-$(SERVER_PORT)}; \
	if lsof -i :$$PORT >/dev/null 2>&1; then \
		echo "❌ Port $$PORT is already in use"; \
		echo "Processes using port $$PORT:"; \
		lsof -i :$$PORT; \
		echo "Run 'make kill-port' to free the port"; \
		exit 1; \
	else \
		echo "✅ Port $$PORT is available"; \
	fi

## Kill process using the configured port
kill-port:
	@PORT=$${SERVER_PORT:-$(SERVER_PORT)}; \
	echo "Killing processes on port $$PORT..."; \
	sudo lsof -t -i:$$PORT | xargs kill -9 2>/dev/null || echo "No processes found on port $$PORT"

## Run the application with hot reload
dev: check-port
	@if [ -f "$(shell go env GOPATH)/bin/air" ]; then \
		echo "Using Air for hot reload..."; \
		if [ ! -f .air.toml ]; then $(shell go env GOPATH)/bin/air init; fi; \
		$(shell go env GOPATH)/bin/air; \
	elif command -v CompileDaemon >/dev/null 2>&1; then \
		echo "Using CompileDaemon for hot reload..."; \
		CompileDaemon -command="./$(APP_NAME)" -build="go build -o $(APP_NAME) cmd/main.go"; \
	else \
		echo "No hot reload available, running normally..."; \
		go run cmd/main.go; \
	fi

## Force run (kill port first)
dev-force: kill-port dev

## Run without hot reload
run:
	@echo "Running application..."
	go run cmd/main.go

## Build the application
build:
	@echo "Building application..."
	@mkdir -p bin
	go build -o bin/$(APP_NAME) cmd/main.go

## Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	go tool cover -func=coverage.out

## Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html
	rm -f *.log

## Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

## Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

## Run Docker containers
docker-run:
	@echo "Starting Docker containers..."
	docker compose up -d

## Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	docker compose down

## View Docker logs
docker-logs:
	docker compose logs -f

## Health check
health:
	@echo "Checking application health..."
	@curl -f http://localhost:$(SERVER_PORT)/health || echo "Health check failed"

## Show help
help:
	@echo "🚀 Go Clean Gin API - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  setup          Setup project (first time)"
	@echo "  dev            Run with hot reload"
	@echo "  dev-force      Kill port conflicts and run"
	@echo "  run            Run without hot reload"
	@echo ""
	@echo "Building & Testing:"
	@echo "  build          Build application"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt            Format code"
	@echo "  tidy           Tidy dependencies"
	@echo "  clean          Clean build artifacts"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Start containers"
	@echo "  docker-stop    Stop containers"
	@echo "  docker-logs    View container logs"
	@echo ""
	@echo "Utilities:"
	@echo "  check-port     Check if port is available"
	@echo "  kill-port      Kill processes on port"
	@echo "  health         Check application health"