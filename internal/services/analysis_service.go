package services

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"video-analysis-service/internal/config"
	"video-analysis-service/internal/models"

	"github.com/google/uuid"
)

// AnalysisService handles video analysis operations
type AnalysisService struct {
	db  *sql.DB
	cfg *config.Config
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
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM videos WHERE id = ?)", videoID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check video existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("video not found: %s", videoID)
	}

	// Check if analysis job already exists
	var jobID string
	err = s.db.QueryRow("SELECT id FROM analysis_jobs WHERE video_id = ? AND status IN ('pending', 'running')", videoID).Scan(&jobID)
	if err == nil {
		// Job already exists, return it
		return s.GetAnalysisJob(jobID)
	}

	// Create new analysis job
	jobID = uuid.New().String()
	now := time.Now()

	_, err = s.db.Exec(`
		INSERT INTO analysis_jobs (id, video_id, status, progress, started_at)
		VALUES (?, ?, ?, ?, ?)
	`, jobID, videoID, models.JobStatusPending, 0, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create analysis job: %v", err)
	}

	// Start analysis in background
	go s.runAnalysis(jobID, videoID)

	return &models.AnalysisJob{
		ID:        jobID,
		VideoID:   videoID,
		Status:    models.JobStatusPending,
		Progress:  0,
		StartedAt: now,
	}, nil
}

// GetAnalysisJob retrieves an analysis job by ID
func (s *AnalysisService) GetAnalysisJob(jobID string) (*models.AnalysisJob, error) {
	var job models.AnalysisJob
	var completedAt sql.NullTime
	var errorMsg sql.NullString

	err := s.db.QueryRow(`
		SELECT id, video_id, status, progress, started_at, completed_at, error
		FROM analysis_jobs WHERE id = ?
	`, jobID).Scan(&job.ID, &job.VideoID, &job.Status, &job.Progress, &job.StartedAt, &completedAt, &errorMsg)

	if err != nil {
		return nil, fmt.Errorf("analysis job not found: %v", err)
	}

	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}
	if errorMsg.Valid {
		job.Error = &errorMsg.String
	}

	return &job, nil
}

// GetAnalysisStatus gets the status of an analysis job
func (s *AnalysisService) GetAnalysisStatus(videoID string) (*models.AnalysisJob, error) {
	var job models.AnalysisJob
	var completedAt sql.NullTime
	var errorMsg sql.NullString

	err := s.db.QueryRow(`
		SELECT id, video_id, status, progress, started_at, completed_at, error
		FROM analysis_jobs WHERE video_id = ? ORDER BY started_at DESC LIMIT 1
	`, videoID).Scan(&job.ID, &job.VideoID, &job.Status, &job.Progress, &job.StartedAt, &completedAt, &errorMsg)

	if err != nil {
		return nil, fmt.Errorf("analysis job not found: %v", err)
	}

	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}
	if errorMsg.Valid {
		job.Error = &errorMsg.String
	}

	return &job, nil
}

// GetAnalysisResults gets the analysis results for a video
func (s *AnalysisService) GetAnalysisResults(videoID string) (*models.AnalysisResult, error) {
	var result models.AnalysisResult

	err := s.db.QueryRow(`
		SELECT id, video_id, analysis_job_id, total_frames, total_people, unique_people, created_at
		FROM analysis_results WHERE video_id = ? ORDER BY created_at DESC LIMIT 1
	`, videoID).Scan(&result.ID, &result.VideoID, &result.AnalysisJobID, &result.TotalFrames, &result.TotalPeople, &result.UniquePeople, &result.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("analysis results not found: %v", err)
	}

	// Get people per frame data
	rows, err := s.db.Query(`
		SELECT frame_number, timestamp, count
		FROM people_per_frame WHERE analysis_result_id = ? ORDER BY frame_number
	`, result.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get people per frame data: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ppf models.PeoplePerFrame
		err := rows.Scan(&ppf.FrameNumber, &ppf.Timestamp, &ppf.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan people per frame: %v", err)
		}
		result.PeoplePerFrame = append(result.PeoplePerFrame, ppf)
	}

	// Get tracking data
	rows, err = s.db.Query(`
		SELECT person_id, frame_number, timestamp, x, y, width, height, confidence
		FROM tracking_data WHERE analysis_result_id = ? ORDER BY frame_number, person_id
	`, result.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking data: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var td models.TrackingData
		err := rows.Scan(&td.PersonID, &td.FrameNumber, &td.Timestamp, &td.BoundingBox.X, &td.BoundingBox.Y, &td.BoundingBox.Width, &td.BoundingBox.Height, &td.Confidence)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking data: %v", err)
		}
		result.TrackingData = append(result.TrackingData, td)
	}

	return &result, nil
}

