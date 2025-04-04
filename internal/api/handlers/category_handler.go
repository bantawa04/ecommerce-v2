package handlers

import (
	"net/http"
	"strconv"

	"beautyessentials.com/internal/api/responses"
	"beautyessentials.com/internal/service/interfaces"
	"beautyessentials.com/internal/requests"
	"beautyessentials.com/internal/validators"
	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category-related requests
type CategoryHandler struct {
	categoryService interfaces.CategoryService
	respHelper      *responses.ResponseHelper
	validator       *validators.Validator
}

// NewCategoryHandler creates a new instance of CategoryHandler
func NewCategoryHandler(
	categoryService interfaces.CategoryService,
	respHelper *responses.ResponseHelper,
) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		respHelper:      respHelper,
		validator:       validators.NewValidator(),
	}
}

// GetAllCategories handles the request to get all categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
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

	// Get categories from service
	categories, err := h.categoryService.GetAllCategories(c, filters, appends)
	if err != nil {
		h.respHelper.SendError(c, "Failed to retrieve categories", err.Error(), http.StatusInternalServerError)
		return
	}

	// Send appropriate response based on pagination
	if shouldPaginate {
		// For paginated response, the categories variable should be a map with pagination data
		h.respHelper.PaginatedResponse(c, categories, "Categories retrieved successfully")
	} else {
		// For non-paginated response, the categories variable should be a slice of Category models
		h.respHelper.OkResponse(c, categories, "Categories retrieved successfully")
	}
}

// GetCategory handles the request to get a specific category
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryService.FindCategory(c, id)
	if err != nil {
		h.respHelper.SendError(c, "Category not found", err.Error(), http.StatusNotFound)
		return
	}

	h.respHelper.OkResponse(c, category, "Category retrieved successfully")
}

// CreateCategory handles the request to create a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	// Parse and validate request
	var request requests.CategoryCreateRequest
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

	// Create category directly using the request data
	category, err := h.categoryService.CreateCategory(c, request)
	if err != nil {
		h.respHelper.SendError(c, "Failed to create category", err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created category
	h.respHelper.CreatedResponse(c, category, "Category created successfully")
}

// UpdateCategory handles the request to update a category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	// Parse and validate request
	var request requests.CategoryUpdateRequest
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
	if request.Name != "" {
		data["name"] = request.Name
	}
	if request.Description != "" {
		data["description"] = request.Description
	}
	if request.Status != "" {
		data["status"] = request.Status
	}
	if request.MediaID != "" {
		data["media_id"] = request.MediaID
	}

	category, err := h.categoryService.UpdateCategory(c, data, id)
	if err != nil {
		h.respHelper.SendError(c, "Failed to update category", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, category, "Category updated successfully")
}

// DeleteCategory handles the request to delete a category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	err := h.categoryService.DeleteCategory(c, id)
	if err != nil {
		h.respHelper.SendError(c, "Failed to delete category", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, nil, "Category deleted successfully")
}

// GetActiveCategories handles the request to get all active categories
func (h *CategoryHandler) GetActiveCategories(c *gin.Context) {
	categories, err := h.categoryService.GetActiveCategories(c)
	if err != nil {
		h.respHelper.SendError(c, "Failed to retrieve active categories", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, categories, "Active categories retrieved successfully")
}

// FindCategoryBySlug handles the request to find categories by slug
func (h *CategoryHandler) FindCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	categories, err := h.categoryService.FindCategoryBySlug(c, slug)
	if err != nil {
		h.respHelper.SendError(c, "Failed to retrieve categories by slug", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, categories, "Categories retrieved successfully")
}