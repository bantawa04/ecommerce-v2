package handlers

import (
	"net/http"

	"beautyessentials.com/internal/api/responses"
	"beautyessentials.com/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

// BrandHandler handles brand-related requests
type BrandHandler struct {
	brandService interfaces.BrandService
	respHelper   *responses.ResponseHelper
}

// NewBrandHandler creates a new instance of BrandHandler
func NewBrandHandler(
	brandService interfaces.BrandService,
	respHelper *responses.ResponseHelper,
) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
		respHelper:   respHelper,
	}
}

// GetAllBrands handles the request to get all brands
func (h *BrandHandler) GetAllBrands(c *gin.Context) {
	brands, err := h.brandService.GetAllBrands(c)
	if err != nil {
		h.respHelper.SendError(c, "Failed to retrieve brands", err.Error(), http.StatusInternalServerError)
		return
	}

	h.respHelper.OkResponse(c, brands, "Brands retrieved successfully")
}