# Makefile for vbook-cli cross-platform builds

VERSION ?= 1.0.0
OUTPUT_DIR ?= dist
BINARY_NAME = vbook-cli

# Build flags
LDFLAGS = -ldflags "-s -w -X main.version=$(VERSION)"
BUILD_FLAGS = CGO_ENABLED=0

# Default target
.PHONY: all
all: clean build

# Clean build artifacts
.PHONY: clean
clean:
ifeq ($(OS),Windows_NT)
	if exist $(OUTPUT_DIR) rmdir /s /q $(OUTPUT_DIR)
	if exist vbook-cli.exe del vbook-cli.exe
	if exist vbook-cli del vbook-cli
else
	rm -rf $(OUTPUT_DIR)
	rm -f vbook-cli vbook-cli.exe
endif

# Build for current platform
.PHONY: build
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Build for all platforms
.PHONY: build-all
build-all: clean
ifeq ($(OS),Windows_NT)
	@if not exist $(OUTPUT_DIR) mkdir $(OUTPUT_DIR)
else
	@mkdir -p $(OUTPUT_DIR)
endif
	@echo "Building for all platforms..."
	
	@echo "Building Windows binaries..."
ifeq ($(OS),Windows_NT)
	@set CGO_ENABLED=0&& set GOOS=windows&& set GOARCH=amd64&& go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe . && echo "[OK] Windows amd64" || echo "[FAIL] Windows amd64 failed"
	@set CGO_ENABLED=0&& set GOOS=windows&& set GOARCH=arm64&& go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-arm64.exe . && echo "[OK] Windows arm64" || echo "[FAIL] Windows arm64 failed"
else
	@$(BUILD_FLAGS) GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe . && echo "[OK] Windows amd64" || echo "[FAIL] Windows amd64 failed"
	@$(BUILD_FLAGS) GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-arm64.exe . && echo "[OK] Windows arm64" || echo "[FAIL] Windows arm64 failed"
endif
	@echo "Building Linux binaries..."
ifeq ($(OS),Windows_NT)
	@set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=amd64&& go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 . && echo "[OK] Linux amd64" || echo "[FAIL] Linux amd64 failed"
	@set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=arm64&& go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64 . && echo "[OK] Linux arm64" || echo "[FAIL] Linux arm64 failed"
else
	@$(BUILD_FLAGS) GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 . && echo "[OK] Linux amd64" || echo "[FAIL] Linux amd64 failed"
	@$(BUILD_FLAGS) GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64 . && echo "[OK] Linux arm64" || echo "[FAIL] Linux arm64 failed"
endif
	@echo "Building macOS binaries..."
ifeq ($(OS),Windows_NT)
	@set CGO_ENABLED=0&& set GOOS=darwin&& set GOARCH=amd64&& go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 . && echo "[OK] macOS amd64" || echo "[FAIL] macOS amd64 failed"
	@set CGO_ENABLED=0&& set GOOS=darwin&& set GOARCH=arm64&& go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 . && echo "[OK] macOS arm64" || echo "[FAIL] macOS arm64 failed"
else
	@$(BUILD_FLAGS) GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 . && echo "[OK] macOS amd64" || echo "[FAIL] macOS amd64 failed"
	@$(BUILD_FLAGS) GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 . && echo "[OK] macOS arm64" || echo "[FAIL] macOS arm64 failed"
endif
	
	@echo "Build complete! Check $(OUTPUT_DIR)/ for binaries."

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Development build (with debug info)
.PHONY: dev
dev:
	go build -o $(BINARY_NAME) .

# Install to GOPATH/bin
.PHONY: install
install:
	go install $(LDFLAGS) .

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build for current platform"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all supported platforms"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  dev          - Development build"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  help         - Show this help"