package services

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"video-analysis-service/internal/config"
	"video-analysis-service/internal/models"

	"github.com/google/uuid"
)

// VideoService handles video-related operations
type VideoService struct {
	db  *sql.DB
	cfg *config.Config
}

// NewVideoService creates a new video service
func NewVideoService(db *sql.DB, cfg *config.Config) *VideoService {
	return &VideoService{
		db:  db,
		cfg: cfg,
	}
}

// UploadVideo uploads and stores a video file
func (s *VideoService) UploadVideo(file *multipart.FileHeader) (*models.Video, error) {
	// Validate file type
	if !s.isValidVideoType(file.Filename) {
		return nil, fmt.Errorf("unsupported video format")
	}

	// Validate file size
	if file.Size > s.cfg.Storage.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size")
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	filePath := filepath.Join(s.cfg.Storage.VideosDir, filename)

	// Ensure directory exists
	if err := os.MkdirAll(s.cfg.Storage.VideosDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create videos directory: %v", err)
	}

	// Save file
	if err := s.saveUploadedFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Create video record
	video := &models.Video{
		ID:           uuid.New().String(),
		Filename:     filename,
		OriginalName: file.Filename,
		FilePath:     filePath,
		FileSize:     file.Size,
		Status:       models.VideoStatusUploaded,
		UploadedAt:   time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	_, err := s.db.Exec(`
		INSERT INTO videos (id, filename, original_name, file_path, file_size, status, uploaded_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, video.ID, video.Filename, video.OriginalName, video.FilePath, video.FileSize, video.Status, video.UploadedAt, video.UpdatedAt)

	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save video record: %v", err)
	}

	return video, nil
}

// GetVideo retrieves a video by ID
func (s *VideoService) GetVideo(videoID string) (*models.Video, error) {
	var video models.Video
	var duration sql.NullFloat64
	var frameCount sql.NullInt64
	var width sql.NullInt64
	var height sql.NullInt64

	err := s.db.QueryRow(`
		SELECT id, filename, original_name, file_path, file_size, duration, frame_count, width, height, status, uploaded_at, updated_at
		FROM videos WHERE id = ?
	`, videoID).Scan(&video.ID, &video.Filename, &video.OriginalName, &video.FilePath, &video.FileSize, &duration, &frameCount, &width, &height, &video.Status, &video.UploadedAt, &video.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("video not found: %v", err)
	}

	// Handle nullable fields
	if duration.Valid {
		video.Duration = duration.Float64
	}
	if frameCount.Valid {
		video.FrameCount = int(frameCount.Int64)
	}
	if width.Valid {
		video.Width = int(width.Int64)
	}
	if height.Valid {
		video.Height = int(height.Int64)
	}

	return &video, nil
}

// ListVideos retrieves all videos with pagination
func (s *VideoService) ListVideos(limit, offset int) ([]models.Video, int, error) {
	// Get total count
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM videos").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get video count: %v", err)
	}

	// Get videos with pagination
	rows, err := s.db.Query(`
		SELECT id, filename, original_name, file_path, file_size, duration, frame_count, width, height, status, uploaded_at, updated_at
		FROM videos ORDER BY uploaded_at DESC LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query videos: %v", err)
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var video models.Video
		var duration sql.NullFloat64
		var frameCount sql.NullInt64
		var width sql.NullInt64
		var height sql.NullInt64

		err := rows.Scan(&video.ID, &video.Filename, &video.OriginalName, &video.FilePath, &video.FileSize, &duration, &frameCount, &width, &height, &video.Status, &video.UploadedAt, &video.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan video: %v", err)
		}

		// Handle nullable fields
		if duration.Valid {
			video.Duration = duration.Float64
		}
		if frameCount.Valid {
			video.FrameCount = int(frameCount.Int64)
		}
		if width.Valid {
			video.Width = int(width.Int64)
		}
		if height.Valid {
			video.Height = int(height.Int64)
		}

		videos = append(videos, video)
	}

	if videos == nil {
		return []models.Video{}, total, nil
	}

	return videos, total, nil
}

// DeleteVideo deletes a video and its associated file
func (s *VideoService) DeleteVideo(videoID string) error {
	// Get video info first
	var filePath string
	err := s.db.QueryRow("SELECT file_path FROM videos WHERE id = ?", videoID).Scan(&filePath)
	if err != nil {
		return fmt.Errorf("video not found: %v", err)
	}

	// Delete from database
	_, err = s.db.Exec("DELETE FROM videos WHERE id = ?", videoID)
	if err != nil {
		return fmt.Errorf("failed to delete video: %v", err)
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to delete file %s: %v\n", filePath, err)
	}

	return nil
}

// GetVideoFilePath gets the file path for a video
func (s *VideoService) GetVideoFilePath(videoID string) (string, error) {
	var filePath string
	err := s.db.QueryRow("SELECT file_path FROM videos WHERE id = ?", videoID).Scan(&filePath)
	if err != nil {
		return "", fmt.Errorf("video not found: %v", err)
	}
	return filePath, nil
}

// UpdateVideoStatus updates the status of a video
func (s *VideoService) UpdateVideoStatus(videoID, status string) error {
	_, err := s.db.Exec("UPDATE videos SET status = ?, updated_at = ? WHERE id = ?", status, time.Now(), videoID)
	if err != nil {
		return fmt.Errorf("failed to update video status: %v", err)
	}
	return nil
}

// UpdateVideoMetadata updates video metadata
func (s *VideoService) UpdateVideoMetadata(videoID string, duration float64, frameCount int, width, height int) error {
	_, err := s.db.Exec(`
		UPDATE videos SET duration = ?, frame_count = ?, width = ?, height = ?, updated_at = ? WHERE id = ?
	`, duration, frameCount, width, height, time.Now(), videoID)
	if err != nil {
		return fmt.Errorf("failed to update video metadata: %v", err)
	}
	return nil
}

// isValidVideoType checks if the file type is a valid video
func (s *VideoService) isValidVideoType(filename string) bool {
	ext := filepath.Ext(filename)
	for _, allowedType := range s.cfg.Storage.AllowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// saveUploadedFile saves an uploaded file to the specified path
func (s *VideoService) saveUploadedFile(file *multipart.FileHeader, filePath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy file content
	_, err = dst.ReadFrom(src)
	return err
}
