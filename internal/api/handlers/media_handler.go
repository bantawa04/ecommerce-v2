package handlers

import (
	"net/http"
	"strconv"

	"beautyessentials.com/internal/api/responses"
	"beautyessentials.com/internal/requests"
	"beautyessentials.com/internal/service/interfaces"
	"beautyessentials.com/internal/validators"
	"github.com/gin-gonic/gin"
)

// MediaHandler handles media-related requests
type MediaHandler struct {
	mediaService interfaces.MediaService
	respHelper   *responses.ResponseHelper
	validator    *validators.Validator
}

// NewMediaHandler creates a new instance of MediaHandler
func NewMediaHandler(
	mediaService interfaces.MediaService,
	respHelper *responses.ResponseHelper,
) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
		respHelper:   respHelper,
		validator:    validators.NewValidator(),
	}
}

// GetAllMedia handles the request to get all media
func (h *MediaHandler) GetAllMedia(c *gin.Context) {
	// Extract query parameters for filtering and pagination
	filters := make(map[string]interface{})
	appends := make(map[string]interface{})

	// Add trashed filter if provided
	if trashed := c.Query("trashed"); trashed != "" {
		filters["trashed"] = trashed == "true"
	}

	// Add pagination parameters - default to true (matching Laravel)
	shouldPaginate := true
	paginateParam := c.DefaultQuery("paginate", "false")
	if paginateParam == "false" {
		shouldPaginate = false
	}
	appends["paginate"] = paginateParam

	if perPage := c.Query("per_page"); perPage != "" {
		// Convert perPage to int
		perPageInt := 15 // Default value
		if val, err := strconv.Atoi(perPage); err == nil {
			perPageInt = val
		}
		appends["per_page"] = perPageInt
	}

	if page := c.Query("page"); page != "" {
		// Convert page to int
		pageInt := 1 // Default value
		if val, err := strconv.Atoi(page); err == nil {
			pageInt = val
		}
		appends["page"] = pageInt
	}

	// Get media from service
	media, err := h.mediaService.GetAllMedia(c, filters, appends)
	if err != nil {
		h.respHelper.SendError(c, "Failed to retrieve media", err.Error(), http.StatusInternalServerError)
		return
	}

	// Send appropriate response based on pagination
	if shouldPaginate {
		// For paginated response, the media variable should be a map with pagination data
		h.respHelper.PaginatedResponse(c, media, "Medias fetched successfully")
	} else {
		// For non-paginated response, the media variable should be a slice of Media models
		h.respHelper.OkResponse(c, media, "Medias fetched successfully")
	}
}

// CreateMedia handles the request to create a new media
func (h *MediaHandler) CreateMedia(c *gin.Context) {
	// Parse and validate request
	var request requests.MediaCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.respHelper.SendError(c, "Invalid request format", err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the request
	if err := h.validator.Struct(request); err != nil {
		validationErrors := h.validator.GenerateValidationErrors(err)
		h.respHelper.ValidationError(c, validationErrors, "Validation failed")
		return
	}

	// Create media directly using the request data
	media, err := h.mediaService.CreateMedia(c, request)
	if err != nil {
		h.respHelper.SendError(c, "Failed to create media", err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created media
	h.respHelper.CreatedResponse(c, media, "Media created successfully")
}

// DeleteMedia handles the request to delete a media
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	err := h.mediaService.DeleteMedia(c, id)
	if err != nil {
		h.respHelper.SendError(c, "Failed to delete media", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, nil, "Media deleted successfully")
}
