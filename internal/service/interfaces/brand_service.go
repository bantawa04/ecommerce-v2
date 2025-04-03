package interfaces

import (
	"context"

	"beautyessentials.com/internal/models"
)

// BrandService defines the interface for brand business logic
type BrandService interface {
	GetAllBrands(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error)
	FindBrand(ctx context.Context, id string) (models.Brand, error)
	CreateBrand(ctx context.Context, data map[string]interface{}) (models.Brand, error)
	UpdateBrand(ctx context.Context, data map[string]interface{}, id string) (models.Brand, error)
	DeleteBrand(ctx context.Context, id string) error
	GetActiveBrands(ctx context.Context) ([]models.Brand, error)
	GetGroupedBrands(ctx context.Context) (map[string][]models.Brand, error)
}