package implementations

import (
	"context"

	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/repository/interfaces"
	serviceInterfaces "beautyessentials.com/internal/service/interfaces"
)

// BrandService implements the BrandService interface
type BrandService struct {
	brandRepo interfaces.BrandRepository
}

// NewBrandService creates a new instance of BrandService
func NewBrandService(brandRepo interfaces.BrandRepository) serviceInterfaces.BrandService {
	return &BrandService{
		brandRepo: brandRepo,
	}
}

// GetAllBrands retrieves all brands
func (s *BrandService) GetAllBrands(ctx context.Context) ([]models.Brand, error) {
	return s.brandRepo.GetAllBrands(ctx)
}