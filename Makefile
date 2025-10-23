# Go A-B-C Testing Utility Makefile

# Default target
.DEFAULT_GOAL := help

# Test commands
.PHONY: test-main
test-main: ## Run all tests from main_test.go (automatically finds all Test* functions except integration)
	go test -v -run TestRules -timeout 10m

.PHONY: test-integration
test-integration: ## Run integration tests
	go test -v -run TestExecutorLogFileCreation -timeout 5m

.PHONY: test-eslint-negative
test-eslint-negative: ## Run ESLint negative case tests
	go test -v -run TestEsLintExecutorNegativeCase -timeout 3m

.PHONY: test-executors
test-executors: ## Run executor unit tests
	go test -v ./executor

.PHONY: test-analyzers
test-analyzers: ## Run log analyzer tests
	go test -v ./logsAnalyzer

.PHONY: test-all
test-all: ## Run all tests
	go test -v ./...

# Build commands
.PHONY: build
build: ## Build the project
	go build .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

# Clean commands
.PHONY: clean
clean: ## Clean execution directories
	rm -rf __executions__

# Help command
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)