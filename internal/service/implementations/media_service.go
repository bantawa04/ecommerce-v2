package implementations

import (
	"context"

	"beautyessentials.com/internal/dto"
	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/repository/interfaces"
	"beautyessentials.com/internal/requests"
	"beautyessentials.com/internal/service/external"
	serviceInterfaces "beautyessentials.com/internal/service/interfaces"
)

// MediaService implements the MediaService interface
type MediaService struct {
	mediaRepo       interfaces.MediaRepository
	imageKitService *external.ImageKitService
}

// NewMediaService creates a new instance of MediaService
func NewMediaService(mediaRepo interfaces.MediaRepository, imageKitService *external.ImageKitService) serviceInterfaces.MediaService {
	return &MediaService{
		mediaRepo:       mediaRepo,
		imageKitService: imageKitService,
	}
}

// GetAllMedia retrieves all media with filtering and pagination
func (s *MediaService) GetAllMedia(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error) {
	result, err := s.mediaRepo.GetAllMedia(ctx, filters, appends)
	if err != nil {
		return nil, err
	}

	// Check if the result is paginated
	if paginatedResult, ok := result.(map[string]interface{}); ok {
		return dto.TransformMediaPagination(paginatedResult), nil
	}

	// If not paginated, it should be a slice of Media models
	if media, ok := result.([]models.Media); ok {
		return dto.TransformMediaCollection(media), nil
	}

	// Return as is if we can't transform
	return result, nil
}

// FindMedia finds a media by ID
func (s *MediaService) FindMedia(ctx context.Context, id string) (dto.MediaDTO, error) {
	media, err := s.mediaRepo.FindMedia(ctx, id)
	if err != nil {
		return dto.MediaDTO{}, err
	}
	return dto.FromMediaModel(media), nil
}

// CreateMedia creates a new media
func (s *MediaService) CreateMedia(ctx context.Context, request requests.MediaCreateRequest) (dto.MediaDTO, error) {
	// Convert request to data map
	data := map[string]interface{}{
		"file_id":   request.FileID,
		"url":       request.URL,
		"thumb_url": request.ThumbURL,
	}

	media, err := s.mediaRepo.CreateMedia(ctx, data)
	if err != nil {
		return dto.MediaDTO{}, err
	}

	return dto.FromMediaModel(media), nil
}

// DeleteMedia deletes a media
func (s *MediaService) DeleteMedia(ctx context.Context, id string) error {
	// Find the media first to get the file ID
	media, err := s.mediaRepo.FindMedia(ctx, id)
	if err != nil {
		return err
	}

	// Delete the file from ImageKit if file_id exists
	if media.FileID != "" {
		if err := s.imageKitService.DeleteFile(media.FileID); err != nil {
			return err
		}
	}

	// Delete the media record from the database
	return s.mediaRepo.DeleteMedia(ctx, id)
}
