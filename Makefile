.PHONY: help build run clean test tcp-client http-client deps

help:
	@echo "whereyouat - Geolocation RPC Service"
	@echo ""
	@echo "Usage:"
	@echo "  make build        Build server and client binaries"
	@echo "  make run          Run the server"
	@echo "  make tcp-client   Run TCP client example"
	@echo "  make test         Run tests"
	@echo "  make clean        Remove build artifacts"
	@echo "  make deps         Download and tidy dependencies"

build:
	@echo "Building binaries..."
	@go build -o bin/whereyouat ./cmd/whereyouat
	@go build -o bin/client ./cmd/client
	@echo "Build complete!"

run:
	@go run ./cmd/whereyouat

tcp-client:
	@go run ./cmd/client

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete!"

test:
	@go test -v ./...

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated!"
