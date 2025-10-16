.PHONY: build clean run fmt test install

# Build variables
BINARY_NAME=minra
BUILD_DIR=bin
CMD_DIR=cmd/minra

# Build the project
build: clean
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run the application
run: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install ./$(CMD_DIR)
	@echo "Install complete"

# Development mode (watch and rebuild)
dev:
	@echo "Running in development mode..."
	@while true; do \
		make build; \
		./$(BUILD_DIR)/$(BINARY_NAME); \
		sleep 1; \
	done

# Help
help:
	@echo "Minra - TUI Text Editor"
	@echo ""
	@echo "Available targets:"
	@echo "  build    - Build the project"
	@echo "  clean    - Remove build artifacts"
	@echo "  run      - Build and run the application"
	@echo "  fmt      - Format code"
	@echo "  test     - Run tests"
	@echo "  deps     - Install deps"