// GetEnhancedAnalysisResults gets enhanced analysis results with person data
func (s *AnalysisService) GetEnhancedAnalysisResults(videoID string) (*models.EnhancedAnalysisResult, error) {
	// Get basic analysis results
	basicResult, err := s.GetAnalysisResults(videoID)
	if err != nil {
		return nil, err
	}

	enhancedResult := &models.EnhancedAnalysisResult{
		AnalysisResult: basicResult,
		Persons:        []models.Person{},
	}

	// Get persons with their faces
	persons, err := s.getPersonsWithFaces(videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get persons: %v", err)
	}

	enhancedResult.Persons = persons
	return enhancedResult, nil
}

// getPersonsWithFaces retrieves all persons for a video with their faces
func (s *AnalysisService) getPersonsWithFaces(videoID string) ([]models.Person, error) {
	rows, err := s.db.Query(`
		SELECT id, video_id, person_id, first_seen, last_seen, total_frames, created_at
		FROM persons WHERE video_id = ? ORDER BY first_seen
	`, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.ID, &person.VideoID, &person.PersonID, &person.FirstSeen, &person.LastSeen, &person.TotalFrames, &person.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Get best face for this person
		bestFace, err := s.getBestFaceForPerson(person.ID)
		if err == nil {
			person.BestFace = bestFace
		}

		// Get all faces for this person
		faces, err := s.getFacesForPerson(person.ID)
		if err == nil {
			person.Faces = faces
		}

		persons = append(persons, person)
	}

	return persons, nil
}

