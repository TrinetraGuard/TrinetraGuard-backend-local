package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"video-processing-backend/models"

	"github.com/gin-gonic/gin"
)

// VideoUploadResponse represents the response structure
type VideoUploadResponse struct {
	UniqueFacesCount int      `json:"unique_faces_count"`
	Faces            []string `json:"faces"`
	Message          string   `json:"message"`
	ProcessingTime   float64  `json:"processing_time_seconds"`
}

// uploadVideoHandler handles video upload and processing
func UploadVideoHandler(c *gin.Context) {
	startTime := time.Now()

	// Get the uploaded file
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No video file provided",
		})
		return
	}

	// Validate file type
	if !isValidVideoFile(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid video file format. Supported formats: mp4, avi, mov, mkv",
		})
		return
	}

	// Create unique ID and filename
	videoID := fmt.Sprintf("video_%d", time.Now().Unix())
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename))
	videoPath := filepath.Join("videos", filename)

	// Create video record
	videoRecord := &models.VideoRecord{
		ID:               videoID,
		OriginalFilename: file.Filename,
		StoredPath:       videoPath,
		UploadTime:       time.Now(),
		Status:           "processing",
	}

	// Save the uploaded file
	if err := c.SaveUploadedFile(file, videoPath); err != nil {
		log.Printf("Error saving file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save video file",
		})
		return
	}

	// Save record to storage
	storage := GetVideoStorage()
	if err := storage.AddRecord(videoRecord); err != nil {
		log.Printf("Error saving video record: %v", err)
	}

	log.Printf("Video saved: %s", videoPath)

	// Process video with Python script
	response, err := processVideoWithPython(videoPath)
	if err != nil {
		log.Printf("Error processing video: %v", err)

		// Update record with error
		videoRecord.Status = "failed"
		videoRecord.ErrorMessage = err.Error()
		storage.UpdateRecord(videoRecord)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process video",
		})
		return
	}

	// Calculate processing time
	processingTime := time.Since(startTime).Seconds()
	response.ProcessingTime = processingTime

	// Update record with results
	videoRecord.Status = "completed"
	videoRecord.ProcessingTime = processingTime
	videoRecord.UniqueFacesCount = response.UniqueFacesCount
	videoRecord.FaceImages = response.Faces
	storage.UpdateRecord(videoRecord)

	c.JSON(http.StatusOK, response)
}

// healthCheckHandler provides a simple health check endpoint
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// processVideoWithPython calls the Python script to process the video
func processVideoWithPython(videoPath string) (*VideoUploadResponse, error) {
	// Get the absolute path to the Python script
	pythonScriptPath := filepath.Join("python", "face_detect.py")

	// Ensure the script exists
	if _, err := os.Stat(pythonScriptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Python script not found: %s", pythonScriptPath)
	}

	// Execute Python script
	cmd := exec.Command("python3", pythonScriptPath, videoPath)
	cmd.Dir = "." // Set working directory to backend root

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Python script error: %v", err)
		log.Printf("Python output: %s", string(output))
		return nil, fmt.Errorf("Python script execution failed: %v", err)
	}

	// Parse JSON response from Python script
	var response VideoUploadResponse

	// Clean the output by finding the last JSON object
	outputStr := string(output)
	lastBraceIndex := strings.LastIndex(outputStr, "}")
	if lastBraceIndex != -1 {
		// Find the start of the JSON object
		startIndex := strings.LastIndex(outputStr[:lastBraceIndex+1], "{")
		if startIndex != -1 {
			jsonStr := outputStr[startIndex : lastBraceIndex+1]
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				log.Printf("Failed to parse Python output: %s", jsonStr)
				return nil, fmt.Errorf("failed to parse Python script output: %v", err)
			}
		} else {
			return nil, fmt.Errorf("no valid JSON found in Python output")
		}
	} else {
		return nil, fmt.Errorf("no JSON object found in Python output")
	}

	return &response, nil
}

// isValidVideoFile checks if the uploaded file is a valid video format
func isValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := []string{".mp4", ".avi", ".mov", ".mkv", ".wmv", ".flv", ".webm"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}
