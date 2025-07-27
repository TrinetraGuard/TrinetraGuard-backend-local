package services

import (
	"database/sql"
	"encoding/json"
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

// FinderService handles person finding operations
type FinderService struct {
	db  *sql.DB
	cfg *config.Config
	log *zap.Logger
}

// NewFinderService creates a new finder service
func NewFinderService(db *sql.DB, cfg *config.Config) *FinderService {
	return &FinderService{
		db:  db,
		cfg: cfg,
	}
}

// UploadReferenceImage uploads and stores a reference image
func (s *FinderService) UploadReferenceImage(file *multipart.FileHeader, description string) (*models.ReferenceImage, error) {
	// Validate file type
	if !s.isValidImageType(file.Filename) {
		return nil, fmt.Errorf("invalid image file type")
	}

	// Validate file size
	if file.Size > s.cfg.Storage.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size")
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	filepath := filepath.Join(s.cfg.Storage.FinderDir, filename)

	// Ensure directory exists
	if err := os.MkdirAll(s.cfg.Storage.FinderDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create finder directory: %w", err)
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

	// Create reference image record in database
	refImage := &models.ReferenceImage{
		ID:               uuid.New().String(),
		Filename:         filename,
		OriginalFilename: file.Filename,
		FileSize:         file.Size,
		Description:      description,
		CreatedAt:        time.Now(),
	}

	query := `INSERT INTO reference_images (id, filename, original_filename, file_size, description, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	_, err = s.db.Exec(query, refImage.ID, refImage.Filename, refImage.OriginalFilename,
		refImage.FileSize, refImage.Description, refImage.CreatedAt)
	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to save reference image record: %w", err)
	}

	return refImage, nil
}

// ListReferenceImages retrieves all reference images
func (s *FinderService) ListReferenceImages() ([]models.ReferenceImage, error) {
	query := `SELECT id, filename, original_filename, file_size, description, created_at 
			  FROM reference_images ORDER BY created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query reference images: %w", err)
	}
	defer rows.Close()

	var images []models.ReferenceImage
	for rows.Next() {
		var image models.ReferenceImage
		err := rows.Scan(
			&image.ID, &image.Filename, &image.OriginalFilename,
			&image.FileSize, &image.Description, &image.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reference image: %w", err)
		}
		images = append(images, image)
	}

	if images == nil {
		images = []models.ReferenceImage{}
	}
	return images, nil
}

// DeleteReferenceImage deletes a reference image and its associated file
func (s *FinderService) DeleteReferenceImage(id string) error {
	// Get image info first
	query := `SELECT filename FROM reference_images WHERE id = ?`
	var filename string
	err := s.db.QueryRow(query, id).Scan(&filename)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("reference image not found")
		}
		return fmt.Errorf("failed to get reference image: %w", err)
	}

	// Delete from database
	query = `DELETE FROM reference_images WHERE id = ?`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete reference image from database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reference image not found")
	}

	// Delete file
	filepath := filepath.Join(s.cfg.Storage.FinderDir, filename)
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete reference image file: %w", err)
	}

	return nil
}

