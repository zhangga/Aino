BIN_NAME=aino
BIN_DIR=bin

.PHONY: build run clean test

build:
	@echo "Building..."
	@go build -o $(BIN_DIR)/$(BIN_NAME) github.com/zhangga/aino

run: build
	@echo "Running..."
	@$(BIN_DIR)/$(BIN_NAME)

clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...