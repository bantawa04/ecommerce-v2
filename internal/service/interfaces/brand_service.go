package interfaces

import (
	"context"

	"beautyessentials.com/internal/models"
)

// BrandService defines the interface for brand business logic
type BrandService interface {
	GetAllBrands(ctx context.Context) ([]models.Brand, error)
}