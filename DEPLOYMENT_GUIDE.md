# 🚀 Deployment Guide - TrinetraGuard Backend

## Quick Start

### 1. **Local Development**
```bash
# Install dependencies
go mod tidy

# Start server
go run main.go

# Or use the deployment script
./deploy.sh
```

### 2. **Render Deployment (Recommended)**
- Connect your GitHub repository to Render
- Use the `render.yaml` configuration file
- Automatic deployment on every push

### 3. **Docker Deployment**
```bash
# Build and run
docker build -t trinetraguard-backend .
docker run -p 8080:8080 trinetraguard-backend
```

## 📋 Prerequisites

- Go 1.21 or higher
- Git
- Docker (optional)
- Render account (for cloud deployment)

## 🏠 Local Development

### Basic Setup
```bash
# Clone repository
git clone <your-repo-url>
cd Trinetr-backend

# Install dependencies
go mod tidy

# Start server
go run main.go
```

### Using Deployment Script
```bash
# Make executable
chmod +x deploy.sh

# Start with default settings
./deploy.sh

# Start with custom port
./deploy.sh -p 8081

# Start in production mode
./deploy.sh -e production
```

### Environment Variables
```bash
# Development
export PORT=8080
export ENVIRONMENT=development

# Production
export PORT=8080
export ENVIRONMENT=production
export GIN_MODE=release
```

## ☁️ Render Deployment

### Automatic Deployment (Recommended)

1. **Connect Repository**
   - Go to [Render Dashboard](https://dashboard.render.com)
   - Click "New +" → "Web Service"
   - Connect your GitHub repository

2. **Configure Service**
   - **Name**: `trinetraguard-backend`
   - **Environment**: `Go`
   - **Build Command**: `go mod tidy`
   - **Start Command**: `go run main.go`
   - **Plan**: Free tier (or paid for production)

3. **Environment Variables**
   - `ENVIRONMENT`: `production`
   - `GIN_MODE`: `release`

4. **Deploy**
   - Click "Create Web Service"
   - Render will automatically deploy

### Using render.yaml (Auto-deploy)

The `render.yaml` file is already configured for automatic deployment:

```yaml
services:
  - type: web
    name: trinetraguard-backend
    env: go
    plan: free
    buildCommand: go mod tidy
    startCommand: go run main.go
    envVars:
      - key: ENVIRONMENT
        value: production
      - key: GIN_MODE
        value: release
    healthCheckPath: /health
    autoDeploy: true
```

### Manual Configuration

If not using `render.yaml`:

1. **Build Command**: `go mod tidy`
2. **Start Command**: `go run main.go`
3. **Health Check Path**: `/health`
4. **Environment Variables**:
   - `ENVIRONMENT`: `production`
   - `GIN_MODE`: `release`

## 🐳 Docker Deployment

### Build Image
```bash
# Build Docker image
docker build -t trinetraguard-backend .

# Or use deployment script
./deploy.sh -b
```

### Run Container
```bash
# Basic run
docker run -p 8080:8080 trinetraguard-backend

# With environment variables
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e ENVIRONMENT=production \
  trinetraguard-backend

# Or use deployment script
./deploy.sh -d
```

### Docker Compose (Optional)
```yaml
# docker-compose.yml
version: '3.8'
services:
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - ENVIRONMENT=production
    volumes:
      - ./uploads:/root/uploads
```

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `ENVIRONMENT` | `development` | Environment mode |
| `GIN_MODE` | `debug` | Gin framework mode |

### Production Settings

```bash
# Environment variables for production
export PORT=8080
export ENVIRONMENT=production
export GIN_MODE=release
```

## 📊 Health Checks

### Health Check Endpoint
```bash
# Test health check
curl http://localhost:8080/health

# Expected response
{
  "status": "ok",
  "message": "TrinetraGuard Backend is running",
  "timestamp": "2025-08-29T07:29:20.407929Z",
  "version": "1.0.0"
}
```

### System Status
```bash
# System status
curl http://localhost:8080/api/v1/system/status

# System info
curl http://localhost:8080/api/v1/system/info
```

## 🚨 Important Notes for Production

### File Storage
- **Local Development**: Files stored in `uploads/` directory
- **Render**: Files are ephemeral (lost on restart)
- **Production Recommendation**: Use cloud storage (AWS S3, Google Cloud Storage)

### Database
- **Current**: JSON file (`database.json`)
- **Production Recommendation**: Use persistent database (PostgreSQL, MongoDB)

### CORS Configuration
Update CORS settings in `internal/routes/routes.go` for your frontend domain:

```go
// Add your frontend domain
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://your-frontend-domain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type"},
    AllowCredentials: true,
}))
```

## 🔍 Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Kill process using port 8080
   lsof -ti:8080 | xargs kill -9
   
   # Or use different port
   PORT=8081 go run main.go
   ```

2. **Permission Denied**
   ```bash
   # Make deployment script executable
   chmod +x deploy.sh
   ```

3. **Docker Build Fails**
   ```bash
   # Clean Docker cache
   docker system prune -a
   
   # Rebuild
   docker build --no-cache -t trinetraguard-backend .
   ```

4. **Render Deployment Fails**
   - Check build logs in Render dashboard
   - Ensure `go.mod` and `go.sum` are committed
   - Verify build command: `go mod tidy`

### Logs and Debugging

```bash
# View server logs
go run main.go

# Docker logs
docker logs <container-id>

# Render logs
# Check in Render dashboard under your service
```

## 📈 Monitoring

### Health Check
- **Endpoint**: `/health`
- **Expected**: 200 OK with status information
- **Frequency**: Every 30 seconds (Render)

### API Endpoints
- **Base URL**: `http://localhost:8080/api/v1`
- **Health**: `GET /health`
- **Status**: `GET /api/v1/system/status`

## 🎯 Deployment Checklist

### Before Deployment
- [ ] All tests passing: `go test ./internal/handlers/ -v`
- [ ] Health check working: `curl http://localhost:8080/health`
- [ ] API endpoints responding
- [ ] Environment variables configured
- [ ] CORS settings updated for production

### After Deployment
- [ ] Health check endpoint responding
- [ ] API endpoints accessible
- [ ] Image upload working
- [ ] Search functionality working
- [ ] Frontend can connect to backend

## 🚀 Quick Commands Reference

```bash
# Development
go run main.go
./deploy.sh

# Production
./deploy.sh -e production
docker run -p 8080:8080 trinetraguard-backend

# Testing
go test ./internal/handlers/ -v
curl http://localhost:8080/health

# Build
go build -o trinetraguard-backend main.go
docker build -t trinetraguard-backend .
```

## 📞 Support

If you encounter issues:

1. Check the logs for error messages
2. Verify environment variables are set correctly
3. Test health check endpoint
4. Review Render deployment logs
5. Check Docker build logs if using containers

**Your backend is ready for deployment!** 🎉
