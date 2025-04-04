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

	// Add search filter if provided
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	// Add trashed filter if provided
	if trashed := c.Query("trashed"); trashed != "" {
		filters["trashed"] = trashed == "true"
	}

	// Add sorting parameters
	filters["sort_by"] = c.DefaultQuery("sort_by", "created_at")
	filters["sort_direction"] = c.DefaultQuery("sort_direction", "desc")

	// Add pagination parameters - default to false
	shouldPaginate := false
	paginateParam := c.DefaultQuery("paginate", "false")
	if paginateParam == "true" {
		shouldPaginate = true
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
		h.respHelper.PaginatedResponse(c, media, "Media retrieved successfully")
	} else {
		// For non-paginated response, the media variable should be a slice of Media models
		h.respHelper.OkResponse(c, media, "Media retrieved successfully")
	}
}

// GetMedia handles the request to get a specific media
func (h *MediaHandler) GetMedia(c *gin.Context) {
	id := c.Param("id")
	media, err := h.mediaService.FindMedia(c, id)
	if err != nil {
		h.respHelper.SendError(c, "Media not found", err.Error(), http.StatusNotFound)
		return
	}

	h.respHelper.OkResponse(c, media, "Media retrieved successfully")
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

// UpdateMedia handles the request to update a media
func (h *MediaHandler) UpdateMedia(c *gin.Context) {
	id := c.Param("id")

	// Parse and validate request
	var request requests.MediaUpdateRequest
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

	// Convert validated request to map for service
	data := make(map[string]interface{})
	if request.FileID != "" {
		data["file_id"] = request.FileID
	}
	if request.FileName != "" {
		data["file_name"] = request.FileName
	}
	if request.URL != "" {
		data["url"] = request.URL
	}
	if request.ThumbURL != "" {
		data["thumb_url"] = request.ThumbURL
	}
	if request.FileType != "" {
		data["file_type"] = request.FileType
	}
	if request.Size > 0 {
		data["size"] = request.Size
	}
	if request.Description != "" {
		data["description"] = request.Description
	}

	media, err := h.mediaService.UpdateMedia(c, data, id)
	if err != nil {
		h.respHelper.SendError(c, "Failed to update media", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, media, "Media updated successfully")
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

// UploadFile handles the request to upload a file
func (h *MediaHandler) UploadFile(c *gin.Context) {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		h.respHelper.SendError(c, "Failed to get file from request", err.Error(), http.StatusBadRequest)
		return
	}

	// Upload file
	media, err := h.mediaService.UploadFile(c, file)
	if err != nil {
		h.respHelper.SendError(c, "Failed to upload file", err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created media
	h.respHelper.CreatedResponse(c, media, "File uploaded successfully")
}

// UploadFromURL handles the request to upload a file from a URL
func (h *MediaHandler) UploadFromURL(c *gin.Context) {
	// Parse and validate request
	var request requests.MediaUploadURLRequest
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

	// Upload file from URL
	media, err := h.mediaService.UploadFromURL(c, request.URL, request.FileName)
	if err != nil {
		h.respHelper.SendError(c, "Failed to upload file from URL", err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created media
	h.respHelper.CreatedResponse(c, media, "File uploaded successfully")
}