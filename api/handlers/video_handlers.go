package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"video-processing-backend/models"

	"github.com/gin-gonic/gin"
)

var searchHistory *models.SearchHistory

// VideoUploadResponse represents the response structure
type VideoUploadResponse struct {
	UniqueFacesCount int      `json:"unique_faces_count"`
	Faces            []string `json:"faces"`
	Message          string   `json:"message"`
	ProcessingTime   float64  `json:"processing_time_seconds"`
}

// FaceSearchResponse represents the face search response structure
type FaceSearchResponse struct {
	Matches []FaceMatch `json:"matches"`
	Message string      `json:"message"`
}

// FaceMatch represents a match found in a video
type FaceMatch struct {
	Video        *models.VideoRecord `json:"video"`
	MatchedFaces []string            `json:"matched_faces"`
	Similarity   float64             `json:"similarity"`
}

// UploadVideoHandler handles video upload and processing
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

	// Get location information from form data
	locationName := c.PostForm("location_name")
	latitudeStr := c.PostForm("latitude")
	longitudeStr := c.PostForm("longitude")

	// Parse latitude and longitude
	var latitude, longitude float64
	var err1, err2 error

	if latitudeStr != "" {
		latitude, err1 = strconv.ParseFloat(latitudeStr, 64)
		if err1 != nil {
			log.Printf("Warning: Invalid latitude value: %s", latitudeStr)
		}
	}

	if longitudeStr != "" {
		longitude, err2 = strconv.ParseFloat(longitudeStr, 64)
		if err2 != nil {
			log.Printf("Warning: Invalid longitude value: %s", longitudeStr)
		}
	}

	// Create unique ID and filename
	videoID := fmt.Sprintf("video_%d", time.Now().Unix())
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename))
	videoPath := filepath.Join("../storage/videos", filename)

	// Create video record
	videoRecord := &models.VideoRecord{
		ID:               videoID,
		OriginalFilename: file.Filename,
		StoredPath:       videoPath,
		UploadTime:       time.Now(),
		Status:           "processing",
		LocationName:     locationName,
		Latitude:         latitude,
		Longitude:        longitude,
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

	log.Printf("Video saved: %s (Location: %s, Lat: %f, Lon: %f)",
		videoPath, locationName, latitude, longitude)

	// Process video with Python script
	response, err := processVideoWithPython(videoPath, videoID)
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

// SearchByFaceHandler handles face search functionality
func SearchByFaceHandler(c *gin.Context) {
	// Get the uploaded search image
	file, err := c.FormFile("search_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No search image provided",
		})
		return
	}

	// Validate file type
	if !isValidImageFile(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid image file format. Supported formats: jpg, jpeg, png",
		})
		return
	}

	// Save the search image temporarily
	searchImagePath := filepath.Join("../storage/temp", fmt.Sprintf("search_%d_%s", time.Now().Unix(), filepath.Base(file.Filename)))

	// Create temp directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(searchImagePath), 0755); err != nil {
		log.Printf("Error creating temp directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create temporary directory",
		})
		return
	}

	if err := c.SaveUploadedFile(file, searchImagePath); err != nil {
		log.Printf("Error saving search image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save search image",
		})
		return
	}

	// Get all videos with faces
	storage := GetVideoStorage()
	allVideos := storage.ListRecords()

	matches := []FaceMatch{} // Initialize as empty slice, not nil

	// Search through each video's faces
	log.Printf("Searching through %d videos", len(allVideos))
	for _, video := range allVideos {
		log.Printf("Checking video %s: status=%s, faces=%d", video.ID, video.Status, len(video.FaceImages))
		if video.Status == "completed" && len(video.FaceImages) > 0 {
			// Compare search image with faces in this video
			matchedFaces, err := compareFacesWithSearchImage(searchImagePath, video.FaceImages)
			if err != nil {
				log.Printf("Error comparing faces for video %s: %v", video.ID, err)
				continue
			}

			log.Printf("Video %s: found %d matched faces", video.ID, len(matchedFaces))
			if len(matchedFaces) > 0 {
				matches = append(matches, FaceMatch{
					Video:        video,
					MatchedFaces: matchedFaces,
					Similarity:   0.85, // Default similarity score
				})
			}
		}
	}

	// Clean up temporary search image
	defer os.Remove(searchImagePath)

	// Add debug logging
	log.Printf("Search completed. Found %d matches", len(matches))
	for i, match := range matches {
		log.Printf("Match %d: Video %s, %d matched faces", i+1, match.Video.ID, len(match.MatchedFaces))
	}

	response := FaceSearchResponse{
		Matches: matches,
		Message: fmt.Sprintf("Found %d video(s) with matching faces", len(matches)),
	}

	// Ensure matches is always an array, not null
	if response.Matches == nil {
		response.Matches = []FaceMatch{}
	}

	// Debug: Print the response structure
	responseJSON, _ := json.Marshal(response)
	log.Printf("Response JSON: %s", string(responseJSON))

	c.JSON(http.StatusOK, response)
}

// HealthCheckHandler provides a simple health check endpoint
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// processVideoWithPython calls the Python script to process the video
func processVideoWithPython(videoPath string, videoID string) (*VideoUploadResponse, error) {
	// Get the absolute path to the Python script
	pythonScriptPath := filepath.Join("python", "face_detect.py")

	// Ensure the script exists
	if _, err := os.Stat(pythonScriptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Python script not found: %s", pythonScriptPath)
	}

	// Execute Python script with virtual environment and video ID
	cmd := exec.Command("venv/bin/python3", pythonScriptPath, videoPath, "--video-id", videoID)
	cmd.Dir = "." // Set working directory to api root

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

// compareFacesWithSearchImage compares a search image with stored face images
func compareFacesWithSearchImage(searchImagePath string, faceImages []string) ([]string, error) {
	// Get the absolute path to the Python script
	pythonScriptPath := filepath.Join("python", "face_search.py")

	// Ensure the script exists
	if _, err := os.Stat(pythonScriptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Python face search script not found: %s", pythonScriptPath)
	}

	// Prepare face images list
	faceImagesStr := strings.Join(faceImages, ",")

	// Execute Python script for face comparison
	cmd := exec.Command("venv/bin/python3", pythonScriptPath, searchImagePath, "--face-images", faceImagesStr)
	cmd.Dir = "." // Set working directory to api root

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Face search Python script error: %v", err)
		log.Printf("Face search Python output: %s", string(output))
		return nil, fmt.Errorf("face search script execution failed: %v", err)
	}

	// Parse JSON response
	var result struct {
		MatchedFaces []string `json:"matched_faces"`
		Error        string   `json:"error,omitempty"`
	}

	outputStr := string(output)
	lastBraceIndex := strings.LastIndex(outputStr, "}")
	if lastBraceIndex != -1 {
		startIndex := strings.LastIndex(outputStr[:lastBraceIndex+1], "{")
		if startIndex != -1 {
			jsonStr := outputStr[startIndex : lastBraceIndex+1]
			if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
				log.Printf("Failed to parse face search output: %s", jsonStr)
				return nil, fmt.Errorf("failed to parse face search output: %v", err)
			}
		}
	}

	if result.Error != "" {
		return nil, fmt.Errorf("face search error: %s", result.Error)
	}

	return result.MatchedFaces, nil
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

// isValidImageFile checks if the uploaded file is a valid image format
func isValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// generateImageHash generates an MD5 hash of an image file
func generateImageHash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
