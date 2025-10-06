.PHONY: help generate build test lint clean install-tools

help:
	@echo "Available targets:"
	@echo "  generate      - Regenerate client code from OpenAPI spec"
	@echo "  build         - Build the project"
	@echo "  test          - Run tests"
	@echo "  lint          - Run linter"
	@echo "  clean         - Clean generated files"
	@echo "  install-tools - Install required tools"

install-tools:
	@echo "Installing oapi-codegen..."
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

generate:
	@echo "Generating client from OpenAPI spec..."
	oapi-codegen -config oapi-codegen.yaml openapi.yaml

build:
	@echo "Building..."
	go build ./...

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

lint:
	@echo "Running linter..."
	golangci-lint run

clean:
	@echo "Cleaning generated files..."
	rm -f api/client.gen.go
	rm -f coverage.out
