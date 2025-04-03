package interfaces

import (
	"context"

	"beautyessentials.com/internal/models"
)

// BrandRepository defines the interface for brand data operations
type BrandRepository interface {
	GetAllBrands(ctx context.Context) ([]models.Brand, error)
}