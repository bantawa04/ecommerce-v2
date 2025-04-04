package interfaces

import (
	"context"

	"beautyessentials.com/internal/models"
)

// MediaRepository defines the interface for media-related database operations
type MediaRepository interface {
	GetAllMedia(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error)
	CreateMedia(ctx context.Context, data map[string]interface{}) (models.Media, error)
	DeleteMedia(ctx context.Context, id string) error
	FindMedia(ctx context.Context, id string) (models.Media, error) // Needed for delete operation
}