# TrinetraGuard Backend

e of A Go-based backend service for managing lost person reports with image upload capabilities.

## Features

- **Lost Person Reports**: Create, read, update, and delete lost person reports
- **Image Upload**: Support for JPEG, JPG, PNG, and GIF image uploads (max 10MB)
- **Search Functionality**: Search lost persons by name or Aadhar number
- **JSON Database**: Local JSON file storage for simplicity
- **RESTful API**: Clean REST API design with proper error handling
- **File Management**: Automatic file organization in uploads directory

## API Endpoints

### Lost Person Management

- `POST /api/v1/lost-persons/` - Create a new lost person report
- `GET /api/v1/lost-persons/` - Get all lost person reports
- `GET /api/v1/lost-persons/:id` - Get a specific lost person report
- `PUT /api/v1/lost-persons/:id` - Update a lost person report
- `DELETE /api/v1/lost-persons/:id` - Delete a lost person report
- `GET /api/v1/lost-persons/search?q=query` - Search lost persons
- `GET /api/v1/images/:path` - Serve uploaded images

### System Endpoints

- `GET /health` - Health check
- `GET /api/v1/system/status` - System status
- `GET /api/v1/system/info` - System information

## Data Model

Each lost person report contains:

- **ID**: Unique identifier (UUID)
- **Name**: Full name of the lost person
- **Aadhar Number**: Aadhar number (required)
- **Contact Number**: Contact number for when person is found (optional)
- **Place Lost**: Location where the person was lost
- **Permanent Address**: Permanent address of the person
- **Image Path**: Path to the uploaded image file
- **Upload Timestamp**: When the report was created

## Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd Trinetr-backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the server**:
   ```bash
   go run main.go
   ```

   The server will start on port 8080 by default. You can change this by setting the `PORT` environment variable.

## Usage Examples

### Create a Lost Person Report

```bash
curl -X POST http://localhost:8080/api/v1/lost-persons/ \
  -F "name=John Doe" \
  -F "aadhar_number=123456789012" \
  -F "contact_number=9876543210" \
  -F "place_lost=Mumbai Central Station" \
  -F "permanent_address=123 Main Street, Mumbai" \
  -F "image=@/path/to/person_photo.jpg"
```

### Get All Lost Person Reports

```bash
curl -X GET http://localhost:8080/api/v1/lost-persons/
```

### Search Lost Persons

```bash
curl -X GET "http://localhost:8080/api/v1/lost-persons/search?q=John%20Doe"
```

### Get a Specific Report

```bash
curl -X GET http://localhost:8080/api/v1/lost-persons/{id}
```

## Flutter Integration

The backend is designed to work seamlessly with Flutter applications. See `API_DOCUMENTATION.md` for detailed Flutter integration examples.

## File Storage

- **Images**: Stored in the `uploads/` directory with UUID-based filenames
- **Database**: JSON file (`database.json`) in the project root
- **File Types**: JPEG, JPG, PNG, GIF (max 10MB)

## Testing

Run the test suite:

```bash
go test ./internal/handlers/ -v
```

## Environment Variables

- `PORT`: Server port (default: 8080)
- `ENVIRONMENT`: Environment mode (development/production)

## Project Structure

```
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── database.json           # JSON database file
├── uploads/                # Image upload directory
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database layer
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Data models
│   ├── routes/            # Route definitions
│   └── utils/             # Utility functions
├── API_DOCUMENTATION.md   # Detailed API documentation
└── README.md             # This file
```

## Security Features

- File upload validation (type and size)
- Path traversal protection for image serving
- Input validation for all form fields
- Proper error handling and logging

## Contributing

Please read `CONTRIBUTING.md` for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.