package handlers

import (
	"net/http"

	"video-analysis-service/internal/models"
	"video-analysis-service/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FinderHandler handles finder-related HTTP requests
type FinderHandler struct {
	finderService *services.FinderService
	log           *zap.Logger
}

// NewFinderHandler creates a new finder handler
func NewFinderHandler(finderService *services.FinderService, log *zap.Logger) *FinderHandler {
	return &FinderHandler{
		finderService: finderService,
		log:           log,
	}
}

// UploadReferenceImage godoc
// @Summary Upload a reference image
// @Description Upload a reference image for person finding
// @Tags finder
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Reference image file"
// @Success 200 {object} models.UploadVideoResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/images/upload [post]
func (h *FinderHandler) UploadReferenceImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		h.log.Error("Failed to get uploaded file", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid_file",
			Message: "No file provided or invalid file",
		})
		return
	}

	image, err := h.finderService.UploadReferenceImage(file)
	if err != nil {
		h.log.Error("Failed to upload reference image", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "upload_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Reference image uploaded successfully",
		"image":   image,
	})
}

// ListReferenceImages godoc
// @Summary List reference images
// @Description Get a list of all uploaded reference images
// @Tags finder
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/images [get]
func (h *FinderHandler) ListReferenceImages(c *gin.Context) {
	images, err := h.finderService.ListReferenceImages()
	if err != nil {
		h.log.Error("Failed to list reference images", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "list_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"images":  images,
	})
}

// DeleteReferenceImage godoc
// @Summary Delete a reference image
// @Description Delete a reference image by ID
// @Tags finder
// @Produce json
// @Param id path string true "Image ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/images/{id} [delete]
func (h *FinderHandler) DeleteReferenceImage(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		h.log.Error("Missing image ID parameter")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "missing_image_id",
			Message: "Image ID is required",
		})
		return
	}

	err := h.finderService.DeleteReferenceImage(imageID)
	if err != nil {
		h.log.Error("Failed to delete reference image", zap.String("image_id", imageID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Reference image deleted successfully",
	})
}

// SearchPerson godoc
// @Summary Search for a person
// @Description Start a search job to find a person in videos using a reference image
// @Tags finder
// @Accept json
// @Produce json
// @Param request body models.SearchPersonRequest true "Search request"
// @Success 200 {object} models.SearchPersonResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/search [post]
func (h *FinderHandler) SearchPerson(c *gin.Context) {
	var request models.SearchPersonRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Error("Failed to bind request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	job, err := h.finderService.SearchPerson(request.ReferenceImageID)
	if err != nil {
		h.log.Error("Failed to start person search", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "search_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SearchPersonResponse{
		Success: true,
		Message: "Person search started successfully",
		Job:     *job,
	})
}

// GetSearchStatus godoc
// @Summary Get search status
// @Description Get the status of a person search job
// @Tags finder
// @Produce json
// @Param jobId path string true "Search job ID"
// @Success 200 {object} models.AnalysisStatusResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/search/{jobId}/status [get]
func (h *FinderHandler) GetSearchStatus(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		h.log.Error("Missing job ID parameter")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "missing_job_id",
			Message: "Job ID is required",
		})
		return
	}

	job, err := h.finderService.GetSearchStatus(jobID)
	if err != nil {
		h.log.Error("Failed to get search status", zap.String("job_id", jobID), zap.Error(err))
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "search_job_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SearchStatusResponse{
		Success: true,
		Job:     *job,
	})
}

// GetSearchResults godoc
// @Summary Get search results
// @Description Get the results of a completed person search
// @Tags finder
// @Produce json
// @Param jobId path string true "Search job ID"
// @Success 200 {object} models.SearchResultsResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/search/{jobId}/results [get]
func (h *FinderHandler) GetSearchResults(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		h.log.Error("Missing job ID parameter")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "missing_job_id",
			Message: "Job ID is required",
		})
		return
	}

	results, err := h.finderService.GetSearchResults(jobID)
	if err != nil {
		h.log.Error("Failed to get search results", zap.String("job_id", jobID), zap.Error(err))
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "search_results_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SearchResultsResponse{
		Success: true,
		Results: *results,
	})
}
