package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// VideoRecord represents a video processing record
type VideoRecord struct {
	ID               string    `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	StoredPath       string    `json:"stored_path"`
	UploadTime       time.Time `json:"upload_time"`
	Status           string    `json:"status"` // "processing", "completed", "failed"
	ProcessingTime   float64   `json:"processing_time,omitempty"`
	UniqueFacesCount int       `json:"unique_faces_count,omitempty"`
	FaceImages       []string  `json:"face_images,omitempty"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	IsArchived       bool      `json:"is_archived"` // New field to mark as history
	LastAccessed     time.Time `json:"last_accessed,omitempty"`
	AccessCount      int       `json:"access_count,omitempty"`
	// Location information
	LocationName string  `json:"location_name,omitempty"`
	Latitude     float64 `json:"latitude,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
}

// VideoStorage manages video records
type VideoStorage struct {
	filepath string
	Records  map[string]*VideoRecord `json:"records"`
}

// NewVideoStorage creates a new video storage instance
func NewVideoStorage(filepath string) *VideoStorage {
	return &VideoStorage{
		filepath: filepath,
		Records:  make(map[string]*VideoRecord),
	}
}

// Load loads video records from JSON file
func (vs *VideoStorage) Load() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(vs.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(vs.filepath); os.IsNotExist(err) {
		// File doesn't exist, create empty storage
		vs.Records = make(map[string]*VideoRecord)
		return vs.Save()
	}

	// Read existing file
	data, err := os.ReadFile(vs.filepath)
	if err != nil {
		return fmt.Errorf("failed to read storage file: %v", err)
	}

	if len(data) == 0 {
		vs.Records = make(map[string]*VideoRecord)
		return nil
	}

	if err := json.Unmarshal(data, &vs); err != nil {
		return fmt.Errorf("failed to unmarshal storage data: %v", err)
	}

	return nil
}

// Save saves video records to JSON file
func (vs *VideoStorage) Save() error {
	data, err := json.MarshalIndent(vs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal storage data: %v", err)
	}

	if err := os.WriteFile(vs.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write storage file: %v", err)
	}

	return nil
}

// AddRecord adds a new video record
func (vs *VideoStorage) AddRecord(record *VideoRecord) error {
	vs.Records[record.ID] = record
	return vs.Save()
}

// GetRecord retrieves a video record by ID
func (vs *VideoStorage) GetRecord(id string) (*VideoRecord, bool) {
	record, exists := vs.Records[id]
	if exists && record != nil {
		// Update access statistics
		record.LastAccessed = time.Now()
		record.AccessCount++
		vs.Save() // Save the updated access info
	}
	return record, exists
}

// UpdateRecord updates an existing video record
func (vs *VideoStorage) UpdateRecord(record *VideoRecord) error {
	if _, exists := vs.Records[record.ID]; !exists {
		return fmt.Errorf("record not found: %s", record.ID)
	}
	vs.Records[record.ID] = record
	return vs.Save()
}

// DeleteRecord deletes a video record (but keeps the files for history)
func (vs *VideoStorage) DeleteRecord(id string) error {
	record, exists := vs.Records[id]
	if !exists {
		return fmt.Errorf("record not found: %s", id)
	}

	// Mark as archived instead of deleting
	record.IsArchived = true
	record.LastAccessed = time.Now()
	vs.Records[id] = record
	return vs.Save()
}

// ListRecords returns all video records
func (vs *VideoStorage) ListRecords() []*VideoRecord {
	var records []*VideoRecord
	for _, record := range vs.Records {
		records = append(records, record)
	}
	return records
}

// ListActiveRecords returns only non-archived records
func (vs *VideoStorage) ListActiveRecords() []*VideoRecord {
	var records []*VideoRecord
	for _, record := range vs.Records {
		if !record.IsArchived {
			records = append(records, record)
		}
	}
	return records
}

// ListArchivedRecords returns only archived records (history)
func (vs *VideoStorage) ListArchivedRecords() []*VideoRecord {
	var records []*VideoRecord
	for _, record := range vs.Records {
		if record.IsArchived {
			records = append(records, record)
		}
	}
	return records
}

// GetStats returns storage statistics
func (vs *VideoStorage) GetStats() map[string]interface{} {
	totalRecords := len(vs.Records)
	activeRecords := 0
	archivedRecords := 0
	totalFaces := 0
	totalProcessingTime := 0.0
	locationsWithGPS := 0

	for _, record := range vs.Records {
		if record.IsArchived {
			archivedRecords++
		} else {
			activeRecords++
		}
		totalFaces += record.UniqueFacesCount
		totalProcessingTime += record.ProcessingTime

		// Count records with GPS coordinates
		if record.Latitude != 0 && record.Longitude != 0 {
			locationsWithGPS++
		}
	}

	return map[string]interface{}{
		"total_records":         totalRecords,
		"active_records":        activeRecords,
		"archived_records":      archivedRecords,
		"total_faces_detected":  totalFaces,
		"total_processing_time": totalProcessingTime,
		"locations_with_gps":    locationsWithGPS,
	}
}

// CleanupOldRecords removes very old archived records (optional, for disk space management)
func (vs *VideoStorage) CleanupOldRecords(daysToKeep int) error {
	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)
	var recordsToDelete []string

	for id, record := range vs.Records {
		if record.IsArchived && record.LastAccessed.Before(cutoffTime) {
			recordsToDelete = append(recordsToDelete, id)
		}
	}

	// Delete old records
	for _, id := range recordsToDelete {
		delete(vs.Records, id)
	}

	if len(recordsToDelete) > 0 {
		return vs.Save()
	}

	return nil
}

// ResetDatabase completely resets the database and removes all files
func (vs *VideoStorage) ResetDatabase() error {
	// Remove all video files
	for _, record := range vs.Records {
		if err := os.Remove(record.StoredPath); err != nil {
			log.Printf("Warning: Could not remove video file %s: %v", record.StoredPath, err)
		}

		// Remove face images
		for _, faceImage := range record.FaceImages {
			facePath := filepath.Join("../storage/faces", filepath.Base(faceImage))
			if err := os.Remove(facePath); err != nil {
				log.Printf("Warning: Could not remove face image %s: %v", facePath, err)
			}
		}
	}

	// Clear all records
	vs.Records = make(map[string]*VideoRecord)

	// Save empty database
	return vs.Save()
}
