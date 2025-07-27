#!/bin/bash

# Video Analysis Service API Test Script
# This script demonstrates the basic functionality of the API

BASE_URL="http://localhost:8080/api/v1"
SERVICE_PID=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Function to start the service
start_service() {
    print_status "Starting Video Analysis Service..."
    ./bin/video-analysis-service > /dev/null 2>&1 &
    SERVICE_PID=$!
    sleep 3
    
    # Check if service is running
    if curl -f http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        print_success "Service started successfully"
    else
        print_error "Failed to start service"
        exit 1
    fi
}

# Function to stop the service
stop_service() {
    if [ ! -z "$SERVICE_PID" ]; then
        print_status "Stopping service..."
        kill $SERVICE_PID 2>/dev/null
        print_success "Service stopped"
    fi
}

# Function to test health endpoint
test_health() {
    print_status "Testing health endpoint..."
    response=$(curl -s http://localhost:8080/api/v1/health)
    if echo "$response" | grep -q "healthy"; then
        print_success "Health check passed"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        print_error "Health check failed"
    fi
}

# Function to test video listing
test_list_videos() {
    print_status "Testing video listing..."
    response=$(curl -s http://localhost:8080/api/v1/videos)
    if echo "$response" | grep -q "success"; then
        print_success "Video listing successful"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        print_error "Video listing failed"
    fi
}

# Function to test reference image listing
test_list_images() {
	print_status "Testing reference image listing..."
	response=$(curl -s http://localhost:8080/api/v1/finder/images)
	if echo "$response" | grep -q "null"; then
		print_success "Reference image listing successful (empty list)"
		echo "$response" | jq '.' 2>/dev/null || echo "$response"
	else
		print_warning "Reference image listing returned unexpected response"
		echo "$response" | jq '.' 2>/dev/null || echo "$response"
	fi
}

# Function to test Swagger documentation
test_swagger() {
    print_status "Testing Swagger documentation..."
    if curl -f http://localhost:8080/swagger/index.html > /dev/null 2>&1; then
        print_success "Swagger documentation is accessible"
        echo "Swagger UI: http://localhost:8080/swagger/index.html"
    else
        print_error "Swagger documentation is not accessible"
    fi
}

# Main test function
run_tests() {
    print_status "Running API tests..."
    echo "=================================="
    
    test_health
    echo "----------------------------------"
    
    test_list_videos
    echo "----------------------------------"
    
    test_list_images
    echo "----------------------------------"
    
    test_swagger
    echo "=================================="
    
    print_success "All tests completed!"
}

# Trap to ensure service is stopped on script exit
trap stop_service EXIT

# Check if binary exists
if [ ! -f "./bin/video-analysis-service" ]; then
    print_error "Binary not found. Please build the project first:"
    echo "  make build"
    exit 1
fi

# Main execution
echo "Video Analysis Service API Test"
echo "================================"

start_service
run_tests

print_status "Test script completed. Service will be stopped automatically." 