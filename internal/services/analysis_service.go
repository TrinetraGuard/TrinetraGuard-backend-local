package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"video-analysis-service/internal/config"
	"video-analysis-service/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AnalysisService handles video analysis operations
type AnalysisService struct {
	db  *sql.DB
	cfg *config.Config
	log *zap.Logger
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(db *sql.DB, cfg *config.Config) *AnalysisService {
	return &AnalysisService{
		db:  db,
		cfg: cfg,
	}
}

// StartAnalysis starts a new analysis job for a video
func (s *AnalysisService) StartAnalysis(videoID string) (*models.AnalysisJob, error) {
	// Check if video exists
	query := `SELECT id FROM videos WHERE id = ?`
	var id string
	err := s.db.QueryRow(query, videoID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("video not found")
		}
		return nil, fmt.Errorf("failed to check video existence: %w", err)
	}

	// Check if there's already a running job for this video
	query = `SELECT id FROM analysis_jobs WHERE video_id = ? AND status IN (?, ?)`
	err = s.db.QueryRow(query, videoID, models.JobStatusPending, models.JobStatusRunning).Scan(&id)
	if err == nil {
		return nil, fmt.Errorf("analysis job already exists for this video")
	}

	// Create new analysis job
	job := &models.AnalysisJob{
		ID:        uuid.New().String(),
		VideoID:   videoID,
		Status:    models.JobStatusPending,
		Progress:  0,
		CreatedAt: time.Now(),
	}

	query = `INSERT INTO analysis_jobs (id, video_id, status, progress, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, job.ID, job.VideoID, job.Status, job.Progress, job.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create analysis job: %w", err)
	}

	// Update video status
	query = `UPDATE videos SET status = ?, updated_at = ? WHERE id = ?`
	_, err = s.db.Exec(query, models.VideoStatusAnalyzing, time.Now(), videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to update video status: %w", err)
	}

	// Start analysis in background
	go s.runAnalysis(job.ID)

	return job, nil
}

// GetAnalysisStatus retrieves the status of an analysis job
func (s *AnalysisService) GetAnalysisStatus(jobID string) (*models.AnalysisJob, error) {
	query := `SELECT id, video_id, status, progress, error_message, started_at, completed_at, created_at 
			  FROM analysis_jobs WHERE id = ?`

	var job models.AnalysisJob
	err := s.db.QueryRow(query, jobID).Scan(
		&job.ID, &job.VideoID, &job.Status, &job.Progress,
		&job.ErrorMessage, &job.StartedAt, &job.CompletedAt, &job.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analysis job not found")
		}
		return nil, fmt.Errorf("failed to get analysis job: %w", err)
	}

	return &job, nil
}

// GetAnalysisResults retrieves the results of a completed analysis
func (s *AnalysisService) GetAnalysisResults(videoID string) (*models.AnalysisResult, error) {
	query := `SELECT ar.id, ar.job_id, ar.video_id, ar.total_frames, ar.total_people, ar.unique_people,
			  ar.people_per_frame, ar.tracking_data, ar.created_at
			  FROM analysis_results ar
			  JOIN analysis_jobs aj ON ar.job_id = aj.id
			  WHERE ar.video_id = ? AND aj.status = ?
			  ORDER BY ar.created_at DESC LIMIT 1`

	var result models.AnalysisResult
	err := s.db.QueryRow(query, videoID, models.JobStatusCompleted).Scan(
		&result.ID, &result.JobID, &result.VideoID, &result.TotalFrames,
		&result.TotalPeople, &result.UniquePeople, &result.PeoplePerFrame,
		&result.TrackingData, &result.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analysis results not found")
		}
		return nil, fmt.Errorf("failed to get analysis results: %w", err)
	}

	return &result, nil
}

// StartBatchAnalysis starts analysis for multiple videos
func (s *AnalysisService) StartBatchAnalysis(videoIDs []string) ([]models.AnalysisJob, error) {
	var jobs []models.AnalysisJob

	for _, videoID := range videoIDs {
		job, err := s.StartAnalysis(videoID)
		if err != nil {
			// Log error but continue with other videos
			s.log.Error("Failed to start analysis for video",
				zap.String("video_id", videoID), zap.Error(err))
			continue
		}
		jobs = append(jobs, *job)
	}

	return jobs, nil
}

// runAnalysis runs the actual video analysis (placeholder implementation)
func (s *AnalysisService) runAnalysis(jobID string) {
	// Update job status to running
	s.updateJobStatus(jobID, models.JobStatusRunning, 0, nil)

	// TODO: Implement actual video analysis
	// This would typically involve:
	// 1. Loading the video file
	// 2. Extracting frames at specified intervals
	// 3. Running person detection on each frame
	// 4. Tracking individuals across frames
	// 5. Storing results

	// Simulate analysis progress
	for i := 0; i <= 100; i += 10 {
		time.Sleep(100 * time.Millisecond) // Simulate processing time
		s.updateJobProgress(jobID, i)
	}

	// Create mock results
	s.createMockResults(jobID)

	// Update job status to completed
	s.updateJobStatus(jobID, models.JobStatusCompleted, 100, nil)
}

// updateJobStatus updates the status of an analysis job
func (s *AnalysisService) updateJobStatus(jobID, status string, progress int, errorMessage *string) {
	var query string
	var args []interface{}

	if status == models.JobStatusRunning {
		query = `UPDATE analysis_jobs SET status = ?, progress = ?, started_at = ? WHERE id = ?`
		args = []interface{}{status, progress, time.Now(), jobID}
	} else if status == models.JobStatusCompleted {
		query = `UPDATE analysis_jobs SET status = ?, progress = ?, completed_at = ? WHERE id = ?`
		args = []interface{}{status, progress, time.Now(), jobID}
	} else if status == models.JobStatusFailed {
		query = `UPDATE analysis_jobs SET status = ?, progress = ?, error_message = ?, completed_at = ? WHERE id = ?`
		args = []interface{}{status, progress, errorMessage, time.Now(), jobID}
	} else {
		query = `UPDATE analysis_jobs SET status = ?, progress = ? WHERE id = ?`
		args = []interface{}{status, progress, jobID}
	}

	_, err := s.db.Exec(query, args...)
	if err != nil {
		s.log.Error("Failed to update job status", zap.String("job_id", jobID), zap.Error(err))
	}
}

// updateJobProgress updates the progress of an analysis job
func (s *AnalysisService) updateJobProgress(jobID string, progress int) {
	query := `UPDATE analysis_jobs SET progress = ? WHERE id = ?`
	_, err := s.db.Exec(query, progress, jobID)
	if err != nil {
		s.log.Error("Failed to update job progress", zap.String("job_id", jobID), zap.Error(err))
	}
}

// createMockResults creates mock analysis results for demonstration
func (s *AnalysisService) createMockResults(jobID string) {
	// Get video ID from job
	var videoID string
	query := `SELECT video_id FROM analysis_jobs WHERE id = ?`
	err := s.db.QueryRow(query, jobID).Scan(&videoID)
	if err != nil {
		s.log.Error("Failed to get video ID for job", zap.String("job_id", jobID), zap.Error(err))
		return
	}

	// Create mock people per frame data
	peoplePerFrame := []models.PeoplePerFrame{
		{FrameNumber: 1, Timestamp: 0.0, Count: 3, Confidence: 0.85},
		{FrameNumber: 2, Timestamp: 1.0, Count: 4, Confidence: 0.87},
		{FrameNumber: 3, Timestamp: 2.0, Count: 2, Confidence: 0.82},
	}

	peoplePerFrameJSON, _ := json.Marshal(peoplePerFrame)

	// Create mock tracking data
	trackingData := []models.TrackingData{
		{
			PersonID:    "person_1",
			FrameNumber: 1,
			Timestamp:   0.0,
			BoundingBox: models.BoundingBox{X: 100, Y: 150, Width: 80, Height: 200},
			Confidence:  0.85,
		},
		{
			PersonID:    "person_2",
			FrameNumber: 1,
			Timestamp:   0.0,
			BoundingBox: models.BoundingBox{X: 300, Y: 200, Width: 70, Height: 180},
			Confidence:  0.87,
		},
	}

	trackingDataJSON, _ := json.Marshal(trackingData)

	// Create analysis result
	result := &models.AnalysisResult{
		ID:             uuid.New().String(),
		JobID:          jobID,
		VideoID:        videoID,
		TotalFrames:    3,
		TotalPeople:    9,
		UniquePeople:   2,
		PeoplePerFrame: peoplePerFrameJSON,
		TrackingData:   trackingDataJSON,
		CreatedAt:      time.Now(),
	}

	query = `INSERT INTO analysis_results (id, job_id, video_id, total_frames, total_people, 
			 unique_people, people_per_frame, tracking_data, created_at) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = s.db.Exec(query, result.ID, result.JobID, result.VideoID, result.TotalFrames,
		result.TotalPeople, result.UniquePeople, result.PeoplePerFrame, result.TrackingData, result.CreatedAt)
	if err != nil {
		s.log.Error("Failed to create analysis results", zap.String("job_id", jobID), zap.Error(err))
		return
	}

	// Create mock persons with faces
	s.createMockPersons(videoID)

	// Update video status to analyzed
	query = `UPDATE videos SET status = ?, updated_at = ? WHERE id = ?`
	_, err = s.db.Exec(query, models.VideoStatusAnalyzed, time.Now(), videoID)
	if err != nil {
		s.log.Error("Failed to update video status", zap.String("video_id", videoID), zap.Error(err))
	}
}

