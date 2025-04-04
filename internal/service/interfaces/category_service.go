package interfaces

import (
	"context"

	"beautyessentials.com/internal/dto"
	"beautyessentials.com/internal/requests"
)

// CategoryService defines the interface for category service operations
type CategoryService interface {
	GetAllCategories(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error)
	FindCategory(ctx context.Context, id string) (dto.CategoryDTO, error)
	CreateCategory(ctx context.Context, request requests.CategoryCreateRequest) (dto.CategoryDTO, error)
	UpdateCategory(ctx context.Context, data map[string]interface{}, id string) (dto.CategoryDTO, error)
	DeleteCategory(ctx context.Context, id string) error
	GetActiveCategories(ctx context.Context) ([]dto.CategoryDTO, error)
	FindCategoryBySlug(ctx context.Context, slug string) ([]dto.CategoryDTO, error)
}