// getBestFaceForPerson gets the best face image for a person
func (s *AnalysisService) getBestFaceForPerson(personID string) (*models.PersonFace, error) {
	var face models.PersonFace
	err := s.db.QueryRow(`
		SELECT id, person_id, video_id, frame_number, timestamp, face_image, x, y, width, height, confidence, is_best_face, created_at
		FROM person_faces WHERE person_id = ? AND is_best_face = TRUE LIMIT 1
	`, personID).Scan(&face.ID, &face.PersonID, &face.VideoID, &face.FrameNumber, &face.Timestamp, &face.FaceImage, &face.BoundingBox.X, &face.BoundingBox.Y, &face.BoundingBox.Width, &face.BoundingBox.Height, &face.Confidence, &face.IsBestFace, &face.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &face, nil
}

// getFacesForPerson gets all face images for a person
func (s *AnalysisService) getFacesForPerson(personID string) ([]models.PersonFace, error) {
	rows, err := s.db.Query(`
		SELECT id, person_id, video_id, frame_number, timestamp, face_image, x, y, width, height, confidence, is_best_face, created_at
		FROM person_faces WHERE person_id = ? ORDER BY timestamp
	`, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faces []models.PersonFace
	for rows.Next() {
		var face models.PersonFace
		err := rows.Scan(&face.ID, &face.PersonID, &face.VideoID, &face.FrameNumber, &face.Timestamp, &face.FaceImage, &face.BoundingBox.X, &face.BoundingBox.Y, &face.BoundingBox.Width, &face.BoundingBox.Height, &face.Confidence, &face.IsBestFace, &face.CreatedAt)
		if err != nil {
			return nil, err
		}
		faces = append(faces, face)
	}

	return faces, nil
}

// StartBatchAnalysis starts analysis for multiple videos
func (s *AnalysisService) StartBatchAnalysis(videoIDs []string) ([]*models.AnalysisJob, error) {
	var jobs []*models.AnalysisJob

	for _, videoID := range videoIDs {
		job, err := s.StartAnalysis(videoID)
		if err != nil {
			return nil, fmt.Errorf("failed to start analysis for video %s: %v", videoID, err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// runAnalysis simulates the analysis process
func (s *AnalysisService) runAnalysis(jobID, videoID string) {
	// Update job status to running
	_, err := s.db.Exec("UPDATE analysis_jobs SET status = ?, progress = ? WHERE id = ?", models.JobStatusRunning, 10, jobID)
	if err != nil {
		return
	}

	// Simulate analysis progress
	for progress := 20; progress <= 90; progress += 10 {
		time.Sleep(500 * time.Millisecond)
		_, err := s.db.Exec("UPDATE analysis_jobs SET progress = ? WHERE id = ?", progress, jobID)
		if err != nil {
			return
		}
	}

	// Create mock results
	err = s.createMockResults(jobID, videoID)
	if err != nil {
		errorMsg := err.Error()
		_, err = s.db.Exec("UPDATE analysis_jobs SET status = ?, error = ?, completed_at = ? WHERE id = ?", models.JobStatusFailed, errorMsg, time.Now(), jobID)
		return
	}

	// Update job status to completed
	_, err = s.db.Exec("UPDATE analysis_jobs SET status = ?, progress = ?, completed_at = ? WHERE id = ?", models.JobStatusCompleted, 100, time.Now(), jobID)
	if err != nil {
		return
	}

	// Update video status
	_, err = s.db.Exec("UPDATE videos SET status = ? WHERE id = ?", models.VideoStatusAnalyzed, videoID)
}

// createMockResults creates mock analysis results
func (s *AnalysisService) createMockResults(jobID, videoID string) error {
	// Create analysis result
	resultID := uuid.New().String()
	now := time.Now()

	_, err := s.db.Exec(`
		INSERT INTO analysis_results (id, video_id, analysis_job_id, total_frames, total_people, unique_people, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, resultID, videoID, jobID, 300, 1500, 5, now)

	if err != nil {
		return fmt.Errorf("failed to create analysis result: %v", err)
	}

	// Create mock people per frame data
	for frame := 0; frame < 300; frame += 10 {
		count := rand.Intn(8) + 1          // 1-8 people per frame
		timestamp := float64(frame) / 30.0 // Assuming 30 fps

		_, err := s.db.Exec(`
			INSERT INTO people_per_frame (id, analysis_result_id, frame_number, timestamp, count)
			VALUES (?, ?, ?, ?, ?)
		`, uuid.New().String(), resultID, frame, timestamp, count)

		if err != nil {
			return fmt.Errorf("failed to create people per frame data: %v", err)
		}
	}

	// Create mock tracking data and persons
	err = s.createMockPersons(videoID, resultID)
	if err != nil {
		return fmt.Errorf("failed to create mock persons: %v", err)
	}

	return nil
}

// createMockPersons creates mock person data with faces
func (s *AnalysisService) createMockPersons(videoID, resultID string) error {
	personIDs := []string{"person_1", "person_2", "person_3", "person_4", "person_5"}

	for i, personID := range personIDs {
		// Create person record
		personUUID := uuid.New().String()
		firstSeen := float64(i * 60)         // Each person appears at different times
		lastSeen := firstSeen + 120.0        // 2 minutes duration
		totalFrames := rand.Intn(1800) + 600 // 600-2400 frames

		_, err := s.db.Exec(`
			INSERT INTO persons (id, video_id, person_id, first_seen, last_seen, total_frames, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, personUUID, videoID, personID, firstSeen, lastSeen, totalFrames, time.Now())

		if err != nil {
			return fmt.Errorf("failed to create person: %v", err)
		}

		// Create tracking data for this person
		for frame := int(firstSeen * 30); frame < int(lastSeen*30); frame += 30 {
			timestamp := float64(frame) / 30.0
			x := rand.Float64() * 0.8
			y := rand.Float64() * 0.8
			width := 0.1 + rand.Float64()*0.2
			height := 0.1 + rand.Float64()*0.2
			confidence := 0.7 + rand.Float64()*0.3

			_, err := s.db.Exec(`
				INSERT INTO tracking_data (id, analysis_result_id, person_id, frame_number, timestamp, x, y, width, height, confidence)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, uuid.New().String(), resultID, personID, frame, timestamp, x, y, width, height, confidence)

			if err != nil {
				return fmt.Errorf("failed to create tracking data: %v", err)
			}
		}

		// Create face images for this person
		err = s.createMockFaces(personUUID, videoID, personID, firstSeen, lastSeen)
		if err != nil {
			return fmt.Errorf("failed to create faces for person %s: %v", personID, err)
		}
	}

	return nil
}

// createMockFaces creates mock face images for a person
func (s *AnalysisService) createMockFaces(personUUID, videoID, personID string, firstSeen, lastSeen float64) error {
	// Create a simple SVG face image as base64
	faceSVG := fmt.Sprintf(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
		<circle cx="50" cy="50" r="40" fill="#FFD700" stroke="#000" stroke-width="2"/>
		<circle cx="35" cy="40" r="3" fill="#000"/>
		<circle cx="65" cy="40" r="3" fill="#000"/>
		<path d="M 30 60 Q 50 70 70 60" stroke="#000" stroke-width="2" fill="none"/>
		<text x="50" y="85" text-anchor="middle" font-size="8" fill="#000">%s</text>
	</svg>`, personID)

	faceBase64 := base64.StdEncoding.EncodeToString([]byte(faceSVG))

	// Create multiple face detections for this person
	numFaces := rand.Intn(10) + 5 // 5-15 face detections per person

	for i := 0; i < numFaces; i++ {
		timestamp := firstSeen + (lastSeen-firstSeen)*float64(i)/float64(numFaces-1)
		frameNumber := int(timestamp * 30)
		x := rand.Float64() * 0.8
		y := rand.Float64() * 0.8
		width := 0.05 + rand.Float64()*0.1
		height := 0.05 + rand.Float64()*0.1
		confidence := 0.8 + rand.Float64()*0.2
		isBestFace := i == numFaces/2 // Middle face is the best

		_, err := s.db.Exec(`
			INSERT INTO person_faces (id, person_id, video_id, frame_number, timestamp, face_image, x, y, width, height, confidence, is_best_face, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, uuid.New().String(), personUUID, videoID, frameNumber, timestamp, faceBase64, x, y, width, height, confidence, isBestFace, time.Now())

		if err != nil {
			return fmt.Errorf("failed to create face: %v", err)
		}
	}

	return nil
}
