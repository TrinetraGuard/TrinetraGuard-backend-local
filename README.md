# TrinetraGuard Backend - Video Processing API

A high-accuracy video face detection and recognition system built with Go and Python, featuring intelligent deduplication and comprehensive storage management. This is a standalone backend API designed to work with any frontend application.

## 🏗️ Project Structure

```
Trinetr-backend/
├── 📁 api/                          # Backend API (Go + Python)
│   ├── 📁 handlers/                 # HTTP request handlers
│   │   ├── 📄 video_handlers.go     # Video upload and processing
│   │   └── 📄 storage_handlers.go   # Storage management
│   ├── 📁 models/                   # Data models and storage
│   │   ├── 📄 video_storage.go      # Video record management
│   │   └── 📄 search_history.go     # Search history management
│   ├── 📁 python/                   # Python ML components
│   │   ├── 📄 face_detect.py        # Main face detection script
│   │   ├── 📄 face_search.py        # Face search and comparison
│   │   └── 📄 requirements.txt      # Python dependencies
│   ├── 📁 venv/                     # Python virtual environment
│   ├── 📄 main.go                   # Server entry point
│   ├── 📄 go.mod                    # Go dependencies
│   └── 📄 go.sum                    # Go dependency checksums
│
├── 📁 storage/                      # Data storage
│   ├── 📁 videos/                   # Uploaded video files
│   ├── 📁 faces/                    # Extracted face images
│   ├── 📁 temp/                     # Temporary files
│   └── 📁 data/                     # JSON storage files
│       ├── 📄 videos.json           # Video records database
│       └── 📄 search_history.json   # Search history database
│
├── 📁 scripts/                      # Utility scripts
│   ├── 📁 setup/                    # Setup and installation
│   │   ├── 📄 setup.sh              # Main setup script
│   │   └── 📄 test_api.sh           # API testing script
│   └── 📁 cleanup/                  # Maintenance scripts
│       └── 📄 cleanup.py            # File cleanup utility
│
├── 📁 docs/                         # Documentation
│   ├── 📄 API_DOCUMENTATION.md      # Complete API documentation
│   ├── 📁 deployment/               # Deployment guides
│   │   └── 📄 DEPLOYMENT.md         # Deployment instructions
│   ├── 📄 README.md                 # Original README
│   └── 📄 SUMMARY.md                # System summary
│
├── 📁 config/                       # Configuration files
│   ├── 📄 Dockerfile                # Docker configuration
│   └── 📄 docker-compose.yml        # Docker Compose setup
│
└── 📄 README.md                     # This file
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Python 3.8+
- FFmpeg

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd Trinetr-backend

# Run setup script
./scripts/setup/setup.sh

# Start the server
cd api && go run main.go
```

### API Access
- **API Base URL**: http://localhost:8080
- **Health Check**: http://localhost:8080/api/health
- **API Documentation**: See [docs/API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md)

## 🎯 Features

### 🔧 Core Features
- **High-Accuracy Face Detection**: 98%+ accuracy using dlib
- **Intelligent Deduplication**: Remove duplicate faces automatically
- **Video Processing**: Support for multiple video formats
- **Local Storage**: Complete video and result management
- **Real-time Processing**: Progress tracking and status updates
- **Face Search**: Search for matching faces across all videos

### 📊 Storage Management
- **Video Records**: Complete metadata tracking
- **Face Images**: Automatic face extraction and storage
- **Statistics**: Processing metrics and analytics
- **Cleanup**: Automatic old file management
- **Search History**: Track all face searches

### 🌐 RESTful API
- **Modern API Design**: RESTful endpoints with JSON responses
- **CORS Support**: Cross-origin request support
- **File Upload**: Multipart form data support
- **Error Handling**: Comprehensive error responses
- **Health Monitoring**: Built-in health checks

## 📡 API Endpoints

### Core Endpoints
- `POST /api/upload-video` - Upload and process video
- `POST /api/search-by-face` - Search for matching faces
- `GET /api/health` - Health check

### Storage Endpoints
- `GET /api/videos` - List all videos
- `GET /api/videos/active` - List active videos
- `GET /api/videos/archived` - List archived videos
- `GET /api/videos/:id` - Get video details
- `DELETE /api/videos/:id` - Delete video
- `POST /api/videos/:id/restore` - Restore archived video
- `GET /api/videos/stats` - Get statistics
- `POST /api/videos/cleanup` - Cleanup old videos
- `GET /api/videos/search` - Search videos
- `GET /api/videos/:id/preview` - Get video preview
- `GET /api/videos/:id/file` - Download video file

### Search History Endpoints
- `GET /api/search-history` - Get search history
- `GET /api/search-history/stats` - Get search statistics

### File Serving
- `GET /api/faces/{filename}` - Serve face images

## 🐳 Docker Deployment

```bash
cd config
docker-compose up --build
```

## 📚 Documentation

- [Complete API Documentation](docs/API_DOCUMENTATION.md)
- [Deployment Guide](docs/deployment/DEPLOYMENT.md)
- [Project Structure](PROJECT_STRUCTURE.md)

## 🛠️ Development

### Project Structure Details

#### `/api` - Backend Services
- **Go Server**: RESTful API with Gin framework
- **Python ML**: Face detection and recognition
- **Storage**: Local JSON-based data management
- **Handlers**: HTTP request processing

#### `/storage` - Data Management
- **Videos**: Uploaded video files
- **Faces**: Extracted face images
- **Data**: JSON storage files
- **Temp**: Temporary processing files

#### `/scripts` - Utilities
- **Setup**: Installation and configuration
- **Cleanup**: Maintenance and cleanup

#### `/docs` - Documentation
- **API**: Complete endpoint documentation
- **Deployment**: Production guides

#### `/config` - Configuration
- **Docker**: Containerization setup
- **Environment**: Configuration files

## 🔧 Configuration

### Environment Variables
```bash
PORT=8080                    # Server port
GIN_MODE=release            # Gin mode
PYTHONPATH=/app/python      # Python path
```

### Storage Configuration
- **Video Storage**: `storage/videos/`
- **Face Storage**: `storage/faces/`
- **Data Storage**: `storage/data/videos.json`
- **Search History**: `storage/data/search_history.json`

## 🧪 Testing

```bash
# Test API endpoints
./scripts/setup/test_api.sh

# Test video upload
curl -X POST http://localhost:8080/api/upload-video \
  -F "video=@your_video.mp4" \
  -F "location_name=Office Building" \
  -F "latitude=40.7128" \
  -F "longitude=-74.0060"

# Test face search
curl -X POST http://localhost:8080/api/search-by-face \
  -F "search_image=@face.jpg"
```

## 📊 Performance

- **Face Detection Accuracy**: 98%+
- **Processing Speed**: ~1 frame/second (configurable)
- **Deduplication**: Configurable similarity threshold
- **Storage**: Efficient JSON-based storage
- **API Response Time**: < 100ms for most endpoints

## 🔌 Frontend Integration

This backend is designed to work with any frontend framework. Example integration:

```javascript
// Upload video
const formData = new FormData();
formData.append('video', videoFile);
formData.append('location_name', 'Office Building');

const response = await fetch('http://localhost:8080/api/upload-video', {
  method: 'POST',
  body: formData
});

// Search faces
const searchFormData = new FormData();
searchFormData.append('search_image', imageFile);

const searchResponse = await fetch('http://localhost:8080/api/search-by-face', {
  method: 'POST',
  body: searchFormData
});
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For support and questions:
- Check the [API Documentation](docs/API_DOCUMENTATION.md)
- Review [deployment guides](docs/deployment/)
- Open an issue for bugs or feature requests 