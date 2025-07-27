#!/bin/bash

# Video Analysis Service API Demo
# This script demonstrates the API functionality with actual operations

BASE_URL="http://localhost:8080/api/v1"
SERVICE_PID=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Start the service
start_service() {
    print_status "Starting Video Analysis Service..."
    ./bin/video-analysis-service > /dev/null 2>&1 &
    SERVICE_PID=$!
    sleep 3
    
    if curl -f http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        print_success "Service started successfully"
    else
        print_error "Failed to start service"
        exit 1
    fi
}

# Stop the service
stop_service() {
    if [ ! -z "$SERVICE_PID" ]; then
        print_status "Stopping service..."
        kill $SERVICE_PID 2>/dev/null
        print_success "Service stopped"
    fi
}

# Demo 1: Health Check
demo_health() {
    print_status "Demo 1: Health Check"
    echo "GET $BASE_URL/health"
    response=$(curl -s http://localhost:8080/api/v1/health)
    echo "$response" | jq '.'
    echo ""
}

# Demo 2: List Videos (Empty)
demo_list_videos() {
    print_status "Demo 2: List Videos (Empty Database)"
    echo "GET $BASE_URL/videos"
    response=$(curl -s http://localhost:8080/api/v1/videos)
    echo "$response" | jq '.'
    echo ""
}

# Demo 3: List Reference Images (Empty)
demo_list_images() {
    print_status "Demo 3: List Reference Images (Empty Database)"
    echo "GET $BASE_URL/finder/images"
    response=$(curl -s http://localhost:8080/api/v1/finder/images)
    echo "$response" | jq '.'
    echo ""
}

# Demo 4: Try to Start Analysis on Non-existent Video
demo_analysis_error() {
    print_status "Demo 4: Try to Start Analysis on Non-existent Video"
    echo "POST $BASE_URL/analysis/non-existent-id/start"
    response=$(curl -s -X POST http://localhost:8080/api/v1/analysis/non-existent-id/start)
    echo "$response" | jq '.'
    echo ""
}

# Demo 5: Try to Search with Non-existent Reference Image
demo_search_error() {
    print_status "Demo 5: Try to Search with Non-existent Reference Image"
    echo "POST $BASE_URL/finder/search"
    echo '{"reference_image_id": "non-existent-id", "video_ids": []}'
    response=$(curl -s -X POST http://localhost:8080/api/v1/finder/search \
        -H "Content-Type: application/json" \
        -d '{"reference_image_id": "non-existent-id", "video_ids": []}')
    echo "$response" | jq '.'
    echo ""
}

# Demo 6: Batch Analysis with Empty Video List
demo_batch_analysis() {
    print_status "Demo 6: Batch Analysis with Empty Video List"
    echo "POST $BASE_URL/analysis/batch"
    echo '[]'
    response=$(curl -s -X POST http://localhost:8080/api/v1/analysis/batch \
        -H "Content-Type: application/json" \
        -d '[]')
    echo "$response" | jq '.'
    echo ""
}

# Demo 7: Swagger Documentation
demo_swagger() {
    print_status "Demo 7: Swagger Documentation"
    echo "Swagger UI is available at: http://localhost:8080/swagger/index.html"
    echo "You can access the interactive API documentation there."
    echo ""
}

# Main demo function
run_demo() {
    print_status "Running API Demo..."
    echo "=================================="
    
    demo_health
    demo_list_videos
    demo_list_images
    demo_analysis_error
    demo_search_error
    demo_batch_analysis
    demo_swagger
    
    print_success "Demo completed!"
    echo ""
    print_status "Next steps:"
    echo "1. Upload a video file using: POST $BASE_URL/videos/upload"
    echo "2. Upload a reference image using: POST $BASE_URL/finder/upload"
    echo "3. Start analysis using: POST $BASE_URL/analysis/{video-id}/start"
    echo "4. Search for a person using: POST $BASE_URL/finder/search"
    echo "5. View full documentation at: http://localhost:8080/swagger/index.html"
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
echo "Video Analysis Service API Demo"
echo "================================"

start_service
run_demo

print_status "Demo completed. Service will be stopped automatically." 