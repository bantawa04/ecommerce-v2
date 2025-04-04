package interfaces

import (
	"context"

	"beautyessentials.com/internal/dto"
	"beautyessentials.com/internal/requests"
)

// MediaService defines the interface for media-related operations
type MediaService interface {
	GetAllMedia(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error)
	CreateMedia(ctx context.Context, request requests.MediaCreateRequest) (dto.MediaDTO, error)
	DeleteMedia(ctx context.Context, id string) error
}