// createMockPersons creates mock person and face data
func (s *AnalysisService) createMockPersons(videoID string) {
	// Create Person 1
	person1ID := uuid.New().String()
	query := `INSERT INTO persons (id, video_id, person_number, first_frame, last_frame, first_time, last_time, total_frames, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, person1ID, videoID, 1, 1, 3, 0.0, 2.0, 3, time.Now())
	if err != nil {
		s.log.Error("Failed to create person 1", zap.String("video_id", videoID), zap.Error(err))
		return
	}

	// Create Person 2
	person2ID := uuid.New().String()
	_, err = s.db.Exec(query, person2ID, videoID, 2, 1, 3, 0.0, 2.0, 3, time.Now())
	if err != nil {
		s.log.Error("Failed to create person 2", zap.String("video_id", videoID), zap.Error(err))
		return
	}

	// Create mock face images (base64 encoded placeholder images)
	// Simple face-like SVG converted to base64
	face1Image := "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgdmlld0JveD0iMCAwIDEwMCAxMDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIxMDAiIGhlaWdodD0iMTAwIiBmaWxsPSIjRkZGRkZGIi8+CjxyZWN0IHg9IjIwIiB5PSIyMCIgd2lkdGg9IjEwIiBoZWlnaHQ9IjEwIiBmaWxsPSIjMzMzMzMzIi8+CjxyZWN0IHg9IjcwIiB5PSIyMCIgd2lkdGg9IjEwIiBoZWlnaHQ9IjEwIiBmaWxsPSIjMzMzMzMzIi8+CjxyZWN0IHg9IjQwIiB5PSI0MCIgd2lkdGg9IjIwIiBoZWlnaHQ9IjEwIiBmaWxsPSIjMzMzMzMzIi8+CjxyZWN0IHg9IjMwIiB5PSI2MCIgd2lkdGg9IjQwIiBoZWlnaHQ9IjEwIiBmaWxsPSIjMzMzMzMzIi8+Cjwvc3ZnPgo="
	face2Image := "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgdmlld0JveD0iMCAwIDEwMCAxMDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIxMDAiIGhlaWdodD0iMTAwIiBmaWxsPSIjRkZGRkZGIi8+CjxyZWN0IHg9IjE1IiB5PSIyMCIgd2lkdGg9IjEwIiBoZWlnaHQ9IjEwIiBmaWxsPSIjMzMzMzMzIi8+CjxyZWN0IHg9Ijc1IiB5PSIyMCIgd2lkdGg9IjEwIiBoZWlnaHQ9IjEwIiBmaWxsPSIjMzMzMzMzIi8+CjxyZWN0IHg9IjQwIiB5PSI0MCIgd2lkdGg9IjIwIiBoZWlnaHQ9IjEwIiBmaWxsPSIjMzMzMzMzIi8+CjxyZWN0IHg9IjI1IiB5PSI2MCIgd2lkdGg9IjUwIiBoZWlnaHQ9IjEwIiBmaWxsPSIjMzMzMzMzIi8+Cjwvc3ZnPgo="

	// Create face for Person 1
	face1ID := uuid.New().String()
	query = `INSERT INTO person_faces (id, person_id, video_id, frame_number, timestamp, 
			 bounding_box_x, bounding_box_y, bounding_box_width, bounding_box_height, confidence, face_image, created_at) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, face1ID, person1ID, videoID, 2, 1.0, 100.0, 150.0, 80.0, 200.0, 0.85, face1Image, time.Now())
	if err != nil {
		s.log.Error("Failed to create face 1", zap.String("person_id", person1ID), zap.Error(err))
	}

	// Create face for Person 2
	face2ID := uuid.New().String()
	_, err = s.db.Exec(query, face2ID, person2ID, videoID, 2, 1.0, 300.0, 200.0, 70.0, 180.0, 0.87, face2Image, time.Now())
	if err != nil {
		s.log.Error("Failed to create face 2", zap.String("person_id", person2ID), zap.Error(err))
	}
}

