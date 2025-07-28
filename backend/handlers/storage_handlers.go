package handlers

import (
	"net/http"
	"strconv"

	"video-processing-backend/models"

	"github.com/gin-gonic/gin"
)

var videoStorage *models.VideoStorage

// InitializeStorage initializes the video storage system
func InitializeStorage() {
	videoStorage = models.NewVideoStorage("storage/videos.json")
	if err := videoStorage.Load(); err != nil {
		panic("Failed to load video storage: " + err.Error())
	}
}

// GetVideoStorage returns the video storage instance
func GetVideoStorage() *models.VideoStorage {
	return videoStorage
}

// ListVideosHandler returns all stored video records
func ListVideosHandler(c *gin.Context) {
	records := videoStorage.GetAllRecords()

	c.JSON(http.StatusOK, gin.H{
		"videos": records,
		"count":  len(records),
	})
}

// GetVideoHandler returns a specific video record
func GetVideoHandler(c *gin.Context) {
	id := c.Param("id")

	record, exists := videoStorage.GetRecord(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video not found",
		})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteVideoHandler deletes a video record and associated files
func DeleteVideoHandler(c *gin.Context) {
	id := c.Param("id")

	if err := videoStorage.DeleteRecord(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video deleted successfully",
	})
}

// GetVideoStatsHandler returns statistics about stored videos
func GetVideoStatsHandler(c *gin.Context) {
	records := videoStorage.GetAllRecords()

	totalVideos := len(records)
	totalFaces := 0
	totalProcessingTime := 0.0
	completedVideos := 0
	failedVideos := 0

	for _, record := range records {
		if record.Status == "completed" {
			completedVideos++
			totalFaces += record.UniqueFacesCount
			totalProcessingTime += record.ProcessingTime
		} else if record.Status == "failed" {
			failedVideos++
		}
	}

	avgProcessingTime := 0.0
	if completedVideos > 0 {
		avgProcessingTime = totalProcessingTime / float64(completedVideos)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_videos":          totalVideos,
		"completed_videos":      completedVideos,
		"failed_videos":         failedVideos,
		"total_faces_detected":  totalFaces,
		"avg_processing_time":   avgProcessingTime,
		"total_processing_time": totalProcessingTime,
	})
}

// CleanupOldVideosHandler removes videos older than specified days
func CleanupOldVideosHandler(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid days parameter",
		})
		return
	}

	if err := videoStorage.CleanupOldRecords(days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cleanup completed successfully",
		"days":    days,
	})
}
