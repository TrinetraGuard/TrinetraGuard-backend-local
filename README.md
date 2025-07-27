# Video Analysis Service

A comprehensive Go backend service for video analysis, person detection, tracking, and face matching. This service provides REST API endpoints for uploading videos, analyzing them for person detection, and finding specific individuals across multiple videos using reference images.

## Features

- **Video Management**: Upload, list, and manage video files
- **Person Detection**: Automatically detect and count people in video frames
- **Person Tracking**: Track individuals across frames with unique IDs
- **Face Matching**: Find specific persons using reference images
- **REST API**: Complete REST API with Swagger documentation
- **Job Management**: Asynchronous processing with job status tracking
- **File Storage**: Organized storage for videos and reference images
- **Database**: SQLite/PostgreSQL support with automatic migrations

## Project Structure

```
video-analysis-service/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── env.example             # Environment variables example
├── README.md               # This file
├── videos/                 # Video file storage
├── finder/                 # Reference image storage
└── internal/               # Internal application code
    ├── config/             # Configuration management
    ├── database/           # Database initialization and migrations
    ├── handlers/           # HTTP request handlers
    ├── middleware/         # HTTP middleware
    ├── models/             # Data models and structures
    └── services/           # Business logic services
```

## Prerequisites

- Go 1.21 or higher
- SQLite3 (default) or PostgreSQL
- Git

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd video-analysis-service
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env file with your configuration
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

The service will start on `http://localhost:8080` by default.

## Configuration

The application can be configured using environment variables or a `.env` file:

### Server Configuration
- `ENVIRONMENT`: Application environment (development/production)
- `SERVER_PORT`: Server port (default: 8080)
- `SERVER_HOST`: Server host (default: localhost)

### Database Configuration
- `DB_DRIVER`: Database driver (sqlite3/postgres)
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSLMODE`: SSL mode for PostgreSQL

### Storage Configuration
- `STORAGE_VIDEOS_DIR`: Directory for video files
- `STORAGE_FINDER_DIR`: Directory for reference images
- `STORAGE_MAX_FILE_SIZE`: Maximum file size in bytes

### Analysis Configuration
- `ANALYSIS_MAX_CONCURRENT_JOBS`: Maximum concurrent analysis jobs
- `ANALYSIS_JOB_TIMEOUT`: Job timeout in seconds
- `ANALYSIS_FRAME_RATE`: Frame rate for analysis
- `ANALYSIS_CONFIDENCE`: Confidence threshold for detections

## API Documentation

The API documentation is available at `http://localhost:8080/swagger/index.html` when the server is running.

### Video Management Endpoints

#### Upload Video
```http
POST /api/v1/videos/upload
Content-Type: multipart/form-data

file: [video file]
```

#### List Videos
```http
GET /api/v1/videos?limit=10&offset=0
```

#### Get Video Details
```http
GET /api/v1/videos/{id}
```

#### Delete Video
```http
DELETE /api/v1/videos/{id}
```

#### Download Video
```http
GET /api/v1/videos/{id}/download
```

### Analysis Endpoints

#### Start Analysis
```http
POST /api/v1/analysis/{videoId}/start
```

#### Get Analysis Status
```http
GET /api/v1/analysis/{videoId}/status
```

#### Get Analysis Results
```http
GET /api/v1/analysis/{videoId}/results
```

#### Start Batch Analysis
```http
POST /api/v1/analysis/batch
Content-Type: application/json

["video-id-1", "video-id-2", "video-id-3"]
```

### Finder Endpoints

#### Upload Reference Image
```http
POST /api/v1/finder/upload
Content-Type: multipart/form-data

file: [image file]
description: "Person description"
```

#### List Reference Images
```http
GET /api/v1/finder/images
```

#### Search for Person
```http
POST /api/v1/finder/search
Content-Type: application/json

{
  "reference_image_id": "image-id",
  "video_ids": ["video-id-1", "video-id-2"]
}
```

#### Get Search Status
```http
GET /api/v1/finder/search/{searchId}/status
```

#### Get Search Results
```http
GET /api/v1/finder/search/{searchId}/results
```

#### Delete Reference Image
```http
DELETE /api/v1/finder/images/{id}
```

## Usage Examples

### 1. Upload and Analyze a Video

```bash
# Upload a video
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "file=@sample_video.mp4"

# Start analysis
curl -X POST http://localhost:8080/api/v1/analysis/{video-id}/start

# Check analysis status
curl http://localhost:8080/api/v1/analysis/{video-id}/status

# Get results
curl http://localhost:8080/api/v1/analysis/{video-id}/results
```

### 2. Find a Person

```bash
# Upload reference image
curl -X POST http://localhost:8080/api/v1/finder/upload \
  -F "file=@person.jpg" \
  -F "description=John Doe"

# Search for person
curl -X POST http://localhost:8080/api/v1/finder/search \
  -H "Content-Type: application/json" \
  -d '{
    "reference_image_id": "image-id",
    "video_ids": ["video-id-1", "video-id-2"]
  }'

# Check search status
curl http://localhost:8080/api/v1/finder/search/{search-id}/status

# Get search results
curl http://localhost:8080/api/v1/finder/search/{search-id}/results
```

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o video-analysis-service main.go
```

### Docker Support
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o video-analysis-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/video-analysis-service .
EXPOSE 8080
CMD ["./video-analysis-service"]
```

## Architecture

The application follows a clean architecture pattern with the following layers:

- **Handlers**: HTTP request/response handling
- **Services**: Business logic implementation
- **Models**: Data structures and database models
- **Database**: Data persistence layer
- **Config**: Configuration management
- **Middleware**: Cross-cutting concerns

### Key Components

1. **Video Service**: Handles video file operations and metadata
2. **Analysis Service**: Manages video analysis jobs and results
3. **Finder Service**: Handles person search operations
4. **Database Layer**: SQLite/PostgreSQL with automatic migrations
5. **File Storage**: Organized storage for videos and images

## Future Enhancements

- **Real-time Processing**: WebSocket support for real-time updates
- **Advanced ML Models**: Integration with more sophisticated ML models
- **Cloud Storage**: Support for cloud storage providers
- **Authentication**: JWT-based authentication and authorization
- **Rate Limiting**: API rate limiting and throttling
- **Monitoring**: Metrics and health monitoring
- **Caching**: Redis-based caching for improved performance
- **Microservices**: Split into microservices for scalability

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions, please open an issue on the GitHub repository or contact the development team. # TrinetraGuard-backend-local
