package implementations

import (
	"context"

	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/repository/interfaces"
	"gorm.io/gorm"
)

// BrandRepository implements the BrandRepository interface
type BrandRepository struct {
	db *gorm.DB
}

// NewBrandRepository creates a new instance of BrandRepository
func NewBrandRepository(db *gorm.DB) interfaces.BrandRepository {
	return &BrandRepository{
		db: db,
	}
}

// GetAllBrands retrieves all brands from the database
func (r *BrandRepository) GetAllBrands(ctx context.Context) ([]models.Brand, error) {
	var brands []models.Brand
	result := r.db.WithContext(ctx).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	return brands, nil
}