// SearchPerson starts a search job to find a person in videos
func (s *FinderService) SearchPerson(request *models.SearchPersonRequest) (*models.SearchJob, error) {
	// Check if reference image exists
	query := `SELECT id FROM reference_images WHERE id = ?`
	var id string
	err := s.db.QueryRow(query, request.ReferenceImageID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reference image not found")
		}
		return nil, fmt.Errorf("failed to check reference image existence: %w", err)
	}

	// Create new search job
	job := &models.SearchJob{
		ID:               uuid.New().String(),
		ReferenceImageID: request.ReferenceImageID,
		Status:           models.JobStatusPending,
		Progress:         0,
		CreatedAt:        time.Now(),
	}

	query = `INSERT INTO search_jobs (id, reference_image_id, status, progress, created_at) 
			 VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, job.ID, job.ReferenceImageID, job.Status, job.Progress, job.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create search job: %w", err)
	}

	// Start search in background
	go s.runSearch(job.ID, request.VideoIDs)

	return job, nil
}

// GetSearchStatus retrieves the status of a search job
func (s *FinderService) GetSearchStatus(searchID string) (*models.SearchJob, error) {
	query := `SELECT id, reference_image_id, status, progress, error_message, started_at, completed_at, created_at 
			  FROM search_jobs WHERE id = ?`

	var job models.SearchJob
	err := s.db.QueryRow(query, searchID).Scan(
		&job.ID, &job.ReferenceImageID, &job.Status, &job.Progress,
		&job.ErrorMessage, &job.StartedAt, &job.CompletedAt, &job.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("search job not found")
		}
		return nil, fmt.Errorf("failed to get search job: %w", err)
	}

	return &job, nil
}

// GetSearchResults retrieves the results of a completed search
func (s *FinderService) GetSearchResults(searchID string) ([]models.SearchResult, error) {
	query := `SELECT id, search_job_id, video_id, reference_image_id, matches, first_appearance, 
			  last_appearance, total_appearances, confidence, created_at
			  FROM search_results WHERE search_job_id = ?`

	rows, err := s.db.Query(query, searchID)
	if err != nil {
		return nil, fmt.Errorf("failed to query search results: %w", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var result models.SearchResult
		err := rows.Scan(
			&result.ID, &result.SearchJobID, &result.VideoID, &result.ReferenceImageID,
			&result.Matches, &result.FirstAppearance, &result.LastAppearance,
			&result.TotalAppearances, &result.Confidence, &result.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// runSearch runs the actual person search (placeholder implementation)
func (s *FinderService) runSearch(searchID string, videoIDs []string) {
	// Update job status to running
	s.updateSearchJobStatus(searchID, models.JobStatusRunning, 0, nil)

	// Get all videos if none specified
	if len(videoIDs) == 0 {
		query := `SELECT id FROM videos WHERE status = ?`
		rows, err := s.db.Query(query, models.VideoStatusAnalyzed)
		if err != nil {
			s.log.Error("Failed to get videos for search", zap.String("search_id", searchID), zap.Error(err))
			errorMsg := err.Error()
			s.updateSearchJobStatus(searchID, models.JobStatusFailed, 0, &errorMsg)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var videoID string
			if err := rows.Scan(&videoID); err != nil {
				s.log.Error("Failed to scan video ID", zap.Error(err))
				continue
			}
			videoIDs = append(videoIDs, videoID)
		}
	}

	// TODO: Implement actual person search
	// This would typically involve:
	// 1. Loading the reference image
	// 2. For each video, loading frames and running face detection
	// 3. Comparing detected faces with the reference image
	// 4. Storing matches with timestamps and confidence scores

	// Simulate search progress
	totalVideos := len(videoIDs)
	for i, videoID := range videoIDs {
		progress := (i + 1) * 100 / totalVideos
		s.updateSearchJobProgress(searchID, progress)

		// Create mock results for each video
		s.createMockSearchResults(searchID, videoID)

		time.Sleep(200 * time.Millisecond) // Simulate processing time
	}

	// Update job status to completed
	s.updateSearchJobStatus(searchID, models.JobStatusCompleted, 100, nil)
}

// updateSearchJobStatus updates the status of a search job
func (s *FinderService) updateSearchJobStatus(searchID, status string, progress int, errorMessage *string) {
	var query string
	var args []interface{}

	if status == models.JobStatusRunning {
		query = `UPDATE search_jobs SET status = ?, progress = ?, started_at = ? WHERE id = ?`
		args = []interface{}{status, progress, time.Now(), searchID}
	} else if status == models.JobStatusCompleted {
		query = `UPDATE search_jobs SET status = ?, progress = ?, completed_at = ? WHERE id = ?`
		args = []interface{}{status, progress, time.Now(), searchID}
	} else if status == models.JobStatusFailed {
		query = `UPDATE search_jobs SET status = ?, progress = ?, error_message = ?, completed_at = ? WHERE id = ?`
		args = []interface{}{status, progress, errorMessage, time.Now(), searchID}
	} else {
		query = `UPDATE search_jobs SET status = ?, progress = ? WHERE id = ?`
		args = []interface{}{status, progress, searchID}
	}

	_, err := s.db.Exec(query, args...)
	if err != nil {
		s.log.Error("Failed to update search job status", zap.String("search_id", searchID), zap.Error(err))
	}
}

// updateSearchJobProgress updates the progress of a search job
func (s *FinderService) updateSearchJobProgress(searchID string, progress int) {
	query := `UPDATE search_jobs SET progress = ? WHERE id = ?`
	_, err := s.db.Exec(query, progress, searchID)
	if err != nil {
		s.log.Error("Failed to update search job progress", zap.String("search_id", searchID), zap.Error(err))
	}
}

// createMockSearchResults creates mock search results for demonstration
func (s *FinderService) createMockSearchResults(searchID, videoID string) {
	// Get reference image ID from search job
	var referenceImageID string
	query := `SELECT reference_image_id FROM search_jobs WHERE id = ?`
	err := s.db.QueryRow(query, searchID).Scan(&referenceImageID)
	if err != nil {
		s.log.Error("Failed to get reference image ID for search", zap.String("search_id", searchID), zap.Error(err))
		return
	}

	// Create mock matches
	matches := []models.Match{
		{
			FrameNumber: 10,
			Timestamp:   10.0,
			BoundingBox: models.BoundingBox{X: 150, Y: 200, Width: 90, Height: 220},
			Confidence:  0.92,
		},
		{
			FrameNumber: 25,
			Timestamp:   25.0,
			BoundingBox: models.BoundingBox{X: 300, Y: 180, Width: 85, Height: 210},
			Confidence:  0.88,
		},
		{
			FrameNumber: 40,
			Timestamp:   40.0,
			BoundingBox: models.BoundingBox{X: 200, Y: 250, Width: 95, Height: 230},
			Confidence:  0.95,
		},
	}

	matchesJSON, _ := json.Marshal(matches)

	// Create search result
	result := &models.SearchResult{
		ID:               uuid.New().String(),
		SearchJobID:      searchID,
		VideoID:          videoID,
		ReferenceImageID: referenceImageID,
		Matches:          matchesJSON,
		FirstAppearance:  10.0,
		LastAppearance:   40.0,
		TotalAppearances: 3,
		Confidence:       0.92,
		CreatedAt:        time.Now(),
	}

	query = `INSERT INTO search_results (id, search_job_id, video_id, reference_image_id, matches, 
			 first_appearance, last_appearance, total_appearances, confidence, created_at) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = s.db.Exec(query, result.ID, result.SearchJobID, result.VideoID, result.ReferenceImageID,
		result.Matches, result.FirstAppearance, result.LastAppearance, result.TotalAppearances,
		result.Confidence, result.CreatedAt)
	if err != nil {
		s.log.Error("Failed to create search results", zap.String("search_id", searchID), zap.Error(err))
	}
}

// isValidImageType checks if the file type is a valid image
func (s *FinderService) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validTypes := []string{".jpg", ".jpeg", ".png", ".bmp", ".gif", ".tiff", ".webp"}
	for _, validType := range validTypes {
		if ext == validType {
			return true
		}
	}
	return false
}
