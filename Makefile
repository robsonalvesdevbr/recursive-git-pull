BINARY_NAME=rgp
MAIN_PATH=./cmd/rgp
BUILD_DIR=./bin

.PHONY: all build clean install test fmt vet run

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	@go install $(MAIN_PATH)

test:
	@echo "Running tests..."
	@go test -v ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  clean    - Clean build artifacts"
	@echo "  install  - Install to GOPATH/bin"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  vet      - Run go vet"
	@echo "  run      - Build and run with default options"
	@echo "  help     - Show this help message"