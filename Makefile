# MarkIt Makefile
# Development and build automation

.PHONY: help test test-coverage test-race lint fmt vet clean build install deps check benchmark security

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
deps: ## Download dependencies
	go mod download
	go mod tidy

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint
	golangci-lint run

# Testing
test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-coverage-html: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	go test -race -v ./...

benchmark: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Security
security: ## Run security scan
	gosec ./...

# Quality checks
check: fmt vet lint test ## Run all quality checks

check-ci: deps fmt vet lint test-coverage test-race security ## Run all CI checks

# Build
build: ## Build the project
	go build -v ./...

install: ## Install the project
	go install ./...

# Cleanup
clean: ## Clean build artifacts
	go clean ./...
	rm -f coverage.out coverage.html

# Documentation
docs: ## Generate documentation
	go doc -all ./...

docs-serve: ## Serve documentation locally
	godoc -http=:6060

# Git hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	@cp scripts/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed successfully"

# Release
tag: ## Create a new tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required. Usage: make tag VERSION=v1.0.0"; exit 1; fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

# Development server
dev: ## Start development mode with file watching
	@echo "Starting development mode..."
	@echo "Run 'make test' in another terminal to run tests"
	@while true; do \
		inotifywait -e modify,create,delete -r . --exclude='\.git|coverage\.|\.html$$' 2>/dev/null || sleep 1; \
		echo "Files changed, running tests..."; \
		make test; \
	done

# Performance
profile-cpu: ## Run CPU profiling
	go test -cpuprofile=cpu.prof -bench=. ./...
	go tool pprof cpu.prof

profile-mem: ## Run memory profiling
	go test -memprofile=mem.prof -bench=. ./...
	go tool pprof mem.prof

# Docker (if needed in the future)
docker-build: ## Build Docker image
	docker build -t markit:latest .

docker-test: ## Run tests in Docker
	docker run --rm -v $(PWD):/app -w /app golang:1.22 make test

# Utilities
lines: ## Count lines of code
	@find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1

todo: ## Find TODO comments
	@grep -r "TODO\|FIXME\|XXX" --include="*.go" . || echo "No TODOs found"

# Module management
mod-update: ## Update all dependencies
	go get -u ./...
	go mod tidy

mod-verify: ## Verify dependencies
	go mod verify

mod-why: ## Explain why dependencies are needed (usage: make mod-why PKG=package-name)
	@if [ -z "$(PKG)" ]; then echo "PKG is required. Usage: make mod-why PKG=package-name"; exit 1; fi
	go mod why $(PKG)

# Git utilities
git-clean: ## Clean up git branches
	git branch --merged | grep -v "\*\|main\|develop" | xargs -n 1 git branch -d

# Environment info
env: ## Show environment information
	@echo "Go version: $(shell go version)"
	@echo "Go env GOOS: $(shell go env GOOS)"
	@echo "Go env GOARCH: $(shell go env GOARCH)"
	@echo "Go env GOPATH: $(shell go env GOPATH)"
	@echo "Go env GOROOT: $(shell go env GOROOT)"
	@echo "Git version: $(shell git --version)"
	@echo "Current branch: $(shell git branch --show-current)"
	@echo "Last commit: $(shell git log -1 --oneline)" 