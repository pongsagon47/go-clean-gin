#!/bin/bash

echo "Running tests..."

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

echo "Coverage report generated: coverage.html"

# Show coverage summary
go tool cover -func=coverage.out