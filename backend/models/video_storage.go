package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// VideoRecord represents a stored video record
type VideoRecord struct {
	ID               string    `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	StoredPath       string    `json:"stored_path"`
	UploadTime       time.Time `json:"upload_time"`
	ProcessingTime   float64   `json:"processing_time_seconds"`
	UniqueFacesCount int       `json:"unique_faces_count"`
	FaceImages       []string  `json:"face_images"`
	Status           string    `json:"status"` // "processing", "completed", "failed"
	ErrorMessage     string    `json:"error_message,omitempty"`
	VideoDuration    float64   `json:"video_duration,omitempty"`
	FramesProcessed  int       `json:"frames_processed,omitempty"`
}

// VideoStorage handles local storage of video records
type VideoStorage struct {
	storageFile string
	records     map[string]*VideoRecord
}

// NewVideoStorage creates a new video storage instance
func NewVideoStorage(storageFile string) *VideoStorage {
	return &VideoStorage{
		storageFile: storageFile,
		records:     make(map[string]*VideoRecord),
	}
}

// Load loads existing records from storage file
func (vs *VideoStorage) Load() error {
	if _, err := os.Stat(vs.storageFile); os.IsNotExist(err) {
		// Create storage file if it doesn't exist
		return vs.Save()
	}

	data, err := os.ReadFile(vs.storageFile)
	if err != nil {
		return fmt.Errorf("failed to read storage file: %v", err)
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, &vs.records)
}

// Save saves records to storage file
func (vs *VideoStorage) Save() error {
	data, err := json.MarshalIndent(vs.records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal records: %v", err)
	}

	return os.WriteFile(vs.storageFile, data, 0644)
}

// AddRecord adds a new video record
func (vs *VideoStorage) AddRecord(record *VideoRecord) error {
	vs.records[record.ID] = record
	return vs.Save()
}

// GetRecord retrieves a video record by ID
func (vs *VideoStorage) GetRecord(id string) (*VideoRecord, bool) {
	record, exists := vs.records[id]
	return record, exists
}

// GetAllRecords returns all video records
func (vs *VideoStorage) GetAllRecords() []*VideoRecord {
	records := make([]*VideoRecord, 0, len(vs.records))
	for _, record := range vs.records {
		records = append(records, record)
	}
	return records
}

// UpdateRecord updates an existing video record
func (vs *VideoStorage) UpdateRecord(record *VideoRecord) error {
	vs.records[record.ID] = record
	return vs.Save()
}

// DeleteRecord deletes a video record
func (vs *VideoStorage) DeleteRecord(id string) error {
	if record, exists := vs.records[id]; exists {
		// Remove associated files
		if record.StoredPath != "" {
			os.Remove(record.StoredPath)
		}
		for _, facePath := range record.FaceImages {
			os.Remove(facePath)
		}
		delete(vs.records, id)
		return vs.Save()
	}
	return fmt.Errorf("record not found: %s", id)
}

// CleanupOldRecords removes records older than specified days
func (vs *VideoStorage) CleanupOldRecords(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	deleted := 0

	for id, record := range vs.records {
		if record.UploadTime.Before(cutoff) {
			vs.DeleteRecord(id)
			deleted++
		}
	}

	if deleted > 0 {
		return vs.Save()
	}
	return nil
}
