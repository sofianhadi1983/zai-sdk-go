.PHONY: all build test test-cover test-integration lint lint-fix format vet clean help install-tools pre-commit-install

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
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: make install-tools" && exit 1)
	golangci-lint run --config=.golangci.yml ./...

## lint-fix: Run linters and auto-fix issues
lint-fix:
	@echo "Running linters with auto-fix..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: make install-tools" && exit 1)
	golangci-lint run --config=.golangci.yml --fix ./...

## format: Format code
format:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@which goimports > /dev/null && goimports -w -local github.com/z-ai/zai-sdk-go . || echo "goimports not found, run: make install-tools"

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

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing mockgen..."
	@go install go.uber.org/mock/mockgen@latest
	@echo "✅ All tools installed!"

## pre-commit-install: Install pre-commit hook
pre-commit-install:
	@echo "Installing pre-commit hook..."
	@cp .pre-commit-hook.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Pre-commit hook installed!"

## check: Run all checks (format, vet, lint, test)
check: format vet lint test
	@echo "✅ All checks passed!"

## help: Display this help message
help:
	@echo "Available targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/## /  /' | column -t -s ':'
