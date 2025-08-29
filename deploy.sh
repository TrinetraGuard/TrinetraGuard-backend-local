#!/bin/bash

# TrinetraGuard Backend Deployment Script

echo "🚀 TrinetraGuard Backend Deployment Script"
echo "=========================================="

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -p, --port PORT     Set server port (default: 8080)"
    echo "  -e, --env ENV       Set environment (development/production)"
    echo "  -d, --docker        Run with Docker"
    echo "  -b, --build         Build Docker image"
    echo "  -h, --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run with default settings"
    echo "  $0 -p 8081           # Run on port 8081"
    echo "  $0 -e production     # Run in production mode"
    echo "  $0 -d                # Run with Docker"
    echo "  $0 -b                # Build Docker image"
}

# Default values
PORT=8080
ENVIRONMENT="development"
USE_DOCKER=false
BUILD_DOCKER=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        -e|--env)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -d|--docker)
            USE_DOCKER=true
            shift
            ;;
        -b|--build)
            BUILD_DOCKER=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Build Docker image if requested
if [ "$BUILD_DOCKER" = true ]; then
    echo "🔨 Building Docker image..."
    docker build -t trinetraguard-backend .
    if [ $? -eq 0 ]; then
        echo "✅ Docker image built successfully!"
    else
        echo "❌ Docker build failed!"
        exit 1
    fi
fi

# Run the application
if [ "$USE_DOCKER" = true ]; then
    echo "🐳 Starting with Docker on port $PORT..."
    docker run -p $PORT:8080 -e PORT=8080 -e ENVIRONMENT=$ENVIRONMENT trinetraguard-backend
else
    echo "🚀 Starting Go server on port $PORT..."
    echo "Environment: $ENVIRONMENT"
    echo "Port: $PORT"
    echo ""
    
    # Set environment variables
    export PORT=$PORT
    export ENVIRONMENT=$ENVIRONMENT
    
    if [ "$ENVIRONMENT" = "production" ]; then
        export GIN_MODE=release
        echo "Running in production mode..."
    fi
    
    # Start the server
    go run main.go
fi
