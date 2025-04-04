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

// BrandHandler handles brand-related requests
type BrandHandler struct {
	brandService interfaces.BrandService
	respHelper   *responses.ResponseHelper
	validator    *validators.Validator
}

// NewBrandHandler creates a new instance of BrandHandler
func NewBrandHandler(
	brandService interfaces.BrandService,
	respHelper *responses.ResponseHelper,
) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
		respHelper:   respHelper,
		validator:    validators.NewValidator(),
	}
}

// GetAllBrands handles the request to get all brands
func (h *BrandHandler) GetAllBrands(c *gin.Context) {
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

	// Get brands from service
	brands, err := h.brandService.GetAllBrands(c, filters, appends)
	if err != nil {
		h.respHelper.SendError(c, "Failed to retrieve brands", err.Error(), http.StatusInternalServerError)
		return
	}

	// Send appropriate response based on pagination
	if shouldPaginate {
		// For paginated response, the brands variable should be a map with pagination data
		h.respHelper.PaginatedResponse(c, brands, "Brands retrieved successfully")
	} else {
		// For non-paginated response, the brands variable should be a slice of Brand models
		h.respHelper.OkResponse(c, brands, "Brands retrieved successfully")
	}
}

// GetBrand handles the request to get a specific brand
func (h *BrandHandler) GetBrand(c *gin.Context) {
	id := c.Param("id")
	brand, err := h.brandService.FindBrand(c, id)
	if err != nil {
		h.respHelper.SendError(c, "Brand not found", err.Error(), http.StatusNotFound)
		return
	}

	h.respHelper.OkResponse(c, brand, "Brand retrieved successfully")
}

// CreateBrand handles the request to create a new brand
func (h *BrandHandler) CreateBrand(c *gin.Context) {
	// Parse and validate request
	var request requests.BrandCreateRequest
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

	// Create brand directly using the request data
	brand, err := h.brandService.CreateBrand(c, request)
	if err != nil {
		h.respHelper.SendError(c, "Failed to create brand", err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created brand
	h.respHelper.CreatedResponse(c, brand, "Brand created successfully")
}

// UpdateBrand handles the request to update a brand
func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	id := c.Param("id")

	// Parse and validate request
	var request requests.BrandUpdateRequest
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

	brand, err := h.brandService.UpdateBrand(c, data, id)
	if err != nil {
		h.respHelper.SendError(c, "Failed to update brand", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, brand, "Brand updated successfully")
}

// DeleteBrand handles the request to delete a brand
func (h *BrandHandler) DeleteBrand(c *gin.Context) {
	id := c.Param("id")
	err := h.brandService.DeleteBrand(c, id)
	if err != nil {
		h.respHelper.SendError(c, "Failed to delete brand", err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	h.respHelper.SendSuccess(c, "Brand deleted successfully", http.StatusOK)
}

// GetActiveBrands handles the request to get all active brands
func (h *BrandHandler) GetActiveBrands(c *gin.Context) {
	brands, err := h.brandService.GetActiveBrands(c)
	if err != nil {
		h.respHelper.SendError(c, "Failed to retrieve active brands", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, brands, "Active brands retrieved successfully")
}

// GetGroupedBrands handles the request to get brands grouped by first letter
func (h *BrandHandler) GetGroupedBrands(c *gin.Context) {
	brands, err := h.brandService.GetGroupedBrands(c)
	if err != nil {
		h.respHelper.SendError(c, "Failed to retrieve grouped brands", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, brands, "Grouped brands retrieved successfully")
}
