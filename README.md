# Trinetr-backend - Video Processing System

A high-accuracy video face detection and recognition system built with Go and Python, featuring intelligent deduplication and comprehensive storage management.

## 🏗️ Project Structure

```
Trinetr-backend/
├── 📁 api/                          # Backend API (Go + Python)
│   ├── 📁 handlers/                 # HTTP request handlers
│   │   ├── 📄 video_handlers.go     # Video upload and processing
│   │   └── 📄 storage_handlers.go   # Storage management
│   ├── 📁 models/                   # Data models and storage
│   │   └── 📄 video_storage.go      # Video record management
│   ├── 📁 middleware/               # Custom middleware (empty)
│   ├── 📁 utils/                    # Utility functions (empty)
│   ├── 📁 python/                   # Python ML components
│   │   ├── 📄 face_detect.py        # Main face detection script
│   │   └── 📄 requirements.txt      # Python dependencies
│   ├── 📁 venv/                     # Python virtual environment
│   ├── 📄 main.go                   # Server entry point
│   ├── 📄 go.mod                    # Go dependencies
│   └── 📄 go.sum                    # Go dependency checksums
│
├── 📁 frontend/                     # Web interface
│   ├── 📁 pages/                    # HTML pages
│   │   ├── 📄 index.html            # Main upload interface
│   │   └── 📄 storage.html          # Storage management interface
│   └── 📁 assets/                   # CSS, JS, images (empty)
│
├── 📁 storage/                      # Data storage
│   ├── 📁 videos/                   # Uploaded video files
│   ├── 📁 faces/                    # Extracted face images
│   └── 📁 data/                     # JSON storage files
│       └── 📄 videos.json           # Video records database
│
├── 📁 scripts/                      # Utility scripts
│   ├── 📁 setup/                    # Setup and installation
│   │   ├── 📄 setup.sh              # Main setup script
│   │   └── 📄 test_api.sh           # API testing script
│   └── 📁 cleanup/                  # Maintenance scripts
│       └── 📄 cleanup.py            # File cleanup utility
│
├── 📁 docs/                         # Documentation
│   ├── 📁 api/                      # API documentation (empty)
│   ├── 📁 deployment/               # Deployment guides
│   │   └── 📄 DEPLOYMENT.md         # Deployment instructions
│   ├── 📁 user-guide/               # User guides (empty)
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

### Access the Application
- **Main Interface**: http://localhost:8080
- **Storage Management**: http://localhost:8080/storage
- **API Health**: http://localhost:8080/api/health

## 🎯 Features

### 🔧 Core Features
- **High-Accuracy Face Detection**: 98%+ accuracy using dlib
- **Intelligent Deduplication**: Remove duplicate faces automatically
- **Video Processing**: Support for multiple video formats
- **Local Storage**: Complete video and result management
- **Real-time Processing**: Progress tracking and status updates

### 📊 Storage Management
- **Video Records**: Complete metadata tracking
- **Face Images**: Automatic face extraction and storage
- **Statistics**: Processing metrics and analytics
- **Cleanup**: Automatic old file management

### 🌐 Web Interface
- **Modern UI**: Beautiful, responsive design
- **Drag & Drop**: Easy video upload
- **Progress Tracking**: Real-time processing status
- **Results Display**: Grid view of detected faces
- **Storage Dashboard**: Complete video management

## 📡 API Endpoints

### Core Endpoints
- `POST /api/upload-video` - Upload and process video
- `GET /api/health` - Health check

### Storage Endpoints
- `GET /api/videos` - List all videos
- `GET /api/videos/:id` - Get video details
- `DELETE /api/videos/:id` - Delete video
- `GET /api/videos/stats` - Get statistics
- `POST /api/videos/cleanup` - Cleanup old videos

## 🐳 Docker Deployment

```bash
cd config
docker-compose up --build
```

## 📚 Documentation

- [API Documentation](docs/api/)
- [Deployment Guide](docs/deployment/)
- [User Guide](docs/user-guide/)

## 🛠️ Development

### Project Structure Details

#### `/api` - Backend Services
- **Go Server**: RESTful API with Gin framework
- **Python ML**: Face detection and recognition
- **Storage**: Local JSON-based data management

#### `/frontend` - Web Interface
- **Pages**: HTML templates
- **Assets**: CSS, JavaScript, images
- **Responsive**: Mobile-friendly design

#### `/storage` - Data Management
- **Videos**: Uploaded video files
- **Faces**: Extracted face images
- **Data**: JSON storage files

#### `/scripts` - Utilities
- **Setup**: Installation and configuration
- **Cleanup**: Maintenance and cleanup

#### `/docs` - Documentation
- **API**: Endpoint documentation
- **Deployment**: Production guides
- **User Guide**: Usage instructions

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

## 🧪 Testing

```bash
# Test API endpoints
./scripts/setup/test_api.sh

# Test video upload
curl -X POST http://localhost:8080/api/upload-video \
  -F "video=@your_video.mp4"
```

## 📊 Performance

- **Face Detection Accuracy**: 98%+
- **Processing Speed**: ~1 frame/second (configurable)
- **Deduplication**: Configurable similarity threshold
- **Storage**: Efficient JSON-based storage

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
- Check the [documentation](docs/)
- Review [deployment guides](docs/deployment/)
- Open an issue for bugs or feature requests 