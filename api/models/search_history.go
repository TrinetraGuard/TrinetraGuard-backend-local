package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SearchRecord represents a search history record
type SearchRecord struct {
	ID              string    `json:"id"`
	SearchImagePath string    `json:"search_image_path"`
	SearchTime      time.Time `json:"search_time"`
	QueryHash       string    `json:"query_hash"` // Hash of the search image for deduplication
	MatchesFound    int       `json:"matches_found"`
	TotalVideos     int       `json:"total_videos"`
	MatchedVideos   []string  `json:"matched_videos"` // List of video IDs that had matches
	ProcessingTime  float64   `json:"processing_time"`
}

// SearchHistory manages search history records
type SearchHistory struct {
	filepath string
	Records  map[string]*SearchRecord `json:"records"`
}

// NewSearchHistory creates a new search history instance
func NewSearchHistory(filepath string) *SearchHistory {
	return &SearchHistory{
		filepath: filepath,
		Records:  make(map[string]*SearchRecord),
	}
}

// Load loads search history from JSON file
func (sh *SearchHistory) Load() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(sh.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(sh.filepath); os.IsNotExist(err) {
		// File doesn't exist, create empty storage
		sh.Records = make(map[string]*SearchRecord)
		return sh.Save()
	}

	// Read existing file
	data, err := os.ReadFile(sh.filepath)
	if err != nil {
		return fmt.Errorf("failed to read history file: %v", err)
	}

	if len(data) == 0 {
		sh.Records = make(map[string]*SearchRecord)
		return nil
	}

	if err := json.Unmarshal(data, &sh); err != nil {
		return fmt.Errorf("failed to unmarshal history data: %v", err)
	}

	return nil
}

// Save saves search history to JSON file
func (sh *SearchHistory) Save() error {
	data, err := json.MarshalIndent(sh, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history data: %v", err)
	}

	if err := os.WriteFile(sh.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %v", err)
	}

	return nil
}

// AddRecord adds a new search record
func (sh *SearchHistory) AddRecord(record *SearchRecord) error {
	sh.Records[record.ID] = record
	return sh.Save()
}

// GetRecord retrieves a search record by ID
func (sh *SearchHistory) GetRecord(id string) (*SearchRecord, bool) {
	record, exists := sh.Records[id]
	return record, exists
}

// ListRecords returns all search records (sorted by time, newest first)
func (sh *SearchHistory) ListRecords() []*SearchRecord {
	var records []*SearchRecord
	for _, record := range sh.Records {
		records = append(records, record)
	}

	// Sort by search time (newest first)
	for i := 0; i < len(records)-1; i++ {
		for j := i + 1; j < len(records); j++ {
			if records[i].SearchTime.Before(records[j].SearchTime) {
				records[i], records[j] = records[j], records[i]
			}
		}
	}

	return records
}

// GetStats returns search history statistics
func (sh *SearchHistory) GetStats() map[string]interface{} {
	totalSearches := len(sh.Records)
	totalMatches := 0
	successfulSearches := 0

	for _, record := range sh.Records {
		totalMatches += record.MatchesFound
		if record.MatchesFound > 0 {
			successfulSearches++
		}
	}

	return map[string]interface{}{
		"total_searches":      totalSearches,
		"successful_searches": successfulSearches,
		"total_matches_found": totalMatches,
		"success_rate":        float64(successfulSearches) / float64(totalSearches) * 100,
	}
}
