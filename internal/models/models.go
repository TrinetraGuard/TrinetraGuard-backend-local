package models

import (
	"time"
)

// Video represents a video file in the system
type Video struct {
	ID           string    `json:"id" db:"id"`
	Filename     string    `json:"filename" db:"filename"`
	OriginalName string    `json:"original_name" db:"original_name"`
	FilePath     string    `json:"file_path" db:"file_path"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	Duration     float64   `json:"duration" db:"duration"`
	FrameCount   int       `json:"frame_count" db:"frame_count"`
	Width        int       `json:"width" db:"width"`
	Height       int       `json:"height" db:"height"`
	Status       string    `json:"status" db:"status"`
	UploadedAt   time.Time `json:"uploaded_at" db:"uploaded_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// AnalysisJob represents a video analysis job
type AnalysisJob struct {
	ID          string     `json:"id" db:"id"`
	VideoID     string     `json:"video_id" db:"video_id"`
	Status      string     `json:"status" db:"status"`
	Progress    int        `json:"progress" db:"progress"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	Error       *string    `json:"error" db:"error"`
}

// AnalysisResult represents the results of video analysis
type AnalysisResult struct {
	ID             string           `json:"id" db:"id"`
	VideoID        string           `json:"video_id" db:"video_id"`
	AnalysisJobID  string           `json:"analysis_job_id" db:"analysis_job_id"`
	TotalFrames    int              `json:"total_frames" db:"total_frames"`
	TotalPeople    int              `json:"total_people" db:"total_people"`
	UniquePeople   int              `json:"unique_people" db:"unique_people"`
	PeoplePerFrame []PeoplePerFrame `json:"people_per_frame" db:"people_per_frame"`
	TrackingData   []TrackingData   `json:"tracking_data" db:"tracking_data"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
}

// PeoplePerFrame represents people count in a specific frame
type PeoplePerFrame struct {
	FrameNumber int     `json:"frame_number" db:"frame_number"`
	Timestamp   float64 `json:"timestamp" db:"timestamp"`
	Count       int     `json:"count" db:"count"`
}

// TrackingData represents individual person tracking
type TrackingData struct {
	PersonID    string      `json:"person_id" db:"person_id"`
	FrameNumber int         `json:"frame_number" db:"frame_number"`
	Timestamp   float64     `json:"timestamp" db:"timestamp"`
	BoundingBox BoundingBox `json:"bounding_box" db:"bounding_box"`
	Confidence  float64     `json:"confidence" db:"confidence"`
}

// BoundingBox represents a bounding box for detected objects
type BoundingBox struct {
	X      float64 `json:"x" db:"x"`
	Y      float64 `json:"y" db:"y"`
	Width  float64 `json:"width" db:"width"`
	Height float64 `json:"height" db:"height"`
}

// ReferenceImage represents a reference image for person finding
type ReferenceImage struct {
	ID           string    `json:"id" db:"id"`
	Filename     string    `json:"filename" db:"filename"`
	OriginalName string    `json:"original_name" db:"original_name"`
	FilePath     string    `json:"file_path" db:"file_path"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	UploadedAt   time.Time `json:"uploaded_at" db:"uploaded_at"`
}

// SearchJob represents a person search job
type SearchJob struct {
	ID               string     `json:"id" db:"id"`
	ReferenceImageID string     `json:"reference_image_id" db:"reference_image_id"`
	Status           string     `json:"status" db:"status"`
	Progress         int        `json:"progress" db:"progress"`
	StartedAt        time.Time  `json:"started_at" db:"started_at"`
	CompletedAt      *time.Time `json:"completed_at" db:"completed_at"`
	Error            *string    `json:"error" db:"error"`
}

// SearchResult represents the results of a person search
type SearchResult struct {
	ID          string    `json:"id" db:"id"`
	SearchJobID string    `json:"search_job_id" db:"search_job_id"`
	VideoID     string    `json:"video_id" db:"video_id"`
	Matches     []Match   `json:"matches" db:"matches"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Match represents a match found in a video
type Match struct {
	PersonID    string      `json:"person_id" db:"person_id"`
	FirstFrame  int         `json:"first_frame" db:"first_frame"`
	LastFrame   int         `json:"last_frame" db:"last_frame"`
	FirstTime   float64     `json:"first_time" db:"first_time"`
	LastTime    float64     `json:"last_time" db:"last_time"`
	Confidence  float64     `json:"confidence" db:"confidence"`
	BoundingBox BoundingBox `json:"bounding_box" db:"bounding_box"`
}

// PersonFace represents a detected face image
type PersonFace struct {
	ID          string      `json:"id" db:"id"`
	PersonID    string      `json:"person_id" db:"person_id"`
	VideoID     string      `json:"video_id" db:"video_id"`
	FrameNumber int         `json:"frame_number" db:"frame_number"`
	Timestamp   float64     `json:"timestamp" db:"timestamp"`
	FaceImage   string      `json:"face_image" db:"face_image"` // Base64 encoded image
	BoundingBox BoundingBox `json:"bounding_box" db:"bounding_box"`
	Confidence  float64     `json:"confidence" db:"confidence"`
	IsBestFace  bool        `json:"is_best_face" db:"is_best_face"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

// Person represents a unique person detected in a video
type Person struct {
	ID          string       `json:"id" db:"id"`
	VideoID     string       `json:"video_id" db:"video_id"`
	PersonID    string       `json:"person_id" db:"person_id"`
	FirstSeen   float64      `json:"first_seen" db:"first_seen"`
	LastSeen    float64      `json:"last_seen" db:"last_seen"`
	TotalFrames int          `json:"total_frames" db:"total_frames"`
	BestFace    *PersonFace  `json:"best_face" db:"best_face"`
	Faces       []PersonFace `json:"faces" db:"faces"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
}

// EnhancedAnalysisResult extends AnalysisResult with person data
type EnhancedAnalysisResult struct {
	*AnalysisResult
	Persons []Person `json:"persons"`
}

// FaceDetectionResult represents face detection results for a frame
type FaceDetectionResult struct {
	FrameNumber int            `json:"frame_number"`
	Timestamp   float64        `json:"timestamp"`
	Faces       []DetectedFace `json:"faces"`
}

// DetectedFace represents a single detected face
type DetectedFace struct {
	PersonID    string      `json:"person_id"`
	BoundingBox BoundingBox `json:"bounding_box"`
	Confidence  float64     `json:"confidence"`
	FaceImage   string      `json:"face_image"` // Base64 encoded
}

// PersonActivity represents activity timeline for a person
type PersonActivity struct {
	PersonID    string      `json:"person_id"`
	FirstSeen   float64     `json:"first_seen"`
	LastSeen    float64     `json:"last_seen"`
	TotalFrames int         `json:"total_frames"`
	TimeRanges  []TimeRange `json:"time_ranges"`
	BestFace    *PersonFace `json:"best_face"`
}

// TimeRange represents a time range when a person appears
type TimeRange struct {
	StartTime  float64 `json:"start_time"`
	EndTime    float64 `json:"end_time"`
	StartFrame int     `json:"start_frame"`
	EndFrame   int     `json:"end_frame"`
	Duration   float64 `json:"duration"`
}

// API Request/Response Models

// UploadVideoResponse represents the response for video upload
type UploadVideoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Video   Video  `json:"video"`
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

// EnhancedAnalysisResultsResponse represents the response for enhanced analysis results
type EnhancedAnalysisResultsResponse struct {
	Success bool                   `json:"success"`
	Results EnhancedAnalysisResult `json:"results"`
}

// SearchPersonRequest represents the request for person search
type SearchPersonRequest struct {
	ReferenceImageID string `json:"reference_image_id" binding:"required"`
}

// SearchPersonResponse represents the response for person search
type SearchPersonResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Job     SearchJob `json:"job"`
}

// SearchResultsResponse represents the response for search results
type SearchResultsResponse struct {
	Success bool         `json:"success"`
	Results SearchResult `json:"results"`
}

// SearchStatusResponse represents the response for search status
type SearchStatusResponse struct {
	Success bool      `json:"success"`
	Job     SearchJob `json:"job"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// Job and Video Status Constants
const (
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"

	VideoStatusUploaded = "uploaded"
	VideoStatusAnalyzed = "analyzed"
	VideoStatusFailed   = "failed"
)
