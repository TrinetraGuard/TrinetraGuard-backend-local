package models

import (
	"encoding/json"
	"time"
)

// Video represents a video file in the system
type Video struct {
	ID               string    `json:"id" db:"id"`
	Filename         string    `json:"filename" db:"filename"`
	OriginalFilename string    `json:"original_filename" db:"original_filename"`
	FileSize         int64     `json:"file_size" db:"file_size"`
	Duration         float64   `json:"duration,omitempty" db:"duration"`
	FrameCount       int       `json:"frame_count,omitempty" db:"frame_count"`
	Width            int       `json:"width,omitempty" db:"width"`
	Height           int       `json:"height,omitempty" db:"height"`
	Status           string    `json:"status" db:"status"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// AnalysisJob represents an analysis job for a video
type AnalysisJob struct {
	ID           string     `json:"id" db:"id"`
	VideoID      string     `json:"video_id" db:"video_id"`
	Status       string     `json:"status" db:"status"`
	Progress     int        `json:"progress" db:"progress"`
	ErrorMessage string     `json:"error_message,omitempty" db:"error_message"`
	StartedAt    *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// AnalysisResult represents the results of video analysis
type AnalysisResult struct {
	ID             string          `json:"id" db:"id"`
	JobID          string          `json:"job_id" db:"job_id"`
	VideoID        string          `json:"video_id" db:"video_id"`
	TotalFrames    int             `json:"total_frames" db:"total_frames"`
	TotalPeople    int             `json:"total_people" db:"total_people"`
	UniquePeople   int             `json:"unique_people" db:"unique_people"`
	PeoplePerFrame json.RawMessage `json:"people_per_frame" db:"people_per_frame"`
	TrackingData   json.RawMessage `json:"tracking_data" db:"tracking_data"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

// PeoplePerFrame represents people count per frame
type PeoplePerFrame struct {
	FrameNumber int     `json:"frame_number"`
	Timestamp   float64 `json:"timestamp"`
	Count       int     `json:"count"`
	Confidence  float64 `json:"confidence"`
}

// TrackingData represents individual person tracking
type TrackingData struct {
	PersonID    string      `json:"person_id"`
	FrameNumber int         `json:"frame_number"`
	Timestamp   float64     `json:"timestamp"`
	BoundingBox BoundingBox `json:"bounding_box"`
	Confidence  float64     `json:"confidence"`
}

// BoundingBox represents a bounding box around a detected person
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ReferenceImage represents a reference image for person finding
type ReferenceImage struct {
	ID               string    `json:"id" db:"id"`
	Filename         string    `json:"filename" db:"filename"`
	OriginalFilename string    `json:"original_filename" db:"original_filename"`
	FileSize         int64     `json:"file_size" db:"file_size"`
	Description      string    `json:"description,omitempty" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// SearchJob represents a search job for finding a person
type SearchJob struct {
	ID               string     `json:"id" db:"id"`
	ReferenceImageID string     `json:"reference_image_id" db:"reference_image_id"`
	Status           string     `json:"status" db:"status"`
	Progress         int        `json:"progress" db:"progress"`
	ErrorMessage     string     `json:"error_message,omitempty" db:"error_message"`
	StartedAt        *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// SearchResult represents the results of a person search
type SearchResult struct {
	ID               string          `json:"id" db:"id"`
	SearchJobID      string          `json:"search_job_id" db:"search_job_id"`
	VideoID          string          `json:"video_id" db:"video_id"`
	ReferenceImageID string          `json:"reference_image_id" db:"reference_image_id"`
	Matches          json.RawMessage `json:"matches" db:"matches"`
	FirstAppearance  float64         `json:"first_appearance" db:"first_appearance"`
	LastAppearance   float64         `json:"last_appearance" db:"last_appearance"`
	TotalAppearances int             `json:"total_appearances" db:"total_appearances"`
	Confidence       float64         `json:"confidence" db:"confidence"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// Match represents a match found in a video
type Match struct {
	FrameNumber int         `json:"frame_number"`
	Timestamp   float64     `json:"timestamp"`
	BoundingBox BoundingBox `json:"bounding_box"`
	Confidence  float64     `json:"confidence"`
}

// API Response structures

// UploadVideoResponse represents the response for video upload
type UploadVideoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Video   Video  `json:"video,omitempty"`
}

// ListVideosResponse represents the response for listing videos
type ListVideosResponse struct {
	Success bool    `json:"success"`
	Videos  []Video `json:"videos"`
	Total   int     `json:"total"`
}

// AnalysisStatusResponse represents the response for analysis status
type AnalysisStatusResponse struct {
	Success bool        `json:"success"`
	Job     AnalysisJob `json:"job"`
}

// AnalysisResultsResponse represents the response for analysis results
type AnalysisResultsResponse struct {
	Success bool           `json:"success"`
	Results AnalysisResult `json:"results"`
}

// SearchPersonRequest represents the request for person search
type SearchPersonRequest struct {
	ReferenceImageID string   `json:"reference_image_id" binding:"required"`
	VideoIDs         []string `json:"video_ids,omitempty"`
}

// SearchPersonResponse represents the response for person search
type SearchPersonResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	SearchJob SearchJob `json:"search_job"`
}

// SearchResultsResponse represents the response for search results
type SearchResultsResponse struct {
	Success bool           `json:"success"`
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// Job status constants
const (
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"
	JobStatusCancelled = "cancelled"
)

// Video status constants
const (
	VideoStatusUploaded  = "uploaded"
	VideoStatusAnalyzing = "analyzing"
	VideoStatusAnalyzed  = "analyzed"
	VideoStatusFailed    = "failed"
)

// PersonFace represents a detected face of a person
type PersonFace struct {
	ID          string      `json:"id" db:"id"`
	PersonID    string      `json:"person_id" db:"person_id"`
	VideoID     string      `json:"video_id" db:"video_id"`
	FrameNumber int         `json:"frame_number" db:"frame_number"`
	Timestamp   float64     `json:"timestamp" db:"timestamp"`
	BoundingBox BoundingBox `json:"bounding_box" db:"bounding_box"`
	Confidence  float64     `json:"confidence" db:"confidence"`
	FaceImage   string      `json:"face_image" db:"face_image"` // Base64 encoded face image
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

// Person represents a unique person detected in a video
type Person struct {
	ID           string     `json:"id" db:"id"`
	VideoID      string     `json:"video_id" db:"video_id"`
	PersonNumber int        `json:"person_number" db:"person_number"`
	FirstFrame   int        `json:"first_frame" db:"first_frame"`
	LastFrame    int        `json:"last_frame" db:"last_frame"`
	FirstTime    float64    `json:"first_time" db:"first_time"`
	LastTime     float64    `json:"last_time" db:"last_time"`
	TotalFrames  int        `json:"total_frames" db:"total_frames"`
	BestFace     PersonFace `json:"best_face,omitempty" db:"best_face"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// Enhanced AnalysisResult with face data
type EnhancedAnalysisResult struct {
	ID             string          `json:"id" db:"id"`
	JobID          string          `json:"job_id" db:"job_id"`
	VideoID        string          `json:"video_id" db:"video_id"`
	TotalFrames    int             `json:"total_frames" db:"total_frames"`
	TotalPeople    int             `json:"total_people" db:"total_people"`
	UniquePeople   int             `json:"unique_people" db:"unique_people"`
	PeoplePerFrame json.RawMessage `json:"people_per_frame" db:"people_per_frame"`
	TrackingData   json.RawMessage `json:"tracking_data" db:"tracking_data"`
	Persons        []Person        `json:"persons,omitempty" db:"persons"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}
