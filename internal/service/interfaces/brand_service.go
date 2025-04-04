package interfaces

import (
	"context"

	"beautyessentials.com/internal/dto"
	"beautyessentials.com/internal/validators"
)

// BrandService defines the interface for brand business logic
type BrandService interface {
	GetAllBrands(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error)
	FindBrand(ctx context.Context, id string) (dto.BrandDTO, error)
	CreateBrand(ctx context.Context, request validators.BrandCreateRequest) (dto.BrandDTO, error)
	UpdateBrand(ctx context.Context, data map[string]interface{}, id string) (dto.BrandDTO, error)
	DeleteBrand(ctx context.Context, id string) error
	GetActiveBrands(ctx context.Context) ([]dto.BrandDTO, error)
	GetGroupedBrands(ctx context.Context) (map[string][]dto.BrandDTO, error)
}