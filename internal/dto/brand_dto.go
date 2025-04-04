package dto

import (
	"time"

	"beautyessentials.com/internal/constant"
	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/utils/transformer"
	"gorm.io/gorm"
)

// BrandDTO represents the data transfer object for Brand
type BrandDTO struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	Status    string         `json:"status"`
	CreatedAt *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty"`
}

// FromModel converts a Brand model to a BrandDTO
func FromModel(brand models.Brand) BrandDTO {
	return BrandDTO{
		ID:        brand.ID,
		Name:      brand.Name,
		Slug:      brand.Slug,
		Status:    string(brand.Status),
		CreatedAt: &brand.CreatedAt,
		UpdatedAt: &brand.UpdatedAt,
		DeletedAt: brand.DeletedAt,
	}
}

// ToModel converts a BrandDTO to a Brand model
func (dto BrandDTO) ToModel() models.Brand {
	return models.Brand{
		ID:     dto.ID,
		Name:   dto.Name,
		Slug:   dto.Slug,
		Status: constant.StatusEnum(dto.Status),
		CreatedAt: func() time.Time {
			if dto.CreatedAt != nil {
				return *dto.CreatedAt
			}
			return time.Now()
		}(),
		UpdatedAt: func() time.Time {
			if dto.UpdatedAt != nil {
				return *dto.UpdatedAt
			}
			return time.Now()
		}(),
		DeletedAt: dto.DeletedAt,
	}
}

// TransformBrandCollection transforms a slice of Brand models to a slice of BrandDTOs
func TransformBrandCollection(brands []models.Brand) []BrandDTO {
	return transformer.TransformCollection(brands, FromModel)
}

// TransformBrandPagination transforms a paginated result of Brand models to a paginated result of BrandDTOs
func TransformBrandPagination(paginatedResult map[string]interface{}) map[string]interface{} {
	return transformer.TransformPagination(paginatedResult, FromModel)
}
