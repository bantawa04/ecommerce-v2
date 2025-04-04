package interfaces

import (
	"context"

	"beautyessentials.com/internal/models"
)

// MediaRepository defines the interface for media repository operations
type MediaRepository interface {
	GetAllMedia(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error)
	FindMedia(ctx context.Context, id string) (models.Media, error)
	CreateMedia(ctx context.Context, data map[string]interface{}) (models.Media, error)
	UpdateMedia(ctx context.Context, data map[string]interface{}, id string) (models.Media, error)
	DeleteMedia(ctx context.Context, id string) error
	FindMediaByFileID(ctx context.Context, fileID string) (models.Media, error)
}