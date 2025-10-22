# Makefile for HNG Step 0 API

.PHONY: test test-verbose test-coverage test-helpers test-api test-nl build run clean

# Default target
all: test build

# Run all tests
test:
	@echo "Running all tests..."
	go test ./tests/...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./tests/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./tests/...

# Run tests with detailed coverage report
test-coverage-html:
	@echo "Running tests with HTML coverage report..."
	go test -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run only helper function tests
test-helpers:
	@echo "Running helper function tests..."
	go test ./tests/helpers_test.go ./tests/test_helper.go

# Run only API endpoint tests
test-api:
	@echo "Running API endpoint tests..."
	go test ./tests/api_test.go ./tests/test_helper.go

# Run only natural language filtering tests
test-nl:
	@echo "Running natural language filtering tests..."
	go test ./tests/natural_language_test.go ./tests/test_helper.go

# Build the application
build:
	@echo "Building application..."
	go build -o main .

# Run the application
run:
	@echo "Starting application..."
	go run .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f main
	rm -f coverage.out
	rm -f coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Run all checks (format, lint, test)
check: fmt lint test

# Help target
help:
	@echo "Available targets:"
	@echo "  test              - Run all tests"
	@echo "  test-verbose      - Run tests with verbose output"
	@echo "  test-coverage     - Run tests with coverage"
	@echo "  test-coverage-html- Run tests with HTML coverage report"
	@echo "  test-helpers      - Run only helper function tests"
	@echo "  test-api          - Run only API endpoint tests"
	@echo "  test-nl           - Run only natural language filtering tests"
	@echo "  build             - Build the application"
	@echo "  run               - Run the application"
	@echo "  clean             - Clean build artifacts"
	@echo "  deps              - Install dependencies"
	@echo "  fmt               - Format code"
	@echo "  lint              - Lint code"
	@echo "  check             - Run all checks (format, lint, test)"
	@echo "  help              - Show this help message"
