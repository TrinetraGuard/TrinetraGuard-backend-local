# Video Analysis Service - API Quick Reference

## 🚀 Quick Start

**Base URL**: `http://localhost:8080`  
**API Version**: `v1`  
**Full Base Path**: `http://localhost:8080/api/v1`

## 📋 All Endpoints Summary

### System Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |

### Video Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/videos/upload` | Upload video file |
| `GET` | `/videos` | List all videos |
| `GET` | `/videos/{id}` | Get video details |
| `DELETE` | `/videos/{id}` | Delete video |
| `GET` | `/videos/{id}/download` | Download video |

### Analysis Operations
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/analysis/{videoId}/start` | Start video analysis |
| `GET` | `/analysis/{videoId}/status` | Get analysis status |
| `GET` | `/analysis/{videoId}/results` | Get analysis results |
| `POST` | `/analysis/batch` | Start batch analysis |

### Finder Operations
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/finder/upload` | Upload reference image |
| `GET` | `/finder/images` | List reference images |
| `POST` | `/finder/search` | Search for person |
| `GET` | `/finder/search/{searchId}/status` | Get search status |
| `GET` | `/finder/search/{searchId}/results` | Get search results |
| `DELETE` | `/finder/images/{id}` | Delete reference image |

## 🔧 Common Request Headers

```http
Content-Type: application/json
Accept: application/json
X-Request-ID: <unique-request-id>  # Optional, auto-generated
```

## 📊 Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully"
}
```

### Error Response
```json
{
  "success": false,
  "error": "error_code",
  "message": "Human-readable error message",
  "request_id": "unique-request-identifier"
}
```

## 🚨 HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 404 | Not Found |
| 409 | Conflict |
| 413 | Payload Too Large |
| 500 | Internal Server Error |

## 📝 Quick Examples

### Upload a Video
```bash
curl -X POST http://localhost:8080/api/v1/videos/upload \
  -F "file=@sample_video.mp4"
```

### Start Analysis
```bash
curl -X POST http://localhost:8080/api/v1/analysis/{video-id}/start
```

### Search for Person
```bash
curl -X POST http://localhost:8080/api/v1/finder/search \
  -H "Content-Type: application/json" \
  -d '{
    "reference_image_id": "image-id",
    "video_ids": ["video-id-1", "video-id-2"]
  }'
```

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

## 📚 Documentation Links

- **Full API Documentation**: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
- **Interactive Swagger UI**: http://localhost:8080/swagger/index.html
- **Postman Collection**: [Video_Analysis_Service.postman_collection.json](./Video_Analysis_Service.postman_collection.json)
- **Postman Environment**: [Video_Analysis_Service.postman_environment.json](./Video_Analysis_Service.postman_environment.json)

## 🔐 Authentication

**Current Status**: No authentication required  
**Future**: JWT Bearer token authentication will be implemented

## 📁 File Upload Limits

- **Videos**: 100MB max, formats: mp4, avi, mov, mkv, wmv, flv, webm
- **Images**: 100MB max, formats: jpg, jpeg, png, bmp, gif, tiff, webp

## 🗄️ Database Schema

### Core Tables
- `videos` - Video metadata
- `analysis_jobs` - Analysis job tracking
- `analysis_results` - Analysis results
- `reference_images` - Reference image metadata
- `search_jobs` - Search job tracking
- `search_results` - Search results

### Key Relationships
- Videos → Analysis Jobs → Analysis Results
- Reference Images → Search Jobs → Search Results
- Videos ↔ Search Results (many-to-many)

## 🛠️ Development Tools

### Testing Scripts
```bash
# Run basic API tests
./test_api.sh

# Run comprehensive demo
./demo_api.sh
```

### Build & Run
```bash
# Build the project
make build

# Run the service
make run

# Run with Docker
docker-compose up
```

## 📞 Support

- **Issues**: Create GitHub issue
- **Documentation**: See README.md and API_DOCUMENTATION.md
- **Interactive Testing**: Use Swagger UI or Postman collection

---

*Last Updated: 2025-07-27*  
*Version: 1.0.0* 