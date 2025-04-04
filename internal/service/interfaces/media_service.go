package interfaces

import (
	"context"
	"mime/multipart"

	"beautyessentials.com/internal/dto"
	"beautyessentials.com/internal/requests"
)

// MediaService defines the interface for media service operations
type MediaService interface {
	GetAllMedia(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error)
	FindMedia(ctx context.Context, id string) (dto.MediaDTO, error)
	CreateMedia(ctx context.Context, request requests.MediaCreateRequest) (dto.MediaDTO, error)
	UpdateMedia(ctx context.Context, data map[string]interface{}, id string) (dto.MediaDTO, error)
	DeleteMedia(ctx context.Context, id string) error
	UploadFile(ctx context.Context, file *multipart.FileHeader) (dto.MediaDTO, error)
	UploadFromURL(ctx context.Context, url string, fileName string) (dto.MediaDTO, error)
}