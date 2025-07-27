package services

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"video-analysis-service/internal/config"
	"video-analysis-service/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// VideoService handles video-related operations
type VideoService struct {
	db  *sql.DB
	cfg *config.Config
	log *zap.Logger
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
		return nil, fmt.Errorf("invalid video file type")
	}

	// Validate file size
	if file.Size > s.cfg.Storage.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size")
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	filepath := filepath.Join(s.cfg.Storage.VideosDir, filename)

	// Ensure directory exists
	if err := os.MkdirAll(s.cfg.Storage.VideosDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create videos directory: %w", err)
	}

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Create video record in database
	video := &models.Video{
		ID:               uuid.New().String(),
		Filename:         filename,
		OriginalFilename: file.Filename,
		FileSize:         file.Size,
		Status:           models.VideoStatusUploaded,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	query := `INSERT INTO videos (id, filename, original_filename, file_size, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = s.db.Exec(query, video.ID, video.Filename, video.OriginalFilename,
		video.FileSize, video.Status, video.CreatedAt, video.UpdatedAt)
	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to save video record: %w", err)
	}

	return video, nil
}

// GetVideo retrieves a video by ID
func (s *VideoService) GetVideo(id string) (*models.Video, error) {
	query := `SELECT id, filename, original_filename, file_size, duration, frame_count, 
			  width, height, status, created_at, updated_at FROM videos WHERE id = ?`

	var video models.Video
	var duration sql.NullFloat64
	var frameCount, width, height sql.NullInt64

	err := s.db.QueryRow(query, id).Scan(
		&video.ID, &video.Filename, &video.OriginalFilename, &video.FileSize,
		&duration, &frameCount, &width, &height,
		&video.Status, &video.CreatedAt, &video.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("video not found")
		}
		return nil, fmt.Errorf("failed to get video: %w", err)
	}
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

// ListVideos retrieves all videos with optional pagination
func (s *VideoService) ListVideos(limit, offset int) ([]models.Video, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM videos`
	err := s.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count videos: %w", err)
	}

	// Get videos
	query := `SELECT id, filename, original_filename, file_size, duration, frame_count, 
			  width, height, status, created_at, updated_at FROM videos 
			  ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var video models.Video
		var duration sql.NullFloat64
		var frameCount, width, height sql.NullInt64
		err := rows.Scan(
			&video.ID, &video.Filename, &video.OriginalFilename, &video.FileSize,
			&duration, &frameCount, &width, &height,
			&video.Status, &video.CreatedAt, &video.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan video: %w", err)
		}
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
		videos = []models.Video{}
	}
	return videos, total, nil
}

// DeleteVideo deletes a video and its associated file
func (s *VideoService) DeleteVideo(id string) error {
	// Get video info first
	video, err := s.GetVideo(id)
	if err != nil {
		return err
	}

	// Delete from database
	query := `DELETE FROM videos WHERE id = ?`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete video from database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("video not found")
	}

	// Delete file
	filepath := filepath.Join(s.cfg.Storage.VideosDir, video.Filename)
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete video file: %w", err)
	}

	return nil
}

// GetVideoFilePath returns the full path to a video file
func (s *VideoService) GetVideoFilePath(filename string) string {
	return filepath.Join(s.cfg.Storage.VideosDir, filename)
}

// UpdateVideoStatus updates the status of a video
func (s *VideoService) UpdateVideoStatus(id, status string) error {
	query := `UPDATE videos SET status = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update video status: %w", err)
	}
	return nil
}

// UpdateVideoMetadata updates video metadata (duration, dimensions, etc.)
func (s *VideoService) UpdateVideoMetadata(id string, duration float64, frameCount, width, height int) error {
	query := `UPDATE videos SET duration = ?, frame_count = ?, width = ?, height = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, duration, frameCount, width, height, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update video metadata: %w", err)
	}
	return nil
}

// isValidVideoType checks if the file type is allowed
func (s *VideoService) isValidVideoType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range s.cfg.Storage.AllowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}
