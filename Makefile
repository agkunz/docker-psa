.PHONY: build install clean clean-all test format-check format lint install-tools dev-setup build-all build-linux build-darwin build-windows setup-venv setup-hooks

# Build configuration
BUILD_DIR := build
BINARY_NAME := docker-psa

build:
	@echo "Building $(BINARY_NAME) to $(BUILD_DIR)/"
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/

install: build
	mkdir -p ~/.docker/cli-plugins
	cp $(BUILD_DIR)/$(BINARY_NAME) ~/.docker/cli-plugins/$(BINARY_NAME)
	chmod +x ~/.docker/cli-plugins/$(BINARY_NAME)
	@echo "Docker PSA plugin installed. You can now run 'docker psa'"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f ~/.docker/cli-plugins/$(BINARY_NAME)
	@echo "Build artifacts cleaned"

clean-all: clean
	@echo "Cleaning all development artifacts..."
	rm -rf .venv
	@echo "All artifacts cleaned"

test:
	go test ./...

# Cross-compilation targets
build-all: build-linux build-darwin build-windows
	@echo "All platform builds completed in $(BUILD_DIR)/"

build-linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/

build-darwin:
	@echo "Building for macOS (amd64 and arm64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/

build-windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/

# Code quality targets
format-check:
	@echo "Running format checks..."
	golangci-lint run

format:
	@echo "Formatting Go code and organizing imports..."
	go fmt ./...
	gofumpt -w .
	goimports -w .
	go mod tidy

lint:
	@echo "Running linters..."
	golangci-lint run

# Development setup
setup-venv:
	@echo "Setting up Python virtual environment..."
	@if [ ! -d ".venv" ]; then \
		python3 -m venv .venv; \
		.venv/bin/pip install pre-commit; \
	fi

setup-hooks: setup-venv
	@echo "Setting up git hooks..."
	.venv/bin/pre-commit install
	.venv/bin/pre-commit install --hook-type commit-msg

install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-setup: install-tools setup-hooks
	@echo "🚀 Development environment is ready!"

.DEFAULT_GOAL := build
