package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/trinetraguard/backend/internal/database"
	"github.com/trinetraguard/backend/internal/models"
	"github.com/trinetraguard/backend/internal/utils"
)

// LostPersonHandler handles lost person-related requests
type LostPersonHandler struct {
	db *database.LostPersonDB
}

// NewLostPersonHandler creates a new LostPersonHandler instance
func NewLostPersonHandler() *LostPersonHandler {
	return &LostPersonHandler{
		db: database.NewLostPersonDB(),
	}
}

// CreateLostPerson handles the creation of a new lost person report
func (h *LostPersonHandler) CreateLostPerson(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		utils.BadRequest(c, "Failed to parse form data")
		return
	}

	// Get form data
	var req models.CreateLostPersonRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequest(c, fmt.Sprintf("Invalid form data: %v", err))
		return
	}

	// Validate required fields
	if req.Name == "" || req.AadharNumber == "" || req.PlaceLost == "" || req.PermanentAddress == "" {
		utils.BadRequest(c, "Name, Aadhar Number, Place Lost, and Permanent Address are required")
		return
	}

	// Handle file upload
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		utils.BadRequest(c, "Image file is required")
		return
	}
	defer file.Close()

	// Validate file type
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif"}
	contentType := header.Header.Get("Content-Type")

	// Validate content type if it's set
	isValidType := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			isValidType = true
			break
		}
	}

	// If content type is not valid, check file extension
	if !isValidType {
		ext := filepath.Ext(header.Filename)
		allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
		isValidExtension := false
		for _, allowedExt := range allowedExtensions {
			if strings.ToLower(ext) == allowedExt {
				isValidExtension = true
				break
			}
		}
		if !isValidExtension {
			utils.BadRequest(c, "Invalid file type. Only JPEG, JPG, PNG, and GIF are allowed")
			return
		}
	}

	// Validate file size (max 10MB)
	if header.Size > 10<<20 {
		utils.BadRequest(c, "File size too large. Maximum size is 10MB")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join("uploads", filename)

	// Create file
	dst, err := os.Create(filePath)
	if err != nil {
		utils.InternalServerError(c, "Failed to create file")
		return
	}
	defer dst.Close()

	// Copy file content
	if err := c.SaveUploadedFile(header, filePath); err != nil {
		utils.InternalServerError(c, "Failed to save file")
		return
	}

	// Create record in database
	lostPerson, err := h.db.CreateLostPerson(req, filePath)
	if err != nil {
		// Clean up file if database operation fails
		os.Remove(filePath)
		utils.InternalServerError(c, fmt.Sprintf("Failed to create record: %v", err))
		return
	}

	// Convert to response format
	response := models.LostPersonResponse{
		ID:               lostPerson.ID,
		Name:             lostPerson.Name,
		AadharNumber:     lostPerson.AadharNumber,
		ContactNumber:    lostPerson.ContactNumber,
		PlaceLost:        lostPerson.PlaceLost,
		PermanentAddress: lostPerson.PermanentAddress,
		ImagePath:        lostPerson.ImagePath,
		UploadTimestamp:  lostPerson.UploadTimestamp,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Lost person report created successfully", response)
}

// GetAllLostPersons returns all lost person reports
func (h *LostPersonHandler) GetAllLostPersons(c *gin.Context) {
	lostPersons, err := h.db.GetAllLostPersons()
	if err != nil {
		utils.InternalServerError(c, fmt.Sprintf("Failed to retrieve records: %v", err))
		return
	}

	// Convert to response format
	var responses []models.LostPersonResponse
	for _, person := range lostPersons {
		response := models.LostPersonResponse{
			ID:               person.ID,
			Name:             person.Name,
			AadharNumber:     person.AadharNumber,
			ContactNumber:    person.ContactNumber,
			PlaceLost:        person.PlaceLost,
			PermanentAddress: person.PermanentAddress,
			ImagePath:        person.ImagePath,
			UploadTimestamp:  person.UploadTimestamp,
		}
		responses = append(responses, response)
	}

	response := models.LostPersonListResponse{
		Total:       len(responses),
		LostPersons: responses,
	}

	utils.SuccessResponse(c, http.StatusOK, "Lost person reports retrieved successfully", response)
}

// GetLostPersonByID returns a specific lost person report by ID
func (h *LostPersonHandler) GetLostPersonByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "ID is required")
		return
	}

	lostPerson, err := h.db.GetLostPersonByID(id)
	if err != nil {
		utils.NotFound(c, fmt.Sprintf("Lost person not found: %v", err))
		return
	}

	// Convert to response format
	response := models.LostPersonResponse{
		ID:               lostPerson.ID,
		Name:             lostPerson.Name,
		AadharNumber:     lostPerson.AadharNumber,
		ContactNumber:    lostPerson.ContactNumber,
		PlaceLost:        lostPerson.PlaceLost,
		PermanentAddress: lostPerson.PermanentAddress,
		ImagePath:        lostPerson.ImagePath,
		UploadTimestamp:  lostPerson.UploadTimestamp,
	}

	utils.SuccessResponse(c, http.StatusOK, "Lost person report retrieved successfully", response)
}

