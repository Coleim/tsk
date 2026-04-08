.PHONY: build run test clean install lint fmt

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

# Build the binary
build:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY) ./cmd/tsk

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
	$(GOCMD) install ./cmd/tsk

# Format code
fmt:
	$(GOFMT) ./...

# Run linter
lint:
	$(GOVET) ./...

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