// GetEnhancedAnalysisResults retrieves analysis results with person and face data
func (s *AnalysisService) GetEnhancedAnalysisResults(videoID string) (*models.EnhancedAnalysisResult, error) {
	// Get basic analysis results
	query := `SELECT id, job_id, video_id, total_frames, total_people, unique_people, 
			  people_per_frame, tracking_data, created_at FROM analysis_results 
			  WHERE video_id = ? ORDER BY created_at DESC LIMIT 1`

	var result models.EnhancedAnalysisResult
	err := s.db.QueryRow(query, videoID).Scan(
		&result.ID, &result.JobID, &result.VideoID, &result.TotalFrames,
		&result.TotalPeople, &result.UniquePeople, &result.PeoplePerFrame,
		&result.TrackingData, &result.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analysis results not found")
		}
		return nil, fmt.Errorf("failed to get analysis results: %w", err)
	}

	// Get persons with their best faces
	persons, err := s.getPersonsWithFaces(videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get persons: %w", err)
	}

	result.Persons = persons
	return &result, nil
}

// getPersonsWithFaces retrieves all persons for a video with their best face
func (s *AnalysisService) getPersonsWithFaces(videoID string) ([]models.Person, error) {
	query := `SELECT id, video_id, person_number, first_frame, last_frame, first_time, last_time, total_frames, created_at 
			  FROM persons WHERE video_id = ? ORDER BY person_number`

	rows, err := s.db.Query(query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to query persons: %w", err)
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var person models.Person
		err := rows.Scan(
			&person.ID, &person.VideoID, &person.PersonNumber, &person.FirstFrame, &person.LastFrame,
			&person.FirstTime, &person.LastTime, &person.TotalFrames, &person.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan person: %w", err)
		}

		// Get best face for this person (highest confidence)
		face, err := s.getBestFaceForPerson(person.ID)
		if err == nil {
			person.BestFace = *face
		}

		persons = append(persons, person)
	}

	return persons, nil
}

// getBestFaceForPerson retrieves the best face (highest confidence) for a person
func (s *AnalysisService) getBestFaceForPerson(personID string) (*models.PersonFace, error) {
	query := `SELECT id, person_id, video_id, frame_number, timestamp, 
			  bounding_box_x, bounding_box_y, bounding_box_width, bounding_box_height, confidence, face_image, created_at 
			  FROM person_faces WHERE person_id = ? ORDER BY confidence DESC LIMIT 1`

	var face models.PersonFace
	err := s.db.QueryRow(query, personID).Scan(
		&face.ID, &face.PersonID, &face.VideoID, &face.FrameNumber, &face.Timestamp,
		&face.BoundingBox.X, &face.BoundingBox.Y, &face.BoundingBox.Width, &face.BoundingBox.Height,
		&face.Confidence, &face.FaceImage, &face.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get best face: %w", err)
	}

	return &face, nil
}
