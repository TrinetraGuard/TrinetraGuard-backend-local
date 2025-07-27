#!/bin/bash

# Video Analysis Service - Complete Project Setup Script

echo "🎥 Video Analysis Service - Complete Project Setup"
echo "=================================================="

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

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if Go is installed
check_go() {
    print_status "Checking Go installation..."
    if command -v go > /dev/null 2>&1; then
        GO_VERSION=$(go version | awk '{print $3}')
        print_success "Go is installed: $GO_VERSION"
        return 0
    else
        print_error "Go is not installed. Please install Go 1.21 or later."
        echo "Download from: https://golang.org/dl/"
        return 1
    fi
}

# Check if required directories exist
setup_directories() {
    print_status "Setting up project directories..."
    
    # Create necessary directories
    mkdir -p videos finder bin logs
    
    # Create .gitkeep files to ensure directories are tracked
    touch videos/.gitkeep finder/.gitkeep logs/.gitkeep
    
    print_success "Directories created: videos/, finder/, bin/, logs/"
}

# Setup environment file
setup_environment() {
    print_status "Setting up environment configuration..."
    
    if [ ! -f ".env" ]; then
        if [ -f "env.example" ]; then
            cp env.example .env
            print_success "Created .env file from env.example"
        else
            print_error "env.example not found. Creating basic .env file..."
            cat > .env << EOF
# Server Configuration
ENVIRONMENT=development
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_DRIVER=sqlite3
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=video_analysis.db
DB_SSLMODE=disable

# Storage Configuration
STORAGE_VIDEOS_DIR=./videos
STORAGE_FINDER_DIR=./finder
STORAGE_MAX_FILE_SIZE=104857600

# Analysis Configuration
ANALYSIS_MAX_CONCURRENT_JOBS=3
ANALYSIS_JOB_TIMEOUT=3600
ANALYSIS_FRAME_RATE=1
ANALYSIS_CONFIDENCE=0.7
EOF
            print_success "Created basic .env file"
        fi
    else
        print_success ".env file already exists"
    fi
}

# Install Go dependencies
install_dependencies() {
    print_status "Installing Go dependencies..."
    
    if [ -f "go.mod" ]; then
        go mod tidy
        if [ $? -eq 0 ]; then
            print_success "Go dependencies installed successfully"
        else
            print_error "Failed to install Go dependencies"
            return 1
        fi
    else
        print_error "go.mod file not found"
        return 1
    fi
}

# Build the application
build_application() {
    print_status "Building the application..."
    
    make build
    if [ $? -eq 0 ]; then
        print_success "Application built successfully"
        return 0
    else
        print_error "Failed to build application"
        return 1
    fi
}

