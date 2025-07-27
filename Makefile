# Video Analysis Service Makefile

.PHONY: help build run test clean deps lint swagger docker-build docker-run

# Default target
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Download dependencies"
	@echo "  lint        - Run linter"
	@echo "  swagger     - Generate Swagger documentation"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run Docker container"

# Build the application
build:
	@echo "Building video-analysis-service..."
	go build -o bin/video-analysis-service main.go

# Run the application
run:
	@echo "Running video-analysis-service..."
	go run main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g main.go

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t video-analysis-service .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 -v $(PWD)/videos:/app/videos -v $(PWD)/finder:/app/finder video-analysis-service

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Setup development environment
setup: install-tools deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp env.example .env; echo "Created .env file from env.example"; fi
	@mkdir -p videos finder
	@echo "Development environment setup complete!"

# Run with hot reload (requires air)
dev:
	@echo "Running with hot reload..."
	air

# Database operations
db-migrate:
	@echo "Running database migrations..."
	go run main.go migrate

# Health check
health:
	@echo "Checking service health..."
	curl -f http://localhost:8080/api/v1/health || echo "Service is not running"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Generate mocks (if using mockery)
mocks:
	@echo "Generating mocks..."
	mockery --all

# Coverage report
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html" 