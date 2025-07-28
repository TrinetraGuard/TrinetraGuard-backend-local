package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"video-processing-backend/models"

	"github.com/gin-gonic/gin"
)

var videoStorage *models.VideoStorage

// InitializeStorage initializes the video storage system
func InitializeStorage() {
	videoStorage = models.NewVideoStorage("../storage/data/videos.json")
	if err := videoStorage.Load(); err != nil {
		panic("Failed to load video storage: " + err.Error())
	}

	searchHistory = models.NewSearchHistory("../storage/data/search_history.json")
	if err := searchHistory.Load(); err != nil {
		log.Printf("Warning: Failed to load search history: %v", err)
	}
}

// GetVideoStorage returns the video storage instance
func GetVideoStorage() *models.VideoStorage {
	return videoStorage
}

// ListVideosHandler returns all video records (active and archived)
func ListVideosHandler(c *gin.Context) {
	records := videoStorage.ListRecords()
	c.JSON(http.StatusOK, gin.H{
		"videos": records,
		"count":  len(records),
	})
}

// ListActiveVideosHandler returns only active video records
func ListActiveVideosHandler(c *gin.Context) {
	records := videoStorage.ListActiveRecords()
	c.JSON(http.StatusOK, gin.H{
		"videos": records,
		"count":  len(records),
		"type":   "active",
	})
}

// ListArchivedVideosHandler returns only archived video records (history)
func ListArchivedVideosHandler(c *gin.Context) {
	records := videoStorage.ListArchivedRecords()
	c.JSON(http.StatusOK, gin.H{
		"videos": records,
		"count":  len(records),
		"type":   "archived",
	})
}

// GetVideoHandler returns a specific video record
func GetVideoHandler(c *gin.Context) {
	id := c.Param("id")
	record, exists := videoStorage.GetRecord(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video record not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"video": record,
	})
}

// DeleteVideoHandler archives a video record (moves to history)
func DeleteVideoHandler(c *gin.Context) {
	id := c.Param("id")

	if err := videoStorage.DeleteRecord(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video record not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video moved to history successfully",
		"id":      id,
	})
}

// RestoreVideoHandler restores an archived video record
func RestoreVideoHandler(c *gin.Context) {
	id := c.Param("id")
	record, exists := videoStorage.GetRecord(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video record not found",
		})
		return
	}

	if !record.IsArchived {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Video is not archived",
		})
		return
	}

	// Restore the record
	record.IsArchived = false
	record.LastAccessed = time.Now()

	if err := videoStorage.UpdateRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to restore video",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video restored successfully",
		"id":      id,
	})
}

// GetVideoStatsHandler returns storage statistics
func GetVideoStatsHandler(c *gin.Context) {
	stats := videoStorage.GetStats()
	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// CleanupOldVideosHandler removes very old archived records (optional)
func CleanupOldVideosHandler(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid days parameter",
		})
		return
	}

	if days < 7 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Minimum cleanup period is 7 days",
		})
		return
	}

	if err := videoStorage.CleanupOldRecords(days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to cleanup old records",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cleanup completed successfully",
		"days":    days,
	})
}

// SearchVideosHandler searches videos by filename or status
func SearchVideosHandler(c *gin.Context) {
	query := c.Query("q")
	status := c.Query("status")
	archived := c.Query("archived")

	var records []*models.VideoRecord

	if archived == "true" {
		records = videoStorage.ListArchivedRecords()
	} else if archived == "false" {
		records = videoStorage.ListActiveRecords()
	} else {
		records = videoStorage.ListRecords()
	}

	// Filter by query if provided
	if query != "" {
		var filtered []*models.VideoRecord
		for _, record := range records {
			if contains(record.OriginalFilename, query) ||
				contains(record.Status, query) ||
				contains(record.ID, query) {
				filtered = append(filtered, record)
			}
		}
		records = filtered
	}

	// Filter by status if provided
	if status != "" {
		var filtered []*models.VideoRecord
		for _, record := range records {
			if record.Status == status {
				filtered = append(filtered, record)
			}
		}
		records = filtered
	}

	c.JSON(http.StatusOK, gin.H{
		"videos":   records,
		"count":    len(records),
		"query":    query,
		"status":   status,
		"archived": archived,
	})
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

// containsSubstring is a simple substring search
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ResetDatabaseHandler completely resets the database and removes all files
func ResetDatabaseHandler(c *gin.Context) {
	// Get confirmation from request - check both form data and query parameters
	confirm := c.PostForm("confirm")
	if confirm == "" {
		confirm = c.Query("confirm")
	}

	if confirm != "true" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please confirm by sending 'confirm=true'",
		})
		return
	}

	// Reset the storage
	if err := videoStorage.ResetDatabase(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reset database: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database reset successfully. All videos and faces have been removed.",
	})
}

// GetSearchHistoryHandler returns search history records
func GetSearchHistoryHandler(c *gin.Context) {
	if searchHistory == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search history not initialized",
		})
		return
	}

	records := searchHistory.ListRecords()
	c.JSON(http.StatusOK, gin.H{
		"searches": records,
		"count":    len(records),
	})
}

// GetSearchHistoryStatsHandler returns search history statistics
func GetSearchHistoryStatsHandler(c *gin.Context) {
	if searchHistory == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search history not initialized",
		})
		return
	}

	stats := searchHistory.GetStats()
	c.JSON(http.StatusOK, stats)
}

// GetVideoPreviewHandler returns video preview information
func GetVideoPreviewHandler(c *gin.Context) {
	id := c.Param("id")
	record, exists := videoStorage.GetRecord(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video record not found",
		})
		return
	}

	// Check if video file exists
	if _, err := os.Stat(record.StoredPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video file not found",
		})
		return
	}

	// Return video preview information
	c.JSON(http.StatusOK, gin.H{
		"video": gin.H{
			"id":                record.ID,
			"original_filename": record.OriginalFilename,
			"upload_time":       record.UploadTime,
			"status":            record.Status,
			"location_name":     record.LocationName,
			"latitude":          record.Latitude,
			"longitude":         record.Longitude,
			"unique_faces":      record.UniqueFacesCount,
			"processing_time":   record.ProcessingTime,
			"video_url":         fmt.Sprintf("/api/videos/%s/file", record.ID),
		},
	})
}

// GetVideoFileHandler serves the actual video file
func GetVideoFileHandler(c *gin.Context) {
	id := c.Param("id")
	record, exists := videoStorage.GetRecord(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video record not found",
		})
		return
	}

	// Check if video file exists
	if _, err := os.Stat(record.StoredPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video file not found",
		})
		return
	}

	// Serve the video file
	c.File(record.StoredPath)
}
