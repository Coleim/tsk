.PHONY: build run test clean install lint fmt bench bench-cpu bench-mem perf-test check pre-commit

# Binary name
BINARY=tsk
BINARY_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Build the binary
build:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY) ./cmd/tsk

# Run the application
run: build
	./$(BINARY_DIR)/$(BINARY)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html

# Install to GOPATH/bin
install:
	$(GOCMD) install $(LDFLAGS) ./cmd/tsk

# Format code
fmt:
	$(GOFMT) ./...

# Run linter (requires golangci-lint: brew install golangci-lint)
lint:
	golangci-lint run ./...

# Pre-commit check (lint + test)
check: lint test

# Alias for pre-commit
pre-commit: check

# Download dependencies
deps:
	$(GOGET) -v ./...
	$(GOCMD) mod tidy

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY)-linux-amd64 ./cmd/tsk

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY)-darwin-amd64 ./cmd/tsk
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY)-darwin-arm64 ./cmd/tsk

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY)-windows-amd64.exe ./cmd/tsk

# Default target
all: fmt lint test build

# Run benchmarks
bench:
	$(GOTEST) -bench=. -benchmem -run="^$$" ./internal/ui/...

# Run benchmarks with CPU profiling
bench-cpu:
	$(GOTEST) -bench=. -benchmem -cpuprofile=cpu.prof -run="^$$" ./internal/ui/...

# Run benchmarks with memory profiling
bench-mem:
	$(GOTEST) -bench=. -benchmem -memprofile=mem.prof -run="^$$" ./internal/ui/...

# Run performance threshold tests
perf-test:
	$(GOTEST) -v -run "TestPerformanceThresholds" ./internal/ui/...

# Generate performance test board
perf-setup:
	@./scripts/perf-test.sh
