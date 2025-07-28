# Multi-stage Docker build for Video Processing Backend
FROM python:3.9-slim as python-base

# Install system dependencies for dlib and OpenCV
RUN apt-get update && apt-get install -y \
    cmake \
    libopenblas-dev \
    liblapack-dev \
    libx11-dev \
    libgtk-3-dev \
    libavcodec-dev \
    libavformat-dev \
    libswscale-dev \
    libv4l-dev \
    libxvidcore-dev \
    libx264-dev \
    libjpeg-dev \
    libpng-dev \
    libtiff-dev \
    libatlas-base-dev \
    gfortran \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
WORKDIR /app/python
COPY python/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Go build stage
FROM golang:1.21-alpine as go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM python:3.9-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libopenblas3 \
    liblapack3 \
    libx11-6 \
    libgtk-3-0 \
    libavcodec58 \
    libavformat58 \
    libswscale5 \
    libv4l-0 \
    libxvidcore4 \
    libx264-160 \
    libjpeg62-turbo \
    libpng16-16 \
    libtiff5 \
    && rm -rf /var/lib/apt/lists/*

# Copy Python dependencies from python-base
COPY --from=python-base /usr/local/lib/python3.9/site-packages /usr/local/lib/python3.9/site-packages
COPY --from=python-base /usr/local/bin /usr/local/bin

# Copy Go binary
WORKDIR /app
COPY --from=go-builder /app/main .

# Copy Python scripts
COPY python/ ./python/

# Create directories
RUN mkdir -p videos faces

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV PYTHONPATH=/app/python

# Run the application
CMD ["./main"] 