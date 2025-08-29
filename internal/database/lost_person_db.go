package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/trinetraguard/backend/internal/models"
)

// LostPersonDB manages the JSON database for lost person records
type LostPersonDB struct {
	filePath string
	mutex    sync.RWMutex
}

// NewLostPersonDB creates a new LostPersonDB instance
func NewLostPersonDB() *LostPersonDB {
	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create uploads directory: %v", err))
	}

	return &LostPersonDB{
		filePath: "database.json",
	}
}

// readData reads all records from the JSON file
func (db *LostPersonDB) readData() ([]models.LostPerson, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	// Check if file exists
	if _, err := os.Stat(db.filePath); os.IsNotExist(err) {
		// Return empty slice if file doesn't exist
		return []models.LostPerson{}, nil
	}

	// Read file
	data, err := ioutil.ReadFile(db.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database file: %v", err)
	}

	// Parse JSON
	var records []models.LostPerson
	if len(data) > 0 {
		if err := json.Unmarshal(data, &records); err != nil {
			return nil, fmt.Errorf("failed to parse database file: %v", err)
		}
	}

	return records, nil
}

// writeData writes all records to the JSON file
func (db *LostPersonDB) writeData(records []models.LostPerson) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Marshal to JSON
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// Write to file
	if err := ioutil.WriteFile(db.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write database file: %v", err)
	}

	return nil
}

// CreateLostPerson creates a new lost person record
func (db *LostPersonDB) CreateLostPerson(req models.CreateLostPersonRequest, imagePath string) (*models.LostPerson, error) {
	// Read existing records
	records, err := db.readData()
	if err != nil {
		return nil, err
	}

	// Create new record
	now := time.Now()
	lostPerson := models.LostPerson{
		ID:               uuid.New().String(),
		Name:             req.Name,
		AadharNumber:     req.AadharNumber,
		ContactNumber:    req.ContactNumber,
		PlaceLost:        req.PlaceLost,
		PermanentAddress: req.PermanentAddress,
		ImagePath:        imagePath,
		UploadTimestamp:  now,
	}

	// Add to records
	records = append(records, lostPerson)

	// Write back to file
	if err := db.writeData(records); err != nil {
		return nil, err
	}

	return &lostPerson, nil
}

// GetAllLostPersons returns all lost person records
func (db *LostPersonDB) GetAllLostPersons() ([]models.LostPerson, error) {
	return db.readData()
}

// GetLostPersonByID returns a specific lost person by ID
func (db *LostPersonDB) GetLostPersonByID(id string) (*models.LostPerson, error) {
	records, err := db.readData()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record.ID == id {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("lost person not found with ID: %s", id)
}

// UpdateLostPerson updates an existing lost person record
func (db *LostPersonDB) UpdateLostPerson(id string, req models.CreateLostPersonRequest, imagePath string) (*models.LostPerson, error) {
	records, err := db.readData()
	if err != nil {
		return nil, err
	}

	// Find and update record
	for i, record := range records {
		if record.ID == id {
			records[i].Name = req.Name
			records[i].AadharNumber = req.AadharNumber
			records[i].ContactNumber = req.ContactNumber
			records[i].PlaceLost = req.PlaceLost
			records[i].PermanentAddress = req.PermanentAddress
			if imagePath != "" {
				records[i].ImagePath = imagePath
			}

			// Write back to file
			if err := db.writeData(records); err != nil {
				return nil, err
			}

			return &records[i], nil
		}
	}

	return nil, fmt.Errorf("lost person not found with ID: %s", id)
}

// DeleteLostPerson deletes a lost person record
func (db *LostPersonDB) DeleteLostPerson(id string) error {
	records, err := db.readData()
	if err != nil {
		return err
	}

	// Find and remove record
	for i, record := range records {
		if record.ID == id {
			// Remove the record
			records = append(records[:i], records[i+1:]...)

			// Write back to file
			if err := db.writeData(records); err != nil {
				return err
			}

			// Try to delete the image file
			if record.ImagePath != "" {
				os.Remove(record.ImagePath)
			}

			return nil
		}
	}

	return fmt.Errorf("lost person not found with ID: %s", id)
}

// SearchLostPersons searches for lost persons by name or Aadhar number
func (db *LostPersonDB) SearchLostPersons(query string) ([]models.LostPerson, error) {
	records, err := db.readData()
	if err != nil {
		return nil, err
	}

	var results []models.LostPerson
	for _, record := range records {
		if record.Name == query || record.AadharNumber == query {
			results = append(results, record)
		}
	}

	return results, nil
}
