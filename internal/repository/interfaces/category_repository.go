package interfaces

import (
	"context"

	"beautyessentials.com/internal/models"
)

// CategoryRepository defines the interface for category repository operations
type CategoryRepository interface {
	GetAllCategories(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error)
	FindCategory(ctx context.Context, id string) (models.Category, error)
	CreateCategory(ctx context.Context, data map[string]interface{}) (models.Category, error)
	UpdateCategory(ctx context.Context, data map[string]interface{}, id string) (models.Category, error)
	DeleteCategory(ctx context.Context, id string) error
	GetActiveCategories(ctx context.Context) ([]models.Category, error)
	FindCategoryBySlug(ctx context.Context, slug string) ([]models.Category, error)
}