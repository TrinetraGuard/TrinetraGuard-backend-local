# TrinetraGuard Lost Person Backend API Documentation

## Overview

This backend provides RESTful API endpoints for managing lost person reports. It supports file uploads for person images and stores data in a local JSON database.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, no authentication is required for the API endpoints.

## Endpoints

### 1. Create Lost Person Report

**POST** `/lost-persons`

Creates a new lost person report with image upload.

**Request:**
- Content-Type: `multipart/form-data`
- Body: Form data with the following fields:
  - `name` (required): Full name of the lost person
  - `aadhar_number` (required): Aadhar number of the lost person
  - `contact_number` (optional): Contact number for when person is found
  - `place_lost` (required): Place where the person was lost
  - `permanent_address` (required): Permanent address of the person
  - `image` (required): Image file (JPEG, JPG, PNG, GIF, max 10MB)

**Example using curl:**
```bash
curl -X POST http://localhost:8080/api/v1/lost-persons \
  -F "name=John Doe" \
  -F "aadhar_number=123456789012" \
  -F "contact_number=9876543210" \
  -F "place_lost=Mumbai Central Station" \
  -F "permanent_address=123 Main Street, Mumbai, Maharashtra" \
  -F "image=@/path/to/person_photo.jpg"
```

**Response:**
```json
{
  "success": true,
  "message": "Lost person report created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "aadhar_number": "123456789012",
    "contact_number": "9876543210",
    "place_lost": "Mumbai Central Station",
    "permanent_address": "123 Main Street, Mumbai, Maharashtra",
    "image_path": "uploads/550e8400-e29b-41d4-a716-446655440000.jpg",
    "upload_timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Get All Lost Person Reports

**GET** `/lost-persons`

Retrieves all lost person reports.

**Response:**
```json
{
  "success": true,
  "message": "Lost person reports retrieved successfully",
  "data": {
    "total": 2,
    "lost_persons": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "John Doe",
        "aadhar_number": "123456789012",
        "contact_number": "9876543210",
        "place_lost": "Mumbai Central Station",
        "permanent_address": "123 Main Street, Mumbai, Maharashtra",
        "image_path": "uploads/550e8400-e29b-41d4-a716-446655440000.jpg",
        "upload_timestamp": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

### 3. Get Lost Person by ID

**GET** `/lost-persons/{id}`

Retrieves a specific lost person report by ID.

**Parameters:**
- `id` (path): Unique identifier of the lost person

**Response:**
```json
{
  "success": true,
  "message": "Lost person report retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "aadhar_number": "123456789012",
    "contact_number": "9876543210",
    "place_lost": "Mumbai Central Station",
    "permanent_address": "123 Main Street, Mumbai, Maharashtra",
    "image_path": "uploads/550e8400-e29b-41d4-a716-446655440000.jpg",
    "upload_timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### 4. Update Lost Person Report

**PUT** `/lost-persons/{id}`

Updates an existing lost person report.

**Parameters:**
- `id` (path): Unique identifier of the lost person

**Request:**
- Content-Type: `multipart/form-data`
- Body: Same as create endpoint (all fields optional except image)

**Response:**
```json
{
  "success": true,
  "message": "Lost person report updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe Updated",
    "aadhar_number": "123456789012",
    "contact_number": "9876543210",
    "place_lost": "Mumbai Central Station",
    "permanent_address": "123 Main Street, Mumbai, Maharashtra",
    "image_path": "uploads/550e8400-e29b-41d4-a716-446655440000.jpg",
    "upload_timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### 5. Delete Lost Person Report

**DELETE** `/lost-persons/{id}`

Deletes a lost person report and its associated image.

**Parameters:**
- `id` (path): Unique identifier of the lost person

**Response:**
```json
{
  "success": true,
  "message": "Lost person report deleted successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### 6. Search Lost Persons

**GET** `/lost-persons/search?q={query}`

Searches for lost persons by name or Aadhar number.

**Parameters:**
- `q` (query): Search query (name or Aadhar number)

**Response:**
```json
{
  "success": true,
  "message": "Search completed successfully",
  "data": {
    "total": 1,
    "lost_persons": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "John Doe",
        "aadhar_number": "123456789012",
        "contact_number": "9876543210",
        "place_lost": "Mumbai Central Station",
        "permanent_address": "123 Main Street, Mumbai, Maharashtra",
        "image_path": "uploads/550e8400-e29b-41d4-a716-446655440000.jpg",
        "upload_timestamp": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

### 7. Get Image

**GET** `/images/{path}`

Serves the uploaded image file.

**Parameters:**
- `path` (path): Image filename

**Response:**
- Content-Type: `image/jpeg` (or appropriate image type)
- Body: Image file content

## Error Responses

All endpoints return consistent error responses:

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

### Common HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## File Upload Requirements

- **Supported formats**: JPEG, JPG, PNG, GIF
- **Maximum file size**: 10MB
- **Field name**: `image`

## Data Storage

- **Database**: Local JSON file (`database.json`)
- **Images**: Stored in `uploads/` directory
- **File naming**: UUID-based unique filenames

## Flutter Integration

### Example Flutter code for creating a lost person report:

```dart
import 'dart:io';
import 'package:http/http.dart' as http;

Future<void> createLostPersonReport({
  required String name,
  required String aadharNumber,
  String? contactNumber,
  required String placeLost,
  required String permanentAddress,
  required File imageFile,
}) async {
  var request = http.MultipartRequest(
    'POST',
    Uri.parse('http://localhost:8080/api/v1/lost-persons'),
  );

  // Add form fields
  request.fields['name'] = name;
  request.fields['aadhar_number'] = aadharNumber;
  if (contactNumber != null) {
    request.fields['contact_number'] = contactNumber;
  }
  request.fields['place_lost'] = placeLost;
  request.fields['permanent_address'] = permanentAddress;

  // Add image file
  request.files.add(
    await http.MultipartFile.fromPath(
      'image',
      imageFile.path,
    ),
  );

  var response = await request.send();
  var responseData = await response.stream.bytesToString();
  
  if (response.statusCode == 201) {
    print('Report created successfully: $responseData');
  } else {
    print('Error creating report: $responseData');
  }
}
```

## Running the Backend

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Run the server:**
   ```bash
   go run main.go
   ```

3. **Server will start on port 8080** (configurable via PORT environment variable)

## Testing

Run the tests:
```bash
go test ./internal/handlers/
```

## Environment Variables

- `PORT`: Server port (default: 8080)
- `ENVIRONMENT`: Environment mode (development/production)

## Security Considerations

- File upload validation (type and size)
- Path traversal protection for image serving
- Input validation for all form fields
- Proper error handling and logging