// UpdateLostPerson updates an existing lost person report
func (h *LostPersonHandler) UpdateLostPerson(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "ID is required")
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		utils.BadRequest(c, "Failed to parse form data")
		return
	}

	// Get form data
	var req models.CreateLostPersonRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequest(c, fmt.Sprintf("Invalid form data: %v", err))
		return
	}

	// Validate required fields
	if req.Name == "" || req.AadharNumber == "" || req.PlaceLost == "" || req.PermanentAddress == "" {
		utils.BadRequest(c, "Name, Aadhar Number, Place Lost, and Permanent Address are required")
		return
	}

	var imagePath string

	// Handle optional file upload
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		// Validate file type
		allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif"}
		contentType := header.Header.Get("Content-Type")

		// Validate content type if it's set
		isValidType := false
		for _, allowedType := range allowedTypes {
			if contentType == allowedType {
				isValidType = true
				break
			}
		}

		// If content type is not valid, check file extension
		if !isValidType {
			ext := filepath.Ext(header.Filename)
			allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
			isValidExtension := false
			for _, allowedExt := range allowedExtensions {
				if strings.ToLower(ext) == allowedExt {
					isValidExtension = true
					break
				}
			}
			if !isValidExtension {
				utils.BadRequest(c, "Invalid file type. Only JPEG, JPG, PNG, and GIF are allowed")
				return
			}
		}

		// Validate file size (max 10MB)
		if header.Size > 10<<20 {
			utils.BadRequest(c, "File size too large. Maximum size is 10MB")
			return
		}

		// Generate unique filename
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		imagePath = filepath.Join("uploads", filename)

		// Save file
		if err := c.SaveUploadedFile(header, imagePath); err != nil {
			utils.InternalServerError(c, "Failed to save file")
			return
		}
	}

	// Update record in database
	lostPerson, err := h.db.UpdateLostPerson(id, req, imagePath)
	if err != nil {
		// Clean up file if database operation fails
		if imagePath != "" {
			os.Remove(imagePath)
		}
		utils.NotFound(c, fmt.Sprintf("Failed to update record: %v", err))
		return
	}

	// Convert to response format
	response := models.LostPersonResponse{
		ID:               lostPerson.ID,
		Name:             lostPerson.Name,
		AadharNumber:     lostPerson.AadharNumber,
		ContactNumber:    lostPerson.ContactNumber,
		PlaceLost:        lostPerson.PlaceLost,
		PermanentAddress: lostPerson.PermanentAddress,
		ImagePath:        lostPerson.ImagePath,
		UploadTimestamp:  lostPerson.UploadTimestamp,
	}

	utils.SuccessResponse(c, http.StatusOK, "Lost person report updated successfully", response)
}

// DeleteLostPerson deletes a lost person report
func (h *LostPersonHandler) DeleteLostPerson(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "ID is required")
		return
	}

	if err := h.db.DeleteLostPerson(id); err != nil {
		utils.NotFound(c, fmt.Sprintf("Failed to delete record: %v", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Lost person report deleted successfully", gin.H{"id": id})
}

// SearchLostPersons searches for lost persons by name or Aadhar number
func (h *LostPersonHandler) SearchLostPersons(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequest(c, "Search query is required")
		return
	}

	// Trim whitespace
	query = strings.TrimSpace(query)
	if query == "" {
		utils.BadRequest(c, "Search query cannot be empty")
		return
	}

	lostPersons, err := h.db.SearchLostPersons(query)
	if err != nil {
		utils.InternalServerError(c, fmt.Sprintf("Failed to search records: %v", err))
		return
	}

	// Convert to response format
	var responses []models.LostPersonResponse
	for _, person := range lostPersons {
		response := models.LostPersonResponse{
			ID:               person.ID,
			Name:             person.Name,
			AadharNumber:     person.AadharNumber,
			ContactNumber:    person.ContactNumber,
			PlaceLost:        person.PlaceLost,
			PermanentAddress: person.PermanentAddress,
			ImagePath:        person.ImagePath,
			UploadTimestamp:  person.UploadTimestamp,
		}
		responses = append(responses, response)
	}

	response := models.LostPersonListResponse{
		Total:       len(responses),
		LostPersons: responses,
	}

	utils.SuccessResponse(c, http.StatusOK, "Search completed successfully", response)
}

// GetImage serves the uploaded image file
func (h *LostPersonHandler) GetImage(c *gin.Context) {
	imagePath := c.Param("path")
	if imagePath == "" {
		utils.BadRequest(c, "Image path is required")
		return
	}

	// Validate path to prevent directory traversal
	if strings.Contains(imagePath, "..") || strings.Contains(imagePath, "/") || strings.Contains(imagePath, "\\") {
		utils.BadRequest(c, "Invalid image path")
		return
	}

	fullPath := filepath.Join("uploads", imagePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		utils.NotFound(c, "Image not found")
		return
	}
	// Serve the file
	c.File(fullPath)
}
