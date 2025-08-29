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

## Installation & Setup

### Local Development

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd Trinetr-backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Start the server**:
   ```bash
   # Default port (8080)
   go run main.go
   
   # Custom port
   PORT=8081 go run main.go
   
   # Production mode
   GIN_MODE=release go run main.go
   ```

### Render Deployment

This backend is configured for deployment on Render. Follow these steps:

1. **Connect your repository to Render**
2. **Create a new Web Service**
3. **Configure the service**:
   - **Build Command**: `go mod tidy`
   - **Start Command**: `go run main.go`
   - **Environment**: Go

4. **Set Environment Variables** (optional):
   - `PORT`: Server port (Render will set this automatically)
   - `ENVIRONMENT`: `production`

### Docker Deployment

You can also deploy using Docker:

```bash
# Build the Docker image
docker build -t trinetraguard-backend .

# Run the container
docker run -p 8080:8080 trinetraguard-backend

# Run with custom port
docker run -p 8081:8080 -e PORT=8080 trinetraguard-backend
```

### Deployment Script

For convenience, use the included deployment script:

```bash
# Make script executable (first time only)
chmod +x deploy.sh

# Quick start
./deploy.sh

# Custom configurations
./deploy.sh -p 8081           # Custom port
./deploy.sh -e production     # Production mode
./deploy.sh -d                # Docker mode
./deploy.sh -b                # Build Docker image
./deploy.sh --help            # Show all options
```

## Quick Start Commands

```bash
# Install dependencies
go mod tidy

# Start server (development)
go run main.go

# Start server (custom port)
PORT=8081 go run main.go

# Start server (production mode)
GIN_MODE=release go run main.go

# Run tests
go test ./internal/handlers/ -v

# Build for production
go build -o trinetraguard-backend main.go

# Run built binary
./trinetraguard-backend

# Docker commands
docker build -t trinetraguard-backend .
docker run -p 8080:8080 trinetraguard-backend

# Using deployment script
./deploy.sh                    # Default settings
./deploy.sh -p 8081           # Custom port
./deploy.sh -e production     # Production mode
./deploy.sh -d                # Docker mode
./deploy.sh -b                # Build Docker image
./deploy.sh --help            # Show help
```

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

### Local Development
- `PORT`: Server port (default: 8080)
- `ENVIRONMENT`: Environment mode (development/production)

### Render Deployment
- `PORT`: Automatically set by Render
- `ENVIRONMENT`: Set to `production` for Render deployment

## Render Deployment Guide

### Prerequisites
- GitHub repository with your code
- Render account

### Deployment Steps

1. **Connect Repository**:
   - Go to [Render Dashboard](https://dashboard.render.com)
   - Click "New +" → "Web Service"
   - Connect your GitHub repository

2. **Configure Service**:
   - **Name**: `trinetraguard-backend` (or your preferred name)
   - **Environment**: `Go`
   - **Build Command**: `go mod tidy`
   - **Start Command**: `go run main.go`
   - **Plan**: Choose appropriate plan (Free tier available)

3. **Environment Variables** (optional):
   - `ENVIRONMENT`: `production`
   - `GIN_MODE`: `release`

4. **Deploy**:
   - Click "Create Web Service"
   - Render will automatically build and deploy your application

### Important Notes for Render

- **File Storage**: Render's file system is ephemeral. For production, consider using cloud storage (AWS S3, Google Cloud Storage) for images
- **Database**: For production, consider using a persistent database (PostgreSQL, MongoDB) instead of JSON file
- **CORS**: Update CORS settings to allow your frontend domain
- **Port**: Render automatically sets the `PORT` environment variable

### Production Considerations

For production deployment, consider these improvements:

1. **Persistent Storage**: Use cloud storage for images
2. **Database**: Use a proper database (PostgreSQL, MongoDB)
3. **Environment Variables**: Set all sensitive configuration via environment variables
4. **Logging**: Implement proper logging for production
5. **Monitoring**: Add health checks and monitoring

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

## 🚀 Deployment Summary

### **Quick Start Commands**

| **Method** | **Command** | **Description** |
|------------|-------------|-----------------|
| **Local** | `go run main.go` | Start development server |
| **Custom Port** | `PORT=8081 go run main.go` | Start on specific port |
| **Production** | `GIN_MODE=release go run main.go` | Start in production mode |
| **Script** | `./deploy.sh` | Use deployment script |
| **Docker** | `docker run -p 8080:8080 trinetraguard-backend` | Run with Docker |

### **Render Deployment (Recommended)**

1. **Connect your GitHub repository to Render**
2. **Use the included `render.yaml` for automatic deployment**
3. **Or manually configure:**
   - **Build Command**: `go mod tidy`
   - **Start Command**: `go run main.go`
   - **Environment**: Go

### **All Deployment Options**

- ✅ **Local Development**: `go run main.go`
- ✅ **Deployment Script**: `./deploy.sh`
- ✅ **Docker**: `docker build && docker run`
- ✅ **Render**: Automatic deployment with `render.yaml`
- ✅ **Production**: Environment variables configured

### **Health Check**

Test your deployment:
```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "message": "TrinetraGuard Backend is running",
  "timestamp": "2025-08-29T07:29:20.407929Z",
  "version": "1.0.0"
}
```

For detailed deployment instructions, see `DEPLOYMENT_GUIDE.md`.