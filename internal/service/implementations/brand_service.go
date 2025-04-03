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

// GetAllBrands retrieves all brands with filtering and pagination
func (s *BrandService) GetAllBrands(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error) {
	return s.brandRepo.GetAllBrands(ctx, filters, appends)
}

// FindBrand finds a brand by ID
func (s *BrandService) FindBrand(ctx context.Context, id string) (models.Brand, error) {
	return s.brandRepo.FindBrand(ctx, id)
}

// CreateBrand creates a new brand
func (s *BrandService) CreateBrand(ctx context.Context, data map[string]interface{}) (models.Brand, error) {
	return s.brandRepo.CreateBrand(ctx, data)
}

// UpdateBrand updates an existing brand
func (s *BrandService) UpdateBrand(ctx context.Context, data map[string]interface{}, id string) (models.Brand, error) {
	return s.brandRepo.UpdateBrand(ctx, data, id)
}

// DeleteBrand deletes a brand
func (s *BrandService) DeleteBrand(ctx context.Context, id string) error {
	return s.brandRepo.DeleteBrand(ctx, id)
}

// GetActiveBrands retrieves all active brands
func (s *BrandService) GetActiveBrands(ctx context.Context) ([]models.Brand, error) {
	return s.brandRepo.GetActiveBrands(ctx)
}

// GetGroupedBrands retrieves brands grouped by first letter
func (s *BrandService) GetGroupedBrands(ctx context.Context) (map[string][]models.Brand, error) {
	return s.brandRepo.GetGroupedBrands(ctx)
}