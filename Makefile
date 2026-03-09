.PHONY: build run install clean test lint help

# Binary name
BINARY_NAME=muster
BUILD_DIR=bin

# API URLs
DEV_API_URL=http://localhost:3000
PROD_API_URL=https://api.muster.stagify.xyz

# ldflags for injecting build-time values
LDFLAGS_DEV=-ldflags "-X main.defaultBaseURL=$(DEV_API_URL)"
LDFLAGS_PROD=-ldflags "-X main.defaultBaseURL=$(PROD_API_URL)"

# Build the CLI (dev — points to localhost)
build:
	@echo "Building $(BINARY_NAME) (dev)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS_DEV) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/muster/main.go
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME) → $(DEV_API_URL)"

# Build the CLI (production — points to prod API)
build-prod:
	@echo "Building $(BINARY_NAME) (production)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS_PROD) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/muster/main.go
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME) → $(PROD_API_URL)"

# Run the CLI directly (dev)
run:
	go run $(LDFLAGS_DEV) cmd/muster/main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies installed"

# Install the binary to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Installed successfully"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Lint the code
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: brew install golangci-lint" && exit 1)
	golangci-lint run ./...

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

# Build for all platforms (production)
build-all:
	@echo "Building for all platforms (production)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS_PROD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 cmd/muster/main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS_PROD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 cmd/muster/main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS_PROD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 cmd/muster/main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS_PROD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 cmd/muster/main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS_PROD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe cmd/muster/main.go
	@echo "✓ Built binaries for all platforms in $(BUILD_DIR)/ → $(PROD_API_URL)"

# Display help
help:
	@echo "Muster CLI - Makefile commands:"
	@echo ""
	@echo "  make build         - Build the CLI binary (dev → localhost)"
	@echo "  make build-prod    - Build the CLI binary (production)"
	@echo "  make run           - Run the CLI directly (dev)"
	@echo "  make deps          - Install dependencies"
	@echo "  make install       - Install binary to /usr/local/bin"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make lint          - Run linters"
	@echo "  make fmt           - Format code"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make help          - Show this help message"
