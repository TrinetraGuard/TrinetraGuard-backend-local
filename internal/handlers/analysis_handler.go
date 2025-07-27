package handlers

import (
	"net/http"

	"video-analysis-service/internal/models"
	"video-analysis-service/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AnalysisHandler handles analysis-related HTTP requests
type AnalysisHandler struct {
	analysisService *services.AnalysisService
	log             *zap.Logger
}

// NewAnalysisHandler creates a new analysis handler
func NewAnalysisHandler(analysisService *services.AnalysisService, log *zap.Logger) *AnalysisHandler {
	return &AnalysisHandler{
		analysisService: analysisService,
		log:             log,
	}
}

// StartAnalysis godoc
// @Summary Start video analysis
// @Description Start analyzing a video for person detection and tracking
// @Tags analysis
// @Accept json
// @Produce json
// @Param videoId path string true "Video ID"
// @Success 200 {object} models.AnalysisStatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analysis/{videoId}/start [post]
func (h *AnalysisHandler) StartAnalysis(c *gin.Context) {
	videoID := c.Param("videoId")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_video_id",
			Message:   "Video ID is required",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	job, err := h.analysisService.StartAnalysis(videoID)
	if err != nil {
		h.log.Error("Failed to start analysis", zap.String("video_id", videoID), zap.Error(err))

		if err.Error() == "video not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success:   false,
				Error:     "not_found",
				Message:   "Video not found",
				RequestID: c.GetString("request_id"),
			})
		} else if err.Error() == "analysis job already exists for this video" {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Success:   false,
				Error:     "job_exists",
				Message:   "Analysis job already exists for this video",
				RequestID: c.GetString("request_id"),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success:   false,
				Error:     "start_failed",
				Message:   "Failed to start analysis",
				RequestID: c.GetString("request_id"),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.AnalysisStatusResponse{
		Success: true,
		Job:     *job,
	})
}

// GetAnalysisStatus godoc
// @Summary Get analysis status
// @Description Get the current status of a video analysis job
// @Tags analysis
// @Produce json
// @Param videoId path string true "Video ID"
// @Success 200 {object} models.AnalysisStatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analysis/{videoId}/status [get]
func (h *AnalysisHandler) GetAnalysisStatus(c *gin.Context) {
	videoID := c.Param("videoId")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_video_id",
			Message:   "Video ID is required",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Get the latest analysis job for this video
	// Note: This is a simplified implementation. In a real scenario, you might want to
	// get the job ID from the URL or query parameters
	job, err := h.analysisService.GetAnalysisStatus(videoID)
	if err != nil {
		h.log.Error("Failed to get analysis status", zap.String("video_id", videoID), zap.Error(err))
		if err.Error() == "analysis job not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success:   false,
				Error:     "not_found",
				Message:   "Analysis job not found",
				RequestID: c.GetString("request_id"),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success:   false,
				Error:     "status_failed",
				Message:   "Failed to get analysis status",
				RequestID: c.GetString("request_id"),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.AnalysisStatusResponse{
		Success: true,
		Job:     *job,
	})
}

// GetAnalysisResults godoc
// @Summary Get analysis results
// @Description Get the results of a completed video analysis
// @Tags analysis
// @Produce json
// @Param videoId path string true "Video ID"
// @Success 200 {object} models.AnalysisResultsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analysis/{videoId}/results [get]
func (h *AnalysisHandler) GetAnalysisResults(c *gin.Context) {
	videoID := c.Param("videoId")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_video_id",
			Message:   "Video ID is required",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	results, err := h.analysisService.GetAnalysisResults(videoID)
	if err != nil {
		h.log.Error("Failed to get analysis results", zap.String("video_id", videoID), zap.Error(err))
		if err.Error() == "analysis results not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success:   false,
				Error:     "not_found",
				Message:   "Analysis results not found",
				RequestID: c.GetString("request_id"),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success:   false,
				Error:     "results_failed",
				Message:   "Failed to get analysis results",
				RequestID: c.GetString("request_id"),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.AnalysisResultsResponse{
		Success: true,
		Results: *results,
	})
}

// StartBatchAnalysis godoc
// @Summary Start batch analysis
// @Description Start analysis for multiple videos
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body []string true "Array of video IDs"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analysis/batch [post]
func (h *AnalysisHandler) StartBatchAnalysis(c *gin.Context) {
	var videoIDs []string
	if err := c.ShouldBindJSON(&videoIDs); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_request",
			Message:   "Invalid request body",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	if len(videoIDs) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "empty_list",
			Message:   "Video IDs list cannot be empty",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	jobs, err := h.analysisService.StartBatchAnalysis(videoIDs)
	if err != nil {
		h.log.Error("Failed to start batch analysis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "batch_failed",
			Message:   "Failed to start batch analysis",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Batch analysis started",
		"jobs":    jobs,
		"total":   len(jobs),
	})
}
