package models

import (
	"time"
)

// LostPerson represents a lost person report in the system
type LostPerson struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	AadharNumber     string    `json:"aadhar_number"`
	ContactNumber    string    `json:"contact_number"`
	PlaceLost        string    `json:"place_lost"`
	PermanentAddress string    `json:"permanent_address"`
	ImagePath        string    `json:"image_path"`
	UploadTimestamp  time.Time `json:"upload_timestamp"`
}

// CreateLostPersonRequest represents the request body for creating a lost person report
type CreateLostPersonRequest struct {
	Name             string `form:"name" binding:"required"`
	AadharNumber     string `form:"aadhar_number" binding:"required"`
	ContactNumber    string `form:"contact_number"`
	PlaceLost        string `form:"place_lost" binding:"required"`
	PermanentAddress string `form:"permanent_address" binding:"required"`
}

// LostPersonResponse represents the response for lost person operations
type LostPersonResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	AadharNumber     string    `json:"aadhar_number"`
	ContactNumber    string    `json:"contact_number"`
	PlaceLost        string    `json:"place_lost"`
	PermanentAddress string    `json:"permanent_address"`
	ImagePath        string    `json:"image_path"`
	UploadTimestamp  time.Time `json:"upload_timestamp"`
}

// LostPersonListResponse represents the response for listing lost persons
type LostPersonListResponse struct {
	Total       int                  `json:"total"`
	LostPersons []LostPersonResponse `json:"lost_persons"`
}
