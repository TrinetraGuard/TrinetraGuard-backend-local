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

// FinderService handles person finding operations
type FinderService struct {
	db  *sql.DB
	cfg *config.Config
}

// NewFinderService creates a new finder service
func NewFinderService(db *sql.DB, cfg *config.Config) *FinderService {
	return &FinderService{
		db:  db,
		cfg: cfg,
	}
}

// UploadReferenceImage uploads a reference image for person finding
func (s *FinderService) UploadReferenceImage(file *multipart.FileHeader) (*models.ReferenceImage, error) {
	// Validate file type
	if !s.isValidImageType(file.Filename) {
		return nil, fmt.Errorf("unsupported image format")
	}

	// Generate unique filename
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	filePath := filepath.Join(s.cfg.Storage.FinderDir, filename)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(s.cfg.Storage.FinderDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	// Save file
	if err := s.saveUploadedFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	// Create reference image record
	refImage := &models.ReferenceImage{
		ID:           uuid.New().String(),
		Filename:     filename,
		OriginalName: file.Filename,
		FilePath:     filePath,
		FileSize:     fileInfo.Size(),
		UploadedAt:   time.Now(),
	}

	// Save to database
	_, err = s.db.Exec(`
		INSERT INTO reference_images (id, filename, original_name, file_path, file_size, uploaded_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, refImage.ID, refImage.Filename, refImage.OriginalName, refImage.FilePath, refImage.FileSize, refImage.UploadedAt)

	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save reference image: %v", err)
	}

	return refImage, nil
}

// ListReferenceImages lists all reference images
func (s *FinderService) ListReferenceImages() ([]models.ReferenceImage, error) {
	rows, err := s.db.Query(`
		SELECT id, filename, original_name, file_path, file_size, uploaded_at
		FROM reference_images ORDER BY uploaded_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query reference images: %v", err)
	}
	defer rows.Close()

	var images []models.ReferenceImage
	for rows.Next() {
		var image models.ReferenceImage
		err := rows.Scan(&image.ID, &image.Filename, &image.OriginalName, &image.FilePath, &image.FileSize, &image.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reference image: %v", err)
		}
		images = append(images, image)
	}

	if images == nil {
		return []models.ReferenceImage{}, nil
	}

	return images, nil
}

// DeleteReferenceImage deletes a reference image
func (s *FinderService) DeleteReferenceImage(imageID string) error {
	// Get file path before deleting
	var filePath string
	err := s.db.QueryRow("SELECT file_path FROM reference_images WHERE id = ?", imageID).Scan(&filePath)
	if err != nil {
		return fmt.Errorf("reference image not found: %v", err)
	}

	// Delete from database
	_, err = s.db.Exec("DELETE FROM reference_images WHERE id = ?", imageID)
	if err != nil {
		return fmt.Errorf("failed to delete reference image: %v", err)
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to delete file %s: %v\n", filePath, err)
	}

	return nil
}

// SearchPerson starts a person search job
func (s *FinderService) SearchPerson(referenceImageID string) (*models.SearchJob, error) {
	// Check if reference image exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM reference_images WHERE id = ?)", referenceImageID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check reference image existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("reference image not found: %s", referenceImageID)
	}

	// Create search job
	jobID := uuid.New().String()
	now := time.Now()

	_, err = s.db.Exec(`
		INSERT INTO search_jobs (id, reference_image_id, status, progress, started_at)
		VALUES (?, ?, ?, ?, ?)
	`, jobID, referenceImageID, models.JobStatusPending, 0, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create search job: %v", err)
	}

	// Start search in background
	go s.runSearch(jobID, referenceImageID)

	return &models.SearchJob{
		ID:               jobID,
		ReferenceImageID: referenceImageID,
		Status:           models.JobStatusPending,
		Progress:         0,
		StartedAt:        now,
	}, nil
}

// GetSearchStatus gets the status of a search job
func (s *FinderService) GetSearchStatus(jobID string) (*models.SearchJob, error) {
	var job models.SearchJob
	var completedAt sql.NullTime
	var errorMsg sql.NullString

	err := s.db.QueryRow(`
		SELECT id, reference_image_id, status, progress, started_at, completed_at, error
		FROM search_jobs WHERE id = ?
	`, jobID).Scan(&job.ID, &job.ReferenceImageID, &job.Status, &job.Progress, &job.StartedAt, &completedAt, &errorMsg)

	if err != nil {
		return nil, fmt.Errorf("search job not found: %v", err)
	}

	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}
	if errorMsg.Valid {
		job.Error = &errorMsg.String
	}

	return &job, nil
}

