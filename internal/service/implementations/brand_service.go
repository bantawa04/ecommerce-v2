package implementations

import (
	"context"

	"beautyessentials.com/internal/dto"
	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/repository/interfaces"
	serviceInterfaces "beautyessentials.com/internal/service/interfaces"
	"beautyessentials.com/internal/validators"
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
	result, err := s.brandRepo.GetAllBrands(ctx, filters, appends)
	if err != nil {
		return nil, err
	}

	// Check if the result is paginated
	if paginatedResult, ok := result.(map[string]interface{}); ok {
		return dto.TransformBrandPagination(paginatedResult), nil
	}

	// If not paginated, it should be a slice of Brand models
	if brands, ok := result.([]models.Brand); ok {
		return dto.TransformBrandCollection(brands), nil
	}

	// Return as is if we can't transform
	return result, nil
}

// GetActiveBrands retrieves all active brands
func (s *BrandService) GetActiveBrands(ctx context.Context) ([]dto.BrandDTO, error) {
	brands, err := s.brandRepo.GetActiveBrands(ctx)
	if err != nil {
		return nil, err
	}
	return dto.TransformBrandCollection(brands), nil
}

// GetGroupedBrands retrieves brands grouped by first letter
func (s *BrandService) GetGroupedBrands(ctx context.Context) (map[string][]dto.BrandDTO, error) {
	groupedBrands, err := s.brandRepo.GetGroupedBrands(ctx)
	if err != nil {
		return nil, err
	}

	// Transform each group
	result := make(map[string][]dto.BrandDTO)
	for letter, brands := range groupedBrands {
		result[letter] = dto.TransformBrandCollection(brands)
	}

	return result, nil
}

// FindBrand finds a brand by ID
func (s *BrandService) FindBrand(ctx context.Context, id string) (dto.BrandDTO, error) {
	brand, err := s.brandRepo.FindBrand(ctx, id)
	if err != nil {
		return dto.BrandDTO{}, err
	}
	return dto.FromModel(brand), nil
}

// CreateBrand creates a new brand
func (s *BrandService) CreateBrand(ctx context.Context, request validators.BrandCreateRequest) (dto.BrandDTO, error) {
	// Convert request to map for repository
	data := map[string]interface{}{
		"name": request.Name,
	}
	
	// Use transaction in repository
	createdBrand, err := s.brandRepo.CreateBrand(ctx, data)
	if err != nil {
		return dto.BrandDTO{}, err
	}
	
	// Convert to DTO and return
	return dto.FromModel(createdBrand), nil
}

// UpdateBrand updates an existing brand
func (s *BrandService) UpdateBrand(ctx context.Context, data map[string]interface{}, id string) (dto.BrandDTO, error) {
	brand, err := s.brandRepo.UpdateBrand(ctx, data, id)
	if err != nil {
		return dto.BrandDTO{}, err
	}
	return dto.FromModel(brand), nil
}

// DeleteBrand deletes a brand
func (s *BrandService) DeleteBrand(ctx context.Context, id string) error {
	return s.brandRepo.DeleteBrand(ctx, id)
}
