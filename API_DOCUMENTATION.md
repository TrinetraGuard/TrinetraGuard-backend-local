# Video Analysis Service - API Documentation

## 📋 Table of Contents

1. [Overview](#overview)
2. [Technologies & Architecture](#technologies--architecture)
3. [Setup Instructions](#setup-instructions)
4. [Authentication & Authorization](#authentication--authorization)
5. [Base URL & Headers](#base-url--headers)
6. [Error Handling](#error-handling)
7. [API Endpoints](#api-endpoints)
   - [System Endpoints](#system-endpoints)
   - [Video Management](#video-management)
   - [Analysis Operations](#analysis-operations)
   - [Finder Operations](#finder-operations)
8. [Database Schema](#database-schema)
9. [Postman Collection](#postman-collection)

---

## 📖 Overview

The **Video Analysis Service** is a comprehensive Go backend service designed for video analysis, person detection, tracking, and face matching. It provides a robust REST API for uploading videos, analyzing them for person detection, and finding specific individuals across multiple videos using reference images.

### Key Features
- **Video Management**: Upload, store, and manage video files
- **Person Detection**: Automatically detect and count people in video frames
- **Person Tracking**: Track individuals across frames with unique IDs
- **Face Matching**: Find specific persons using reference images
- **Asynchronous Processing**: Background job processing with status tracking
- **File Storage**: Organized storage for videos and reference images

---

## 🏗️ Technologies & Architecture

### Technology Stack
- **Backend Framework**: Go 1.21+ with Gin HTTP framework
- **Database**: SQLite (default) / PostgreSQL support
- **Logging**: Structured logging with Zap
- **Documentation**: Swagger/OpenAPI 3.0
- **Containerization**: Docker & Docker Compose
- **File Storage**: Local file system with organized directories

### Architecture Pattern
The service follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────┐
│   Handlers      │  HTTP request/response handling
├─────────────────┤
│   Services      │  Business logic implementation
├─────────────────┤
│   Models        │  Data structures and validation
├─────────────────┤
│   Database      │  Data persistence layer
├─────────────────┤
│   Config        │  Configuration management
└─────────────────┘
```

### Project Structure
```
video-analysis-service/
├── main.go                     # Application entry point
├── go.mod                      # Go module dependencies
├── env.example                 # Environment configuration
├── README.md                   # Project documentation
├── videos/                     # Video file storage
├── finder/                     # Reference image storage
└── internal/                   # Application code
    ├── config/                 # Configuration management
    ├── database/               # Database initialization
    ├── handlers/               # HTTP request handlers
    ├── middleware/             # HTTP middleware
    ├── models/                 # Data models
    └── services/               # Business logic services
```

---

## 🚀 Setup Instructions

### Prerequisites
- Go 1.21 or higher
- SQLite3 (default) or PostgreSQL
- Git

### Local Development Setup

1. **Clone and Navigate**
   ```bash
   git clone <repository-url>
   cd video-analysis-service
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Configure Environment**
   ```bash
   cp env.example .env
   # Edit .env file with your configuration
   ```

4. **Build and Run**
   ```bash
   # Build the project
   make build
   
   # Run the service
   make run
   
   # Or use the binary directly
   ./bin/video-analysis-service
   ```

5. **Verify Installation**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

### Docker Setup
```bash
# Build and run with Docker
docker-compose up

# Or build and run manually
docker build -t video-analysis-service .
docker run -p 8080:8080 video-analysis-service
```

---

## 🔐 Authentication & Authorization

**Current Status**: Authentication is not implemented in the current version. All endpoints are publicly accessible.

**Future Implementation**: The service is designed to support JWT-based authentication with the following structure:
- Bearer token authentication
- Role-based access control (User, Admin, Authority)
- API key authentication for service-to-service communication

**Security Headers**: The service includes CORS middleware configured for cross-origin requests.

---

## 🌐 Base URL & Headers

### Base URL
```
Development: http://localhost:8080
Production: https://your-domain.com
```

### API Base Path
```
/api/v1
```

### Required Headers
```http
Content-Type: application/json
Accept: application/json
X-Request-ID: <unique-request-id>  # Optional, auto-generated if not provided
```

### Optional Headers
```http
Authorization: Bearer <jwt-token>  # For future authentication
```

---

## ⚠️ Error Handling

### Standard Error Response Format
```json
{
  "success": false,
  "error": "error_code",
  "message": "Human-readable error message",
  "request_id": "unique-request-identifier"
}
```

### HTTP Status Codes
| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 413 | Payload Too Large - File size exceeds limit |
| 422 | Unprocessable Entity - Validation failed |
| 500 | Internal Server Error - Server error |
| 503 | Service Unavailable - Service temporarily unavailable |

### Common Error Codes
| Error Code | Description |
|------------|-------------|
| `invalid_file` | Invalid or missing file in upload |
| `upload_failed` | File upload processing failed |
| `not_found` | Requested resource not found |
| `invalid_id` | Invalid ID parameter |
| `job_exists` | Analysis job already exists |
| `timeout` | Request timeout |
| `internal_server_error` | Unexpected server error |

---

## 📡 API Endpoints

### System Endpoints

#### Health Check
**GET** `/api/v1/health`

Check the health status of the service.

**Response**
```json
{
  "status": "healthy",
  "service": "Video Analysis Service",
  "version": "1.0.0",
  "time": "2025-07-27T15:54:35.157107Z"
}
```

**Example**
```bash
curl -X GET http://localhost:8080/api/v1/health
```

---

### Video Management

#### Upload Video
**POST** `/api/v1/videos/upload`

Upload a video file for analysis.

**Content-Type**: `multipart/form-data`

**Form Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file` | File | Yes | Video file to upload |

**Supported Formats**: `.mp4`, `.avi`, `.mov`, `.mkv`, `.wmv`, `.flv`, `.webm`

**Max File Size**: 100MB (configurable)

**Success Response** (200)
```json
{
  "success": true,
  "message": "Video uploaded successfully",
  "video": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "filename": "550e8400-e29b-41d4-a716-446655440000.mp4",
    "original_filename": "sample_video.mp4",
    "file_size": 52428800,
    "duration": null,
    "frame_count": null,
    "width": null,
    "height": null,
    "status": "uploaded",
    "created_at": "2025-07-27T15:54:35.157107Z",
    "updated_at": "2025-07-27T15:54:35.157107Z"
  }
}
```

**Error Responses**
- `400` - Invalid file or missing file
- `413` - File size exceeds limit
- `500` - Upload processing failed

**Example**
```bash
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "file=@sample_video.mp4"
```

#### List Videos
**GET** `/api/v1/videos`

Get a paginated list of all uploaded videos.

**Query Parameters**
| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `limit` | Integer | 10 | 100 | Number of videos to return |
| `offset` | Integer | 0 | - | Number of videos to skip |

**Success Response** (200)
```json
{
  "success": true,
  "videos": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "filename": "550e8400-e29b-41d4-a716-446655440000.mp4",
      "original_filename": "sample_video.mp4",
      "file_size": 52428800,
      "duration": 120.5,
      "frame_count": 3600,
      "width": 1920,
      "height": 1080,
      "status": "analyzed",
      "created_at": "2025-07-27T15:54:35.157107Z",
      "updated_at": "2025-07-27T15:54:35.157107Z"
    }
  ],
  "total": 1
}
```

**Example**
```bash
curl -X GET "http://localhost:8080/api/v1/videos?limit=10&offset=0"
```

#### Get Video Details
**GET** `/api/v1/videos/{id}`

Get detailed information about a specific video.

**Path Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | Yes | Video ID |

**Success Response** (200)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "filename": "550e8400-e29b-41d4-a716-446655440000.mp4",
  "original_filename": "sample_video.mp4",
  "file_size": 52428800,
  "duration": 120.5,
  "frame_count": 3600,
  "width": 1920,
  "height": 1080,
  "status": "analyzed",
  "created_at": "2025-07-27T15:54:35.157107Z",
  "updated_at": "2025-07-27T15:54:35.157107Z"
}
```

**Error Responses**
- `400` - Invalid video ID
- `404` - Video not found
- `500` - Server error

**Example**
```bash
curl -X GET http://localhost:8080/api/v1/videos/550e8400-e29b-41d4-a716-446655440000
```

#### Delete Video
**DELETE** `/api/v1/videos/{id}`

Delete a video and its associated analysis results.

**Path Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | Yes | Video ID |

**Success Response** (200)
```json
{
  "success": true,
  "message": "Video deleted successfully",
  "request_id": "b47c7d3d-7ea4-4de4-910b-e354b3b6d643"
}
```

**Error Responses**
- `400` - Invalid video ID
- `404` - Video not found
- `500` - Server error

**Example**
```bash
curl -X DELETE http://localhost:8080/api/v1/videos/550e8400-e29b-41d4-a716-446655440000
```

#### Download Video
**GET** `/api/v1/videos/{id}/download`

Download the original video file.

**Path Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | Yes | Video ID |

**Success Response** (200)
- File download with original filename

**Error Responses**
- `400` - Invalid video ID
- `404` - Video not found
- `500` - Server error

**Example**
```bash
curl -X GET http://localhost:8080/api/v1/videos/550e8400-e29b-41d4-a716-446655440000/download \
  -o downloaded_video.mp4
```

---

### Analysis Operations

#### Start Analysis
**POST** `/api/v1/analysis/{videoId}/start`

Start analyzing a video for person detection and tracking.

**Path Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `videoId` | String | Yes | Video ID to analyze |

**Success Response** (200)
```json
{
  "success": true,
  "job": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "video_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "pending",
    "progress": 0,
    "error_message": null,
    "started_at": null,
    "completed_at": null,
    "created_at": "2025-07-27T15:54:35.157107Z"
  }
}
```

**Error Responses**
- `400` - Invalid video ID
- `404` - Video not found
- `409` - Analysis job already exists
- `500` - Server error

**Example**
```bash
curl -X POST http://localhost:8080/api/v1/analysis/550e8400-e29b-41d4-a716-446655440000/start
```

#### Get Analysis Status
**GET** `/api/v1/analysis/{videoId}/status`

Get the current status of a video analysis job.

**Path Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `videoId` | String | Yes | Video ID |

**Success Response** (200)
```json
{
  "success": true,
  "job": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "video_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "running",
    "progress": 45,
    "error_message": null,
    "started_at": "2025-07-27T15:54:35.157107Z",
    "completed_at": null,
    "created_at": "2025-07-27T15:54:35.157107Z"
  }
}
```

**Job Status Values**
- `pending` - Job is queued
- `running` - Job is in progress
- `completed` - Job finished successfully
- `failed` - Job failed with error
- `cancelled` - Job was cancelled

**Error Responses**
- `400` - Invalid video ID
- `404` - Analysis job not found
- `500` - Server error

**Example**
```bash
curl -X GET http://localhost:8080/api/v1/analysis/550e8400-e29b-41d4-a716-446655440000/status
```

#### Get Analysis Results
**GET** `/api/v1/analysis/{videoId}/results`

Get the results of a completed video analysis.

**Path Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `videoId` | String | Yes | Video ID |

**Success Response** (200)
```json
{
  "success": true,
  "results": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "job_id": "660e8400-e29b-41d4-a716-446655440001",
    "video_id": "550e8400-e29b-41d4-a716-446655440000",
    "total_frames": 3600,
    "total_people": 1250,
    "unique_people": 8,
    "people_per_frame": [
      {
        "frame_number": 1,
        "timestamp": 0.0,
        "count": 3,
        "confidence": 0.85
      }
    ],
    "tracking_data": [
      {
        "person_id": "person_1",
        "frame_number": 1,
        "timestamp": 0.0,
        "bounding_box": {
          "x": 100,
          "y": 150,
          "width": 80,
          "height": 200
        },
        "confidence": 0.85
      }
    ],
    "created_at": "2025-07-27T15:54:35.157107Z"
  }
}
```

**Error Responses**
- `400` - Invalid video ID
- `404` - Analysis results not found
- `500` - Server error

**Example**
```bash
curl -X GET http://localhost:8080/api/v1/analysis/550e8400-e29b-41d4-a716-446655440000/results
```

#### Start Batch Analysis
**POST** `/api/v1/analysis/batch`

Start analysis for multiple videos.

**Request Body**
```json
[
  "550e8400-e29b-41d4-a716-446655440000",
  "550e8400-e29b-41d4-a716-446655440001",
  "550e8400-e29b-41d4-a716-446655440002"
]
```

**Success Response** (200)
```json
{
  "success": true,
  "message": "Batch analysis started",
  "jobs": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "video_id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "pending",
      "progress": 0,
      "created_at": "2025-07-27T15:54:35.157107Z"
    }
  ],
  "total": 1
}
```

**Error Responses**
- `400` - Invalid request body or empty video list
- `500` - Server error

**Example**
```bash
curl -X POST http://localhost:8080/api/v1/analysis/batch \
  -H "Content-Type: application/json" \
  -d '["550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440001"]'
```

---

### Finder Operations

#### Upload Reference Image
**POST** `/api/v1/finder/upload`

Upload a reference image for person finding.

**Content-Type**: `multipart/form-data`

**Form Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file` | File | Yes | Reference image file |
| `description` | String | No | Description of the person |

**Supported Formats**: `.jpg`, `.jpeg`, `.png`, `.bmp`, `.gif`, `.tiff`, `.webp`

**Max File Size**: 100MB (configurable)

**Success Response** (200)
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440003",
  "filename": "880e8400-e29b-41d4-a716-446655440003.jpg",
  "original_filename": "person.jpg",
  "file_size": 524288,
  "description": "John Doe",
  "created_at": "2025-07-27T15:54:35.157107Z"
}
```

**Error Responses**
- `400` - Invalid file or missing file
- `413` - File size exceeds limit
- `500` - Upload processing failed

**Example**
```bash
curl -X POST http://localhost:8080/api/v1/finder/upload \
  -F "file=@person.jpg" \
  -F "description=John Doe"
```

#### List Reference Images
**GET** `/api/v1/finder/images`

Get a list of all uploaded reference images.

**Success Response** (200)
```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "filename": "880e8400-e29b-41d4-a716-446655440003.jpg",
    "original_filename": "person.jpg",
    "file_size": 524288,
    "description": "John Doe",
    "created_at": "2025-07-27T15:54:35.157107Z"
  }
]
```

**Example**
```bash
curl -X GET http://localhost:8080/api/v1/finder/images
```

#### Search for Person
**POST** `/api/v1/finder/search`

Search for a person in videos using a reference image.

**Request Body**
```json
{
  "reference_image_id": "880e8400-e29b-41d4-a716-446655440003",
  "video_ids": [
    "550e8400-e29b-41d4-a716-446655440000",
    "550e8400-e29b-41d4-a716-446655440001"
  ]
}
```

**Request Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `reference_image_id` | String | Yes | Reference image ID |
| `video_ids` | Array | No | Array of video IDs to search (empty = search all) |

**Success Response** (200)
```json
{
  "success": true,
  "message": "Person search started successfully",
  "search_job": {
    "id": "990e8400-e29b-41d4-a716-446655440004",
    "reference_image_id": "880e8400-e29b-41d4-a716-446655440003",
    "status": "pending",
    "progress": 0,
    "error_message": null,
    "started_at": null,
    "completed_at": null,
    "created_at": "2025-07-27T15:54:35.157107Z"
  }
}
```

**Error Responses**
- `400` - Invalid request body
- `404` - Reference image not found
- `500` - Server error

**Example**
```bash
curl -X POST http://localhost:8080/api/v1/finder/search \
  -H "Content-Type: application/json" \
  -d '{
    "reference_image_id": "880e8400-e29b-41d4-a716-446655440003",
    "video_ids": ["550e8400-e29b-41d4-a716-446655440000"]
  }'
```

#### Get Search Status
**GET** `/api/v1/finder/search/{searchId}/status`

Get the current status of a person search job.

**Path Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `searchId` | String | Yes | Search job ID |

**Success Response** (200)
```json
{
  "id": "990e8400-e29b-41d4-a716-446655440004",
  "reference_image_id": "880e8400-e29b-41d4-a716-446655440003",
  "status": "running",
  "progress": 60,
  "error_message": null,
  "started_at": "2025-07-27T15:54:35.157107Z",
  "completed_at": null,
  "created_at": "2025-07-27T15:54:35.157107Z"
}
```

**Error Responses**
- `400` - Invalid search ID
- `404` - Search job not found
- `500` - Server error

**Example**
```bash
curl -X GET http://localhost:8080/api/v1/finder/search/990e8400-e29b-41d4-a716-446655440004/status
```

#### Get Search Results
**GET** `/api/v1/finder/search/{searchId}/results`

Get the results of a completed person search.

**Path Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `searchId` | String | Yes | Search job ID |

**Success Response** (200)
```json
{
  "success": true,
  "results": [
    {
      "id": "aa0e8400-e29b-41d4-a716-446655440005",
      "search_job_id": "990e8400-e29b-41d4-a716-446655440004",
      "video_id": "550e8400-e29b-41d4-a716-446655440000",
      "reference_image_id": "880e8400-e29b-41d4-a716-446655440003",
      "matches": [
        {
          "frame_number": 10,
          "timestamp": 10.0,
          "bounding_box": {
            "x": 150,
            "y": 200,
            "width": 90,
            "height": 220
          },
          "confidence": 0.92
        }
      ],
      "first_appearance": 10.0,
      "last_appearance": 40.0,
      "total_appearances": 3,
      "confidence": 0.92,
      "created_at": "2025-07-27T15:54:35.157107Z"
    }
  ],
  "total": 1
}
```

**Error Responses**
- `400` - Invalid search ID
- `404` - Search results not found
- `500` - Server error

**Example**
```bash
curl -X GET http://localhost:8080/api/v1/finder/search/990e8400-e29b-41d4-a716-446655440004/results
```

#### Delete Reference Image
**DELETE** `/api/v1/finder/images/{id}`

Delete a reference image and its associated search results.

**Path Parameters**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | Yes | Reference image ID |

**Success Response** (200)
```json
{
  "success": true,
  "message": "Reference image deleted successfully",
  "request_id": "b47c7d3d-7ea4-4de4-910b-e354b3b6d643"
}
```

**Error Responses**
- `400` - Invalid image ID
- `404` - Reference image not found
- `500` - Server error

**Example**
```bash
curl -X DELETE http://localhost:8080/api/v1/finder/images/880e8400-e29b-41d4-a716-446655440003
```

---

## 🗄️ Database Schema

### Tables Overview

The service uses a relational database with the following tables:

1. **videos** - Video file metadata
2. **analysis_jobs** - Analysis job tracking
3. **analysis_results** - Analysis results storage
4. **reference_images** - Reference image metadata
5. **search_jobs** - Search job tracking
6. **search_results** - Search results storage

### Table Relationships

```
videos (1) ──── (N) analysis_jobs
analysis_jobs (1) ──── (1) analysis_results
videos (1) ──── (N) search_results
reference_images (1) ──── (N) search_jobs
search_jobs (1) ──── (N) search_results
```

### Detailed Schema

#### videos
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | TEXT | PRIMARY KEY | Unique video identifier |
| `filename` | TEXT | NOT NULL | Stored filename |
| `original_filename` | TEXT | NOT NULL | Original uploaded filename |
| `file_size` | INTEGER | NOT NULL | File size in bytes |
| `duration` | REAL | NULL | Video duration in seconds |
| `frame_count` | INTEGER | NULL | Total number of frames |
| `width` | INTEGER | NULL | Video width in pixels |
| `height` | INTEGER | NULL | Video height in pixels |
| `status` | TEXT | DEFAULT 'uploaded' | Video processing status |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

#### analysis_jobs
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | TEXT | PRIMARY KEY | Unique job identifier |
| `video_id` | TEXT | NOT NULL, FOREIGN KEY | Reference to videos table |
| `status` | TEXT | DEFAULT 'pending' | Job status |
| `progress` | INTEGER | DEFAULT 0 | Progress percentage (0-100) |
| `error_message` | TEXT | NULL | Error message if failed |
| `started_at` | TIMESTAMP | NULL | Job start timestamp |
| `completed_at` | TIMESTAMP | NULL | Job completion timestamp |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

#### analysis_results
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | TEXT | PRIMARY KEY | Unique result identifier |
| `job_id` | TEXT | NOT NULL, FOREIGN KEY | Reference to analysis_jobs |
| `video_id` | TEXT | NOT NULL, FOREIGN KEY | Reference to videos table |
| `total_frames` | INTEGER | NULL | Total frames analyzed |
| `total_people` | INTEGER | NULL | Total people detected |
| `unique_people` | INTEGER | NULL | Unique people tracked |
| `people_per_frame` | JSON | NULL | People count per frame data |
| `tracking_data` | JSON | NULL | Individual tracking data |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

#### reference_images
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | TEXT | PRIMARY KEY | Unique image identifier |
| `filename` | TEXT | NOT NULL | Stored filename |
| `original_filename` | TEXT | NOT NULL | Original uploaded filename |
| `file_size` | INTEGER | NOT NULL | File size in bytes |
| `description` | TEXT | NULL | Person description |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

#### search_jobs
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | TEXT | PRIMARY KEY | Unique search job identifier |
| `reference_image_id` | TEXT | NOT NULL, FOREIGN KEY | Reference to reference_images |
| `status` | TEXT | DEFAULT 'pending' | Job status |
| `progress` | INTEGER | DEFAULT 0 | Progress percentage (0-100) |
| `error_message` | TEXT | NULL | Error message if failed |
| `started_at` | TIMESTAMP | NULL | Job start timestamp |
| `completed_at` | TIMESTAMP | NULL | Job completion timestamp |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

#### search_results
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | TEXT | PRIMARY KEY | Unique result identifier |
| `search_job_id` | TEXT | NOT NULL, FOREIGN KEY | Reference to search_jobs |
| `video_id` | TEXT | NOT NULL, FOREIGN KEY | Reference to videos table |
| `reference_image_id` | TEXT | NOT NULL, FOREIGN KEY | Reference to reference_images |
| `matches` | JSON | NULL | Match data with timestamps |
| `first_appearance` | REAL | NULL | First appearance timestamp |
| `last_appearance` | REAL | NULL | Last appearance timestamp |
| `total_appearances` | INTEGER | NULL | Total number of appearances |
| `confidence` | REAL | NULL | Overall confidence score |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

### Indexes
- `idx_videos_status` - Optimize video status queries
- `idx_analysis_jobs_video_id` - Optimize job lookups by video
- `idx_analysis_jobs_status` - Optimize job status queries
- `idx_analysis_results_video_id` - Optimize result lookups by video
- `idx_search_jobs_status` - Optimize search job status queries
- `idx_search_results_video_id` - Optimize search result lookups by video
- `idx_search_results_reference_image_id` - Optimize search result lookups by reference image

---

## 📦 Postman Collection

### Collection Structure
```
Video Analysis Service
├── System
│   └── Health Check
├── Video Management
│   ├── Upload Video
│   ├── List Videos
│   ├── Get Video Details
│   ├── Delete Video
│   └── Download Video
├── Analysis Operations
│   ├── Start Analysis
│   ├── Get Analysis Status
│   ├── Get Analysis Results
│   └── Start Batch Analysis
└── Finder Operations
    ├── Upload Reference Image
    ├── List Reference Images
    ├── Search for Person
    ├── Get Search Status
    ├── Get Search Results
    └── Delete Reference Image
```

### Environment Variables
```json
{
  "base_url": "http://localhost:8080",
  "api_version": "v1",
  "video_id": "",
  "analysis_job_id": "",
  "reference_image_id": "",
  "search_job_id": ""
}
```

### Collection Import
You can import this collection into Postman by:

1. **Download the collection file** (provided separately)
2. **Import into Postman**:
   - Open Postman
   - Click "Import"
   - Select the collection file
   - Set up environment variables

### Pre-request Scripts
The collection includes pre-request scripts to:
- Set common headers
- Generate request IDs
- Handle authentication (when implemented)

### Tests
Each request includes tests to:
- Verify response status codes
- Validate response structure
- Extract and store dynamic values (IDs) for subsequent requests

---

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ENVIRONMENT` | development | Application environment |
| `SERVER_PORT` | 8080 | Server port |
| `SERVER_HOST` | localhost | Server host |
| `DB_DRIVER` | sqlite3 | Database driver (sqlite3/postgres) |
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 5432 | Database port |
| `DB_USER` | postgres | Database username |
| `DB_PASSWORD` | | Database password |
| `DB_NAME` | video_analysis | Database name |
| `DB_SSLMODE` | disable | SSL mode for PostgreSQL |
| `STORAGE_VIDEOS_DIR` | ./videos | Video storage directory |
| `STORAGE_FINDER_DIR` | ./finder | Reference image storage directory |
| `STORAGE_MAX_FILE_SIZE` | 104857600 | Maximum file size (100MB) |
| `ANALYSIS_MAX_CONCURRENT_JOBS` | 3 | Maximum concurrent analysis jobs |
| `ANALYSIS_JOB_TIMEOUT` | 3600 | Job timeout in seconds |
| `ANALYSIS_FRAME_RATE` | 1 | Frame rate for analysis |
| `ANALYSIS_CONFIDENCE` | 0.7 | Confidence threshold |

### File Storage
- **Videos**: Stored in `./videos/` directory
- **Reference Images**: Stored in `./finder/` directory
- **File Naming**: UUID-based filenames for security
- **Supported Formats**: Multiple video and image formats

---

## 🚀 Deployment

### Production Considerations
1. **Database**: Use PostgreSQL for production
2. **File Storage**: Consider cloud storage (S3, GCS)
3. **Authentication**: Implement JWT authentication
4. **Rate Limiting**: Enable rate limiting middleware
5. **Monitoring**: Add metrics and health monitoring
6. **SSL/TLS**: Use HTTPS in production
7. **Backup**: Implement database and file backups

### Docker Deployment
```bash
# Build production image
docker build -t video-analysis-service:latest .

# Run with environment variables
docker run -d \
  -p 8080:8080 \
  -e ENVIRONMENT=production \
  -e DB_DRIVER=postgres \
  -e DB_HOST=your-db-host \
  -e DB_NAME=video_analysis \
  -v /path/to/videos:/app/videos \
  -v /path/to/finder:/app/finder \
  video-analysis-service:latest
```

---

## 📞 Support

### Documentation
- **Interactive API Docs**: http://localhost:8080/swagger/index.html
- **Project README**: See README.md for detailed setup instructions
- **Code Examples**: See test_api.sh and demo_api.sh for usage examples

### Contact
- **Issues**: Create an issue on the GitHub repository
- **Questions**: Contact the development team
- **Contributions**: Submit pull requests following the contribution guidelines

---

## 📄 License

This project is licensed under the MIT License. See the LICENSE file for details.

---

*This documentation was generated for Video Analysis Service v1.0.0* 