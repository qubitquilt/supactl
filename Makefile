.PHONY: all build test clean install lint fmt vet help

# Variables
BINARY_NAME=supactl
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_DIR=dist
GO=go
GOFLAGS=-ldflags="-s -w -X github.com/yourusername/supactl/cmd.version=$(VERSION)"

# Colors for output
BLUE=\033[0;34m
NC=\033[0m # No Color

all: clean lint test build ## Run all build steps

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary for the current platform
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)
	@echo "$(BLUE)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

build-all: ## Build binaries for all platforms
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@echo "$(BLUE)All builds complete$(NC)"

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests and generate coverage report
	@echo "$(BLUE)Generating coverage report...$(NC)"
	$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "$(BLUE)Coverage report generated: coverage.html$(NC)"

bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(NC)"
	$(GO) test -bench=. -benchmem ./...

lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run --timeout=5m

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GO) vet ./...

clean: ## Remove build artifacts
	@echo "$(BLUE)Cleaning...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f coverage.txt coverage.html

install: build ## Install the binary to GOPATH/bin
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	$(GO) install $(GOFLAGS)
	@echo "$(BLUE)Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)$(NC)"

deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	$(GO) mod download

tidy: ## Tidy go.mod
	@echo "$(BLUE)Tidying go.mod...$(NC)"
	$(GO) mod tidy

upgrade-deps: ## Upgrade dependencies
	@echo "$(BLUE)Upgrading dependencies...$(NC)"
	$(GO) get -u ./...
	$(GO) mod tidy

run: ## Run the application (requires login)
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	$(GO) run main.go

dev: ## Build and run in development mode
	@echo "$(BLUE)Running in development mode...$(NC)"
	$(GO) run main.go --help

.DEFAULT_GOAL := help
