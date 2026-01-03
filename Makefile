.PHONY: all build test test-cover test-integration lint format clean help

# Variables
BINARY_NAME=zai-sdk-go
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Default target
all: help

## build: Build the project
build:
	@echo "Building..."
	$(GOBUILD) -v ./...

## test: Run unit tests
test:
	@echo "Running unit tests..."
	$(GOTEST) -v -short -race ./...

## test-cover: Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -short -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

## test-integration: Run integration tests (requires API key)
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -race -tags=integration ./test/integration/...

## lint: Run linters
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: brew install golangci-lint" && exit 1)
	golangci-lint run ./...

## format: Format code
format:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@which goimports > /dev/null && goimports -w . || echo "goimports not found, skipping import formatting"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

## tidy: Tidy go modules
tidy:
	@echo "Tidying go modules..."
	$(GOMOD) tidy

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -rf bin/ dist/

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

## deps-update: Update dependencies
deps-update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

## help: Display this help message
help:
	@echo "Available targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/## /  /' | column -t -s ':'
