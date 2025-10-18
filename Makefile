.PHONY: test bench lint fmt clean help

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

test: ## Run all tests
	go test -v -race -cover ./...

test-unit: ## Run unit tests only
	go test -v -race -cover ./tests/unit/...

test-integration: ## Run integration tests only
	go test -v -race -cover ./tests/integration/...

bench: ## Run benchmarks
	go test -bench=. -benchmem ./tests/benchmark/...

lint: ## Run linters
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	go vet ./...

coverage: ## Generate coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -f coverage.out coverage.html
	go clean -testcache

deps: ## Download dependencies
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
