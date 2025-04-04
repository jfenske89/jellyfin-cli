BINARY_NAME := jellyfin-cli
VERSION := 1.0.0
BUILD_DIR := bin

.PHONY: all build clean install lint test fmt upgrade-deps

all: clean fmt lint build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags="-X 'github.com/jfenske89/jellyfin-cli/pkg/cmd.Version=$(VERSION)'" main.go

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

lint:
	@echo "Linting code..."
	@golangci-lint run

test:
	@echo "Running tests..."
	@go test -v ./...

fmt:
	@echo "Formatting code..."
	@goimports -w --local github.com/jfenske89 ./

upgrade-deps:
	@echo "Upgrading dependencies..."
	@go get -u ./...
	@go mod tidy