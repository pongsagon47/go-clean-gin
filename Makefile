.PHONY: build run test clean docker-build docker-run migrate seed dev help install setup check-port kill-port dev-force air

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
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest || echo "golangci-lint installation failed"

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

## Install Air for hot reload
air:
	@echo "Trying to install Air..."
	@go install github.com/air-verse/air@latest || \
	 go install github.com/cosmtrek/air@v1.49.0 || \
	 go install github.com/githubnemo/CompileDaemon@latest || \
	 echo "Hot reload tools installation failed"

## Run the application with hot reload
dev: check-port
	@if [ -f "$(shell go env GOPATH)/bin/air" ]; then \
		echo "Using Air for hot reload..."; \
		if [ ! -f .air.toml ]; then $(shell go env GOPATH)/bin/air init; fi; \
		$(shell go env GOPATH)/bin/air; \
	elif command -v CompileDaemon >/dev/null 2>&1; then \
		echo "Using CompileDaemon for hot reload..."; \
		CompileDaemon -command="./$(APP_NAME)" -build="go build -o $(APP_NAME) cmd/main.go"; \
	elif [ -f scripts/dev.sh ]; then \
		echo "Using custom hot reload script..."; \
		chmod +x scripts/dev.sh; \
		./scripts/dev.sh; \
	else \
		echo "No hot reload available, running normally..."; \
		go run cmd/main.go; \
	fi

## Force run (kill port first)
dev-force: kill-port dev

## Run without hot reload
dev-simple:
	@echo "Running without hot reload..."
	go run cmd/main.go

## Build the application
build:
	@echo "Building application..."
	@mkdir -p bin
	go build -o bin/$(APP_NAME) cmd/main.go

## Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux cmd/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME)-windows.exe cmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP_NAME)-darwin cmd/main.go
	@echo "✅ Multi-platform build completed!"

## Run the application
run: build
	./bin/$(APP_NAME)

## Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@chmod +x scripts/test.sh
	@./scripts/test.sh

## Run integration tests
test-integration:
	@echo "Running integration tests..."
	@if [ -f test/integration_test.go ]; then \
		go test -v -tags=integration ./test/...; \
	else \
		echo "No integration tests found"; \
	fi

## Run specific test
test-unit:
	@echo "Running unit tests only..."
	go test -v -short ./internal/...

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

## Run golangci-lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Run 'make install-tools' first"; \
	fi

## Fix linting issues
lint-fix:
	@echo "Fixing linting issues..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix; \
	else \
		echo "golangci-lint not found. Run 'make install-tools' first"; \
	fi

## Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

## Update dependencies
update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

## Security check
security:
	@echo "Running security check..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

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

## Restart Docker containers
docker-restart: docker-stop docker-run

## Database migration (placeholder)
migrate:
	@echo "Running database migrations..."
	@echo "Implement your migration logic here"

## Seed database (placeholder)
seed:
	@echo "Seeding database..."
	@echo "Implement your seeding logic here"

## Generate API documentation (placeholder)
docs:
	@echo "Generating API documentation..."
	@echo "Implement documentation generation here"

## Health check
health:
	@echo "Checking application health..."
	@curl -f http://localhost:$(SERVER_PORT)/health || echo "Health check failed"

## Show application logs
logs:
	@if [ -f logs/app.log ]; then \
		tail -f logs/app.log; \
	else \
		echo "No log file found"; \
	fi

## Benchmark tests
benchmark:
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem ./...

## Profile CPU
profile-cpu:
	@echo "Running CPU profiling..."
	go test -cpuprofile=cpu.prof -bench=. ./...
	go tool pprof cpu.prof

## Profile memory
profile-mem:
	@echo "Running memory profiling..."
	go test -memprofile=mem.prof -bench=. ./...
	go tool pprof mem.prof

## Show help
help:
	@echo "🚀 Go Clean Gin API - Available commands:"
	@echo ""
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/##//g' | sort