# Test the application
test_application() {
    print_status "Testing the application..."
    
    # Check if binary exists
    if [ ! -f "./bin/video-analysis-service" ]; then
        print_error "Binary not found. Building first..."
        build_application
        if [ $? -ne 0 ]; then
            return 1
        fi
    fi
    
    # Test basic functionality
    print_status "Starting application for testing..."
    ./bin/video-analysis-service > /dev/null 2>&1 &
    APP_PID=$!
    
    # Wait for application to start
    sleep 3
    
    # Test health endpoint
    if curl -f http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        print_success "Application is running and responding"
        
        # Get health response
        HEALTH_RESPONSE=$(curl -s http://localhost:8080/api/v1/health)
        echo "   Health Status: $(echo "$HEALTH_RESPONSE" | jq -r '.status')"
        echo "   Service: $(echo "$HEALTH_RESPONSE" | jq -r '.service')"
        echo "   Version: $(echo "$HEALTH_RESPONSE" | jq -r '.version')"
        
        # Stop the application
        kill $APP_PID 2>/dev/null
        print_success "Application test completed successfully"
        return 0
    else
        print_error "Application failed to start or respond"
        kill $APP_PID 2>/dev/null
        return 1
    fi
}

# Setup frontend
setup_frontend() {
    print_status "Setting up frontend..."
    
    if [ -d "frontend" ]; then
        if [ -f "frontend/index.html" ]; then
            print_success "Frontend files found"
            
            # Make scripts executable
            chmod +x start_frontend.sh demo_frontend.sh 2>/dev/null
            
            print_success "Frontend setup completed"
            return 0
        else
            print_error "Frontend index.html not found"
            return 1
        fi
    else
        print_error "Frontend directory not found"
        return 1
    fi
}

# Create comprehensive README
create_readme() {
    print_status "Creating comprehensive project documentation..."
    
    cat > PROJECT_SETUP.md << 'EOF'
# Video Analysis Service - Complete Project Setup

## 🎯 Project Overview

This is a complete Video Analysis Service with both backend API and frontend interface. The service provides:

- **Video Upload & Management**: Upload, store, and manage video files
- **Person Detection & Analysis**: Analyze videos to detect and count people
- **Person Tracking**: Track individuals across video frames with unique IDs
- **Person Search**: Find specific people across multiple videos using reference images
- **REST API**: Complete RESTful API with Swagger documentation
- **Web Frontend**: Modern, responsive web interface

## 🏗️ Architecture

```
Video Analysis Service/
├── backend/                 # Go backend service
│   ├── main.go             # Application entry point
│   ├── internal/           # Internal packages
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database operations
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── middleware/     # HTTP middleware
│   │   ├── models/         # Data models
│   │   └── services/       # Business logic
│   ├── videos/             # Video file storage
│   ├── finder/             # Reference image storage
│   └── bin/                # Compiled binaries
├── frontend/               # HTML/CSS/JS frontend
│   ├── index.html          # Main frontend application
│   ├── styles.css          # Additional styles
│   └── README.md           # Frontend documentation
└── docs/                   # Documentation
    ├── API_DOCUMENTATION.md
    ├── API_SUMMARY.md
    └── PROJECT_SUMMARY.md
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**: [Download here](https://golang.org/dl/)
- **Modern Browser**: Chrome, Firefox, Safari, Edge
- **curl**: For API testing (usually pre-installed)
- **jq**: For JSON processing (optional but recommended)

### 1. Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd video-analysis-service

# Run the setup script
./setup_project.sh
```

### 2. Start the Service

```bash
# Option 1: Start everything with one command
./start_frontend.sh

# Option 2: Manual start
make build
./bin/video-analysis-service
```

### 3. Access the Application

- **Frontend**: `file:///path/to/frontend/index.html`
- **Backend API**: `http://localhost:8080/api/v1`
- **Swagger Docs**: `http://localhost:8080/swagger/index.html`

## 📋 Complete Workflow

### 1. Upload a Video
1. Open the frontend
2. Go to "Videos" tab
3. Click "Choose video file"
4. Select a video (MP4, AVI, MOV, MKV, WMV, FLV, WEBM)
5. Click "Upload Video"

### 2. Analyze the Video
1. Go to "Analysis" tab
2. Select the uploaded video
3. Click "Start Analysis"
4. Monitor progress with "Check Status"
5. View results with "Get Results"

### 3. Upload Reference Image
1. Go to "Person Finder" tab
2. Click "Choose reference image"
3. Select an image (JPG, PNG, BMP, GIF, TIFF, WEBP)
4. Add description (optional)
5. Click "Upload Image"

### 4. Search for Person
1. Select the reference image
2. Select videos to search in
3. Click "Search for Person"
4. Monitor progress and view results

## 🔧 Configuration

### Environment Variables (.env)

```bash
# Server Configuration
ENVIRONMENT=development
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_DRIVER=sqlite3
DB_NAME=video_analysis.db

# Storage Configuration
STORAGE_VIDEOS_DIR=./videos
STORAGE_FINDER_DIR=./finder
STORAGE_MAX_FILE_SIZE=104857600

# Analysis Configuration
ANALYSIS_MAX_CONCURRENT_JOBS=3
ANALYSIS_JOB_TIMEOUT=3600
ANALYSIS_FRAME_RATE=1
ANALYSIS_CONFIDENCE=0.7
```

### Database

The service uses SQLite by default for simplicity. The database file is created automatically at `video_analysis.db`.

For production, you can switch to PostgreSQL by updating the `.env` file:

```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=video_analysis
```

## 📚 API Documentation

### Available Endpoints

**System:**
- `GET /api/v1/health` - Health check

**Videos:**
- `POST /api/v1/videos/upload` - Upload video
- `GET /api/v1/videos` - List videos
- `GET /api/v1/videos/{id}` - Get video details
- `DELETE /api/v1/videos/{id}` - Delete video
- `GET /api/v1/videos/{id}/download` - Download video

**Analysis:**
- `POST /api/v1/analysis/{videoId}/start` - Start analysis
- `GET /api/v1/analysis/{videoId}/status` - Get analysis status
- `GET /api/v1/analysis/{videoId}/results` - Get analysis results
- `POST /api/v1/analysis/batch` - Batch analysis

**Finder:**
- `POST /api/v1/finder/upload` - Upload reference image
- `GET /api/v1/finder/images` - List reference images
- `POST /api/v1/finder/search` - Search for person
- `GET /api/v1/finder/search/{id}/status` - Get search status
- `GET /api/v1/finder/search/{id}/results` - Get search results
- `DELETE /api/v1/finder/images/{id}` - Delete reference image

### Testing with curl

```bash
# Health check
curl http://localhost:8080/api/v1/health

# List videos
curl http://localhost:8080/api/v1/videos

# Upload video
curl -X POST -F "file=@video.mp4" http://localhost:8080/api/v1/videos/upload

# Start analysis
curl -X POST http://localhost:8080/api/v1/analysis/{videoId}/start
```

## 🛠️ Development

### Project Structure

```
internal/
├── config/         # Configuration loading and validation
├── database/       # Database connection and migrations
├── handlers/       # HTTP request handlers with Swagger docs
├── middleware/     # CORS, logging, authentication middleware
├── models/         # Data structures and database models
└── services/       # Business logic and external integrations
```

### Building

```bash
# Build the application
make build

# Run tests
make test

# Run with hot reload (development)
make dev

# Clean build artifacts
make clean
```

### Adding New Features

1. **Add Models**: Define data structures in `internal/models/`
2. **Add Services**: Implement business logic in `internal/services/`
3. **Add Handlers**: Create HTTP handlers in `internal/handlers/`
4. **Add Routes**: Register routes in `main.go`
5. **Update Frontend**: Add UI components in `frontend/index.html`

## 🐛 Troubleshooting

### Common Issues

**Port 8080 already in use:**
```bash
# Find and kill the process
lsof -ti:8080 | xargs kill -9
```

**Database errors:**
```bash
# Remove and recreate database
rm video_analysis.db
./bin/video-analysis-service
```

**Frontend not loading:**
- Check if backend is running
- Verify CORS configuration
- Check browser console for errors

**File upload issues:**
- Check file size limits (100MB max)
- Verify file format is supported
- Check storage directory permissions

### Debug Mode

Enable debug logging by setting in `.env`:
```bash
ENVIRONMENT=development
```

### Logs

Application logs are written to stdout. For production, consider redirecting to files:
```bash
./bin/video-analysis-service > app.log 2>&1
```

## 📦 Deployment

### Docker Deployment

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
docker-compose up -d
```

### Production Considerations

1. **Environment**: Set `ENVIRONMENT=production`
2. **Database**: Use PostgreSQL for production
3. **Security**: Enable authentication and HTTPS
4. **Monitoring**: Add health checks and metrics
5. **Backup**: Implement database and file backups
6. **Scaling**: Use load balancers and multiple instances

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For issues and questions:
1. Check the troubleshooting section
2. Review the API documentation
3. Check the logs for error messages
4. Open an issue on GitHub

---

*Last Updated: 2025-07-27*
EOF

    print_success "Created PROJECT_SETUP.md with comprehensive documentation"
}

# Main setup function
main_setup() {
    echo ""
    print_status "Starting complete project setup..."
    echo ""
    
    # Check prerequisites
    if ! check_go; then
        exit 1
    fi
    
    # Setup project structure
    setup_directories
    
    # Setup environment
    setup_environment
    
    # Install dependencies
    if ! install_dependencies; then
        exit 1
    fi
    
    # Build application
    if ! build_application; then
        exit 1
    fi
    
    # Test application
    if ! test_application; then
        exit 1
    fi
    
    # Setup frontend
    if ! setup_frontend; then
        exit 1
    fi
    
    # Create documentation
    create_readme
    
    echo ""
    print_success "🎉 Project setup completed successfully!"
    echo ""
    echo "📋 What's been set up:"
    echo "  ✅ Go dependencies installed"
    echo "  ✅ Environment configuration created"
    echo "  ✅ Application built and tested"
    echo "  ✅ Frontend configured"
    echo "  ✅ Project directories created"
    echo "  ✅ Documentation generated"
    echo ""
    echo "🚀 Next Steps:"
    echo "  1. Start the service: ./start_frontend.sh"
    echo "  2. Open the frontend in your browser"
    echo "  3. Upload a video and start analyzing"
    echo "  4. Upload reference images and search for people"
    echo ""
    echo "📚 Documentation:"
    echo "  • PROJECT_SETUP.md - Complete setup guide"
    echo "  • API_DOCUMENTATION.md - API reference"
    echo "  • frontend/README.md - Frontend guide"
    echo "  • Swagger UI: http://localhost:8080/swagger/index.html"
    echo ""
    echo "🔧 Available Commands:"
    echo "  • ./start_frontend.sh - Start everything"
    echo "  • ./demo_frontend.sh - Run demo"
    echo "  • make build - Build application"
    echo "  • make test - Run tests"
    echo "  • make clean - Clean build artifacts"
    echo ""
}

# Run the setup
main_setup 