// GetSearchResults gets the results of a completed search
func (s *FinderService) GetSearchResults(jobID string) (*models.SearchResult, error) {
	var result models.SearchResult

	err := s.db.QueryRow(`
		SELECT id, search_job_id, video_id, created_at
		FROM search_results WHERE search_job_id = ? ORDER BY created_at DESC LIMIT 1
	`, jobID).Scan(&result.ID, &result.SearchJobID, &result.VideoID, &result.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("search results not found: %v", err)
	}

	// Get matches for this search result
	rows, err := s.db.Query(`
		SELECT person_id, first_frame, last_frame, first_time, last_time, confidence, x, y, width, height
		FROM search_matches WHERE search_result_id = ? ORDER BY first_time
	`, result.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get search matches: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var match models.Match
		err := rows.Scan(&match.PersonID, &match.FirstFrame, &match.LastFrame, &match.FirstTime, &match.LastTime, &match.Confidence, &match.BoundingBox.X, &match.BoundingBox.Y, &match.BoundingBox.Width, &match.BoundingBox.Height)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search match: %v", err)
		}
		result.Matches = append(result.Matches, match)
	}

	return &result, nil
}

// isValidImageType checks if the file is a valid image type
func (s *FinderService) isValidImageType(filename string) bool {
	ext := filepath.Ext(filename)
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// saveUploadedFile saves an uploaded file to the specified path
func (s *FinderService) saveUploadedFile(file *multipart.FileHeader, filePath string) error {
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

// runSearch simulates the search process
func (s *FinderService) runSearch(jobID, referenceImageID string) {
	// Update job status to running
	_, err := s.db.Exec("UPDATE search_jobs SET status = ?, progress = ? WHERE id = ?", models.JobStatusRunning, 10, jobID)
	if err != nil {
		return
	}

	// Simulate search progress
	for progress := 20; progress <= 90; progress += 10 {
		time.Sleep(500 * time.Millisecond)
		_, err := s.db.Exec("UPDATE search_jobs SET progress = ? WHERE id = ?", progress, jobID)
		if err != nil {
			return
		}
	}

	// Create mock search results
	err = s.createMockSearchResults(jobID, referenceImageID)
	if err != nil {
		errorMsg := err.Error()
		_, err = s.db.Exec("UPDATE search_jobs SET status = ?, error = ?, completed_at = ? WHERE id = ?", models.JobStatusFailed, errorMsg, time.Now(), jobID)
		return
	}

	// Update job status to completed
	_, err = s.db.Exec("UPDATE search_jobs SET status = ?, progress = ?, completed_at = ? WHERE id = ?", models.JobStatusCompleted, 100, time.Now(), jobID)
}

// createMockSearchResults creates mock search results
func (s *FinderService) createMockSearchResults(jobID, referenceImageID string) error {
	// Get a video ID for mock results
	var videoID string
	err := s.db.QueryRow("SELECT id FROM videos LIMIT 1").Scan(&videoID)
	if err != nil {
		return fmt.Errorf("no videos available for search: %v", err)
	}

	// Create search result
	resultID := uuid.New().String()
	now := time.Now()

	_, err = s.db.Exec(`
		INSERT INTO search_results (id, search_job_id, video_id, created_at)
		VALUES (?, ?, ?, ?)
	`, resultID, jobID, videoID, now)

	if err != nil {
		return fmt.Errorf("failed to create search result: %v", err)
	}

	// Create mock matches
	matches := []struct {
		personID   string
		firstTime  float64
		lastTime   float64
		confidence float64
	}{
		{"person_1", 10.5, 25.3, 0.85},
		{"person_2", 45.2, 67.8, 0.92},
		{"person_3", 120.1, 145.6, 0.78},
	}

	for _, match := range matches {
		_, err := s.db.Exec(`
			INSERT INTO search_matches (id, search_result_id, person_id, first_frame, last_frame, first_time, last_time, confidence, x, y, width, height)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, uuid.New().String(), resultID, match.personID, int(match.firstTime*30), int(match.lastTime*30), match.firstTime, match.lastTime, match.confidence, 0.1, 0.2, 0.3, 0.4)

		if err != nil {
			return fmt.Errorf("failed to create search match: %v", err)
		}
	}

	return nil
}
