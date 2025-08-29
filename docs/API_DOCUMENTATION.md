# TrinetraGuard Backend API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Currently, the API doesn't require authentication. All endpoints are publicly accessible.

## Endpoints

### Health Check
**GET** `/api/health`

Returns the health status of the backend.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": 1703123456
}
```

### Video Upload
**POST** `/api/upload-video`

Upload and process a video file for face detection.

**Form Data:**
- `video` (file): Video file (mp4, avi, mov, mkv, wmv, flv, webm)
- `location_name` (string, optional): Location name
- `latitude` (float, optional): Latitude coordinate
- `longitude` (float, optional): Longitude coordinate

**Response:**
```json
{
  "unique_faces_count": 5,
  "faces": ["face_1.jpg", "face_2.jpg", "face_3.jpg"],
  "message": "Video processed successfully",
  "processing_time_seconds": 12.5
}
```

### Face Search
**POST** `/api/search-by-face`

Search for matching faces across all processed videos.

**Form Data:**
- `search_image` (file): Image file (jpg, jpeg, png, bmp, gif)

**Response:**
```json
{
  "matches": [
    {
      "video": {
        "id": "video_1703123456",
        "original_filename": "sample.mp4",
        "upload_time": "2023-12-21T10:30:00Z",
        "status": "completed",
        "location_name": "Office Building",
        "latitude": 40.7128,
        "longitude": -74.0060,
        "unique_faces_count": 3,
        "processing_time": 8.2
      },
      "matched_faces": ["face_1.jpg", "face_2.jpg"],
      "similarity": 0.85
    }
  ],
  "message": "Found 1 video(s) with matching faces"
}
```

### List All Videos
**GET** `/api/videos`

Get all video records (active and archived).

**Response:**
```json
{
  "videos": [
    {
      "id": "video_1703123456",
      "original_filename": "sample.mp4",
      "upload_time": "2023-12-21T10:30:00Z",
      "status": "completed",
      "location_name": "Office Building",
      "latitude": 40.7128,
      "longitude": -74.0060,
      "unique_faces_count": 3,
      "processing_time": 8.2,
      "is_archived": false
    }
  ],
  "count": 1
}
```

### List Active Videos
**GET** `/api/videos/active`

Get only active video records.

**Response:**
```json
{
  "videos": [...],
  "count": 1,
  "type": "active"
}
```

### List Archived Videos
**GET** `/api/videos/archived`

Get only archived video records.

**Response:**
```json
{
  "videos": [...],
  "count": 0,
  "type": "archived"
}
```

### Get Video Details
**GET** `/api/videos/{id}`

Get details of a specific video.

**Response:**
```json
{
  "video": {
    "id": "video_1703123456",
    "original_filename": "sample.mp4",
    "upload_time": "2023-12-21T10:30:00Z",
    "status": "completed",
    "location_name": "Office Building",
    "latitude": 40.7128,
    "longitude": -74.0060,
    "unique_faces_count": 3,
    "processing_time": 8.2,
    "is_archived": false
  }
}
```

### Delete Video
**DELETE** `/api/videos/{id}`

Archive a video (moves to history).

**Response:**
```json
{
  "message": "Video moved to history successfully",
  "id": "video_1703123456"
}
```

### Restore Video
**POST** `/api/videos/{id}/restore`

Restore an archived video.

**Response:**
```json
{
  "message": "Video restored successfully",
  "id": "video_1703123456"
}
```

### Get Video Statistics
**GET** `/api/videos/stats`

Get storage statistics.

**Response:**
```json
{
  "stats": {
    "total_videos": 10,
    "active_videos": 8,
    "archived_videos": 2,
    "total_faces": 45,
    "total_storage_mb": 1250.5
  }
}
```

### Search Videos
**GET** `/api/videos/search`

Search videos by filename, status, or archived state.

**Query Parameters:**
- `q` (string, optional): Search query
- `status` (string, optional): Filter by status (processing, completed, failed)
- `archived` (string, optional): Filter by archived state (true, false)

**Response:**
```json
{
  "videos": [...],
  "count": 1,
  "query": "sample",
  "status": "completed",
  "archived": "false"
}
```

### Get Video Preview
**GET** `/api/videos/{id}/preview`

Get video preview information.

**Response:**
```json
{
  "video": {
    "id": "video_1703123456",
    "original_filename": "sample.mp4",
    "upload_time": "2023-12-21T10:30:00Z",
    "status": "completed",
    "location_name": "Office Building",
    "latitude": 40.7128,
    "longitude": -74.0060,
    "unique_faces": 3,
    "processing_time": 8.2,
    "video_url": "/api/videos/video_1703123456/file"
  }
}
```

### Get Video File
**GET** `/api/videos/{id}/file`

Download the actual video file.

**Response:** Video file stream

### Get Search History
**GET** `/api/search-history`

Get search history records.

**Response:**
```json
{
  "searches": [
    {
      "id": "search_1703123456",
      "search_time": "2023-12-21T10:30:00Z",
      "results_count": 2,
      "search_image": "search_1703123456.jpg"
    }
  ],
  "count": 1
}
```

### Get Search History Statistics
**GET** `/api/search-history/stats`

Get search history statistics.

**Response:**
```json
{
  "total_searches": 15,
  "total_matches": 8,
  "average_matches_per_search": 0.53,
  "most_recent_search": "2023-12-21T10:30:00Z"
}
```

### Cleanup Old Videos
**POST** `/api/videos/cleanup`

Remove very old archived records.

**Query Parameters:**
- `days` (integer, optional): Days threshold (default: 30, minimum: 7)

**Response:**
```json
{
  "message": "Cleanup completed successfully",
  "days": 30
}
```

### Reset Database
**POST** `/api/videos/reset-database`

Completely reset the database and remove all files.

**Form Data or Query Parameters:**
- `confirm` (string): Must be "true" to confirm

**Response:**
```json
{
  "message": "Database reset successfully. All videos and faces have been removed."
}
```

## Face Images

Face images are served from:
```
GET /api/faces/{filename}
```

## Error Responses

All endpoints return errors in the following format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200`: Success
- `400`: Bad Request (invalid input)
- `404`: Not Found
- `500`: Internal Server Error

## CORS

The API supports CORS and allows requests from any origin with the following headers:
- Origin
- Content-Type
- Accept
- Authorization
- X-Requested-With

## File Upload Limits

- Video files: Supported formats: mp4, avi, mov, mkv, wmv, flv, webm
- Image files: Supported formats: jpg, jpeg, png, bmp, gif

## Example Usage with JavaScript/Fetch

```javascript
// Upload a video
const formData = new FormData();
formData.append('video', videoFile);
formData.append('location_name', 'Office Building');
formData.append('latitude', '40.7128');
formData.append('longitude', '-74.0060');

const response = await fetch('http://localhost:8080/api/upload-video', {
  method: 'POST',
  body: formData
});

const result = await response.json();
console.log(result);

// Search by face
const searchFormData = new FormData();
searchFormData.append('search_image', imageFile);

const searchResponse = await fetch('http://localhost:8080/api/search-by-face', {
  method: 'POST',
  body: searchFormData
});

const searchResult = await searchResponse.json();
console.log(searchResult);
``` 