package handlers

import (
	"net/http"
	"strconv"

	"video-analysis-service/internal/models"
	"video-analysis-service/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VideoHandler handles video-related HTTP requests
type VideoHandler struct {
	videoService *services.VideoService
	log          *zap.Logger
}

// NewVideoHandler creates a new video handler
func NewVideoHandler(videoService *services.VideoService, log *zap.Logger) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
		log:          log,
	}
}

// UploadVideo godoc
// @Summary Upload a video file
// @Description Upload a video file for analysis
// @Tags videos
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Video file to upload"
// @Success 200 {object} models.UploadVideoResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /videos/upload [post]
func (h *VideoHandler) UploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		h.log.Error("Failed to get uploaded file", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_file",
			Message:   "No file provided or invalid file",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	video, err := h.videoService.UploadVideo(file)
	if err != nil {
		h.log.Error("Failed to upload video", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "upload_failed",
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.UploadVideoResponse{
		Success: true,
		Message: "Video uploaded successfully",
		Video:   *video,
	})
}

// ListVideos godoc
// @Summary List all videos
// @Description Get a list of all uploaded videos with pagination
// @Tags videos
// @Produce json
// @Param limit query int false "Number of videos to return (default: 10, max: 100)"
// @Param offset query int false "Number of videos to skip (default: 0)"
// @Success 200 {object} models.ListVideosResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /videos [get]
func (h *VideoHandler) ListVideos(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	videos, total, err := h.videoService.ListVideos(limit, offset)
	if err != nil {
		h.log.Error("Failed to list videos", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "list_failed",
			Message:   "Failed to retrieve videos",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	if videos == nil {
		videos = []models.Video{}
	}
	c.JSON(http.StatusOK, models.ListVideosResponse{
		Success: true,
		Videos:  videos,
		Total:   total,
	})
}

// GetVideo godoc
// @Summary Get video details
// @Description Get detailed information about a specific video
// @Tags videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} models.Video
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /videos/{id} [get]
func (h *VideoHandler) GetVideo(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_id",
			Message:   "Video ID is required",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	video, err := h.videoService.GetVideo(videoID)
	if err != nil {
		h.log.Error("Failed to get video", zap.String("video_id", videoID), zap.Error(err))
		if err.Error() == "video not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success:   false,
				Error:     "not_found",
				Message:   "Video not found",
				RequestID: c.GetString("request_id"),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success:   false,
				Error:     "get_failed",
				Message:   "Failed to retrieve video",
				RequestID: c.GetString("request_id"),
			})
		}
		return
	}

	c.JSON(http.StatusOK, video)
}

// DeleteVideo godoc
// @Summary Delete a video
// @Description Delete a video and its associated analysis results
// @Tags videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /videos/{id} [delete]
func (h *VideoHandler) DeleteVideo(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_id",
			Message:   "Video ID is required",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	err := h.videoService.DeleteVideo(videoID)
	if err != nil {
		h.log.Error("Failed to delete video", zap.String("video_id", videoID), zap.Error(err))
		if err.Error() == "video not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success:   false,
				Error:     "not_found",
				Message:   "Video not found",
				RequestID: c.GetString("request_id"),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success:   false,
				Error:     "delete_failed",
				Message:   "Failed to delete video",
				RequestID: c.GetString("request_id"),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.ErrorResponse{
		Success:   true,
		Message:   "Video deleted successfully",
		RequestID: c.GetString("request_id"),
	})
}

// DownloadVideo godoc
// @Summary Download a video file
// @Description Download the original video file
// @Tags videos
// @Produce octet-stream
// @Param id path string true "Video ID"
// @Success 200 {file} file
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /videos/{id}/download [get]
func (h *VideoHandler) DownloadVideo(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_id",
			Message:   "Video ID is required",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	video, err := h.videoService.GetVideo(videoID)
	if err != nil {
		h.log.Error("Failed to get video for download", zap.String("video_id", videoID), zap.Error(err))
		if err.Error() == "video not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success:   false,
				Error:     "not_found",
				Message:   "Video not found",
				RequestID: c.GetString("request_id"),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success:   false,
				Error:     "download_failed",
				Message:   "Failed to retrieve video",
				RequestID: c.GetString("request_id"),
			})
		}
		return
	}

	filePath := h.videoService.GetVideoFilePath(video.Filename)
	c.FileAttachment(filePath, video.OriginalFilename)
}
