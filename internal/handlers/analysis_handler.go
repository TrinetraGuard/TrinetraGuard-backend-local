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
// @Description Start analysis for a specific video
// @Tags analysis
// @Accept json
// @Produce json
// @Param videoId path string true "Video ID"
// @Success 200 {object} models.AnalysisStatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analysis/{videoId}/start [post]
func (h *AnalysisHandler) StartAnalysis(c *gin.Context) {
	videoID := c.Param("videoId")
	if videoID == "" {
		h.log.Error("Missing video ID parameter")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "missing_video_id",
			Message: "Video ID is required",
		})
		return
	}

	job, err := h.analysisService.StartAnalysis(videoID)
	if err != nil {
		h.log.Error("Failed to start analysis", zap.String("video_id", videoID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "analysis_start_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.AnalysisStatusResponse{
		Success: true,
		Job:     *job,
	})
}

// GetAnalysisStatus godoc
// @Summary Get analysis status
// @Description Get the status of analysis for a specific video
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
		h.log.Error("Missing video ID parameter")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "missing_video_id",
			Message: "Video ID is required",
		})
		return
	}

	job, err := h.analysisService.GetAnalysisStatus(videoID)
	if err != nil {
		h.log.Error("Failed to get analysis status", zap.String("video_id", videoID), zap.Error(err))
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "analysis_job_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.AnalysisStatusResponse{
		Success: true,
		Job:     *job,
	})
}

// GetAnalysisResults godoc
// @Summary Get analysis results
// @Description Get the analysis results for a specific video
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
		h.log.Error("Missing video ID parameter")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "missing_video_id",
			Message: "Video ID is required",
		})
		return
	}

	results, err := h.analysisService.GetAnalysisResults(videoID)
	if err != nil {
		h.log.Error("Failed to get analysis results", zap.String("video_id", videoID), zap.Error(err))
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "analysis_results_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.AnalysisResultsResponse{
		Success: true,
		Results: *results,
	})
}

// GetEnhancedAnalysisResults godoc
// @Summary Get enhanced analysis results
// @Description Get enhanced analysis results with person and face data for a specific video
// @Tags analysis
// @Produce json
// @Param videoId path string true "Video ID"
// @Success 200 {object} models.EnhancedAnalysisResultsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analysis/{videoId}/enhanced-results [get]
func (h *AnalysisHandler) GetEnhancedAnalysisResults(c *gin.Context) {
	videoID := c.Param("videoId")
	if videoID == "" {
		h.log.Error("Missing video ID parameter")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "missing_video_id",
			Message: "Video ID is required",
		})
		return
	}

	results, err := h.analysisService.GetEnhancedAnalysisResults(videoID)
	if err != nil {
		h.log.Error("Failed to get enhanced analysis results", zap.String("video_id", videoID), zap.Error(err))
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "enhanced_analysis_results_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.EnhancedAnalysisResultsResponse{
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
// @Success 200 {array} models.AnalysisJob
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analysis/batch [post]
func (h *AnalysisHandler) StartBatchAnalysis(c *gin.Context) {
	var videoIDs []string
	if err := c.ShouldBindJSON(&videoIDs); err != nil {
		h.log.Error("Failed to bind request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid_request_body",
			Message: "Invalid request body",
		})
		return
	}

	if len(videoIDs) == 0 {
		h.log.Error("Empty video IDs array")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "empty_video_ids",
			Message: "At least one video ID is required",
		})
		return
	}

	jobs, err := h.analysisService.StartBatchAnalysis(videoIDs)
	if err != nil {
		h.log.Error("Failed to start batch analysis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "batch_analysis_failed",
			Message: err.Error(),
		})
		return
	}

	// Convert to array of AnalysisJob
	var jobArray []models.AnalysisJob
	for _, job := range jobs {
		jobArray = append(jobArray, *job)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"jobs":    jobArray,
	})
}
