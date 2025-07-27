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
// @Param description formData string false "Description of the person"
// @Success 200 {object} models.ReferenceImage
// @Failure 400 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/upload [post]
func (h *FinderHandler) UploadReferenceImage(c *gin.Context) {
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

	description := c.PostForm("description")

	refImage, err := h.finderService.UploadReferenceImage(file, description)
	if err != nil {
		h.log.Error("Failed to upload reference image", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "upload_failed",
			Message:   err.Error(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, refImage)
}

// ListReferenceImages godoc
// @Summary List reference images
// @Description Get a list of all uploaded reference images
// @Tags finder
// @Produce json
// @Success 200 {array} models.ReferenceImage
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/images [get]
func (h *FinderHandler) ListReferenceImages(c *gin.Context) {
	images, err := h.finderService.ListReferenceImages()
	if err != nil {
		h.log.Error("Failed to list reference images", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "list_failed",
			Message:   "Failed to retrieve reference images",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	if images == nil {
		images = []models.ReferenceImage{}
	}
	c.JSON(http.StatusOK, images)
}

// SearchPerson godoc
// @Summary Search for a person
// @Description Search for a person in videos using a reference image
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
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_request",
			Message:   "Invalid request body",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	searchJob, err := h.finderService.SearchPerson(&request)
	if err != nil {
		h.log.Error("Failed to start person search", zap.Error(err))

		if err.Error() == "reference image not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success:   false,
				Error:     "not_found",
				Message:   "Reference image not found",
				RequestID: c.GetString("request_id"),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success:   false,
				Error:     "search_failed",
				Message:   "Failed to start person search",
				RequestID: c.GetString("request_id"),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.SearchPersonResponse{
		Success:   true,
		Message:   "Person search started successfully",
		SearchJob: *searchJob,
	})
}

// GetSearchStatus godoc
// @Summary Get search status
// @Description Get the current status of a person search job
// @Tags finder
// @Produce json
// @Param searchId path string true "Search Job ID"
// @Success 200 {object} models.SearchJob
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/search/{searchId}/status [get]
func (h *FinderHandler) GetSearchStatus(c *gin.Context) {
	searchID := c.Param("searchId")
	if searchID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_search_id",
			Message:   "Search ID is required",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	searchJob, err := h.finderService.GetSearchStatus(searchID)
	if err != nil {
		h.log.Error("Failed to get search status", zap.String("search_id", searchID), zap.Error(err))
		if err.Error() == "search job not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success:   false,
				Error:     "not_found",
				Message:   "Search job not found",
				RequestID: c.GetString("request_id"),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success:   false,
				Error:     "status_failed",
				Message:   "Failed to get search status",
				RequestID: c.GetString("request_id"),
			})
		}
		return
	}

	c.JSON(http.StatusOK, searchJob)
}

// GetSearchResults godoc
// @Summary Get search results
// @Description Get the results of a completed person search
// @Tags finder
// @Produce json
// @Param searchId path string true "Search Job ID"
// @Success 200 {object} models.SearchResultsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/search/{searchId}/results [get]
func (h *FinderHandler) GetSearchResults(c *gin.Context) {
	searchID := c.Param("searchId")
	if searchID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_search_id",
			Message:   "Search ID is required",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	results, err := h.finderService.GetSearchResults(searchID)
	if err != nil {
		h.log.Error("Failed to get search results", zap.String("search_id", searchID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "results_failed",
			Message:   "Failed to get search results",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.SearchResultsResponse{
		Success: true,
		Results: results,
		Total:   len(results),
	})
}

// DeleteReferenceImage godoc
// @Summary Delete a reference image
// @Description Delete a reference image and its associated search results
// @Tags finder
// @Produce json
// @Param id path string true "Reference Image ID"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /finder/images/{id} [delete]
func (h *FinderHandler) DeleteReferenceImage(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "invalid_id",
			Message:   "Reference image ID is required",
			RequestID: c.GetString("request_id"),
		})
		return
	}

	err := h.finderService.DeleteReferenceImage(imageID)
	if err != nil {
		h.log.Error("Failed to delete reference image", zap.String("image_id", imageID), zap.Error(err))
		if err.Error() == "reference image not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success:   false,
				Error:     "not_found",
				Message:   "Reference image not found",
				RequestID: c.GetString("request_id"),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success:   false,
				Error:     "delete_failed",
				Message:   "Failed to delete reference image",
				RequestID: c.GetString("request_id"),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.ErrorResponse{
		Success:   true,
		Message:   "Reference image deleted successfully",
		RequestID: c.GetString("request_id"),
	})
}
