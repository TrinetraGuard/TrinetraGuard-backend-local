# Deployment Guide

This guide covers different deployment options for the Video Processing Backend.

## 🚀 Quick Start (Local Development)

### Prerequisites
- Go 1.21+
- Python 3.8+
- FFmpeg

### Setup
```bash
cd backend
./setup.sh
```

### Run
```bash
go run main.go
```

## 🐳 Docker Deployment

### Using Docker Compose (Recommended)

```bash
cd backend
docker-compose up --build
```

### Using Docker directly

```bash
cd backend
docker build -t video-processing-backend .
docker run -p 8080:8080 -v $(pwd)/videos:/app/videos -v $(pwd)/faces:/app/faces video-processing-backend
```

## ☁️ Cloud Deployment

### AWS EC2

1. **Launch EC2 Instance**
   - Use Ubuntu 20.04 LTS
   - t3.medium or larger (for ML processing)
   - At least 20GB storage

2. **Install Dependencies**
```bash
sudo apt-get update
sudo apt-get install -y docker.io docker-compose
sudo usermod -aG docker $USER
```

3. **Deploy Application**
```bash
git clone <your-repo>
cd backend
docker-compose up -d
```

### Google Cloud Run

1. **Build and Push Image**
```bash
docker build -t gcr.io/YOUR_PROJECT/video-processing-backend .
docker push gcr.io/YOUR_PROJECT/video-processing-backend
```

2. **Deploy to Cloud Run**
```bash
gcloud run deploy video-processing-backend \
  --image gcr.io/YOUR_PROJECT/video-processing-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

### Azure Container Instances

```bash
az container create \
  --resource-group myResourceGroup \
  --name video-processing-backend \
  --image video-processing-backend:latest \
  --ports 8080 \
  --dns-name-label myapp
```

## 🔧 Production Configuration

### Environment Variables

```bash
# Server configuration
PORT=8080
GIN_MODE=release

# Python configuration
PYTHONPATH=/app/python
FACE_SIMILARITY_THRESHOLD=0.6
FRAMES_PER_SECOND=1

# Storage configuration
VIDEO_STORAGE_PATH=/app/videos
FACE_STORAGE_PATH=/app/faces
```

### Security Considerations

1. **Authentication**
   - Implement JWT authentication
   - Add API key validation
   - Use HTTPS in production

2. **File Upload Limits**
   - Set maximum file size (e.g., 100MB)
   - Validate file types
   - Implement virus scanning

3. **Rate Limiting**
   - Add rate limiting middleware
   - Implement request throttling

### Monitoring and Logging

1. **Application Logs**
```bash
# Add structured logging
import "github.com/sirupsen/logrus"
```

2. **Health Checks**
```bash
curl http://localhost:8080/api/health
```

3. **Metrics**
   - Add Prometheus metrics
   - Monitor CPU and memory usage
   - Track processing times

## 📊 Performance Optimization

### Horizontal Scaling

1. **Load Balancer**
```nginx
upstream video_backend {
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}
```

2. **Database for State**
   - Use Redis for session storage
   - Store face encodings in database
   - Implement distributed processing

### Vertical Scaling

1. **GPU Acceleration**
```dockerfile
# Add CUDA support
FROM nvidia/cuda:11.0-base
```

2. **Memory Optimization**
   - Process videos in chunks
   - Implement streaming uploads
   - Use memory-mapped files

## 🔄 CI/CD Pipeline

### GitHub Actions

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build and push
        run: |
          docker build -t myapp .
          docker push myapp
      - name: Deploy
        run: |
          # Deploy to your platform
```

## 🧪 Testing

### Unit Tests
```bash
go test ./...
python -m pytest python/tests/
```

### Integration Tests
```bash
# Test API endpoints
curl -X POST http://localhost:8080/api/upload-video \
  -F "video=@test_video.mp4"
```

### Load Testing
```bash
# Using Apache Bench
ab -n 100 -c 10 -p video_file http://localhost:8080/api/upload-video
```

## 🚨 Troubleshooting

### Common Issues

1. **Out of Memory**
   - Increase container memory limits
   - Process videos in smaller chunks
   - Use streaming processing

2. **Slow Processing**
   - Use GPU acceleration
   - Increase CPU cores
   - Optimize face detection parameters

3. **File Upload Failures**
   - Check file size limits
   - Verify file permissions
   - Ensure disk space

### Debug Mode

```bash
# Enable debug logging
export GIN_MODE=debug
export PYTHONPATH=/app/python
python3 -u python/face_detect.py --debug
```

## 📈 Monitoring

### Key Metrics
- Request rate
- Processing time
- Error rate
- Memory usage
- CPU usage
- Disk usage

### Alerts
- High error rate
- Long processing times
- Low disk space
- High memory usage

## 🔒 Security Checklist

- [ ] HTTPS enabled
- [ ] Authentication implemented
- [ ] Rate limiting configured
- [ ] File upload validation
- [ ] Input sanitization
- [ ] Error handling
- [ ] Logging configured
- [ ] Monitoring setup
- [ ] Backup strategy
- [ ] Disaster recovery plan 