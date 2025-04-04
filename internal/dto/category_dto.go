package dto

import (
	"time"

	"beautyessentials.com/internal/constant"
	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/utils/transformer"
	"gorm.io/gorm"
)

// CategoryDTO represents the data transfer object for Category
type CategoryDTO struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	Status    string         `json:"status"`
	CreatedAt *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty"`
}

// FromCategoryModel converts a Category model to a CategoryDTO
func FromCategoryModel(category models.Category) CategoryDTO {
	return CategoryDTO{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,

		Status:    string(category.Status),
		CreatedAt: &category.CreatedAt,
		UpdatedAt: &category.UpdatedAt,
		DeletedAt: category.DeletedAt,
	}
}

// ToCategoryModel converts a CategoryDTO to a Category model
func (dto CategoryDTO) ToCategoryModel() models.Category {

	return models.Category{
		ID:        dto.ID,
		Name:      dto.Name,
		Slug:      dto.Slug,
		Status:    constant.StatusEnum(dto.Status),
		CreatedAt: *dto.CreatedAt,
		UpdatedAt: *dto.UpdatedAt,
		DeletedAt: dto.DeletedAt,
	}
}

// TransformCategoryCollection transforms a slice of Category models to a slice of CategoryDTOs
func TransformCategoryCollection(categories []models.Category) []CategoryDTO {
	return transformer.TransformCollection(categories, FromCategoryModel)
}

// TransformCategoryPagination transforms a paginated result of Category models to a paginated result of CategoryDTOs
func TransformCategoryPagination(paginatedResult map[string]interface{}) map[string]interface{} {
	return transformer.TransformPagination(paginatedResult, FromCategoryModel)
}
