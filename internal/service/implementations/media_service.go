package implementations

import (
	"context"
	"mime/multipart"
	"path/filepath"

	"beautyessentials.com/internal/dto"
	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/repository/interfaces"
	"beautyessentials.com/internal/requests"
	serviceInterfaces "beautyessentials.com/internal/service/interfaces"
	"beautyessentials.com/internal/service/external"
)

// MediaService implements the MediaService interface
type MediaService struct {
	mediaRepo     interfaces.MediaRepository
	imageKitService *external.ImageKitService
}

// NewMediaService creates a new instance of MediaService
func NewMediaService(mediaRepo interfaces.MediaRepository, imageKitService *external.ImageKitService) serviceInterfaces.MediaService {
	return &MediaService{
		mediaRepo:     mediaRepo,
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
		"file_id":     request.FileID,
		"file_name":   request.FileName,
		"url":         request.URL,
		"thumb_url":   request.ThumbURL,
		"file_type":   request.FileType,
		"size":        request.Size,
		"description": request.Description,
	}

	media, err := s.mediaRepo.CreateMedia(ctx, data)
	if err != nil {
		return dto.MediaDTO{}, err
	}

	return dto.FromMediaModel(media), nil
}

// UpdateMedia updates an existing media
func (s *MediaService) UpdateMedia(ctx context.Context, data map[string]interface{}, id string) (dto.MediaDTO, error) {
	media, err := s.mediaRepo.UpdateMedia(ctx, data, id)
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

	// Delete the file from ImageKit
	if err := s.imageKitService.DeleteFile(media.FileID); err != nil {
		return err
	}

	// Delete the media record from the database
	return s.mediaRepo.DeleteMedia(ctx, id)
}

// UploadFile uploads a file to ImageKit and creates a media record
func (s *MediaService) UploadFile(ctx context.Context, file *multipart.FileHeader) (dto.MediaDTO, error) {
	// Upload file to ImageKit
	result, err := s.imageKitService.UploadFile(file)
	if err != nil {
		return dto.MediaDTO{}, err
	}

	// Create media record
	data := map[string]interface{}{
		"file_id":     result["fileId"],
		"file_name":   file.Filename,
		"url":         result["url"],
		"thumb_url":   result["thumbUrl"],
		"file_type":   filepath.Ext(file.Filename),
		"size":        result["size"],
		"description": "",
	}

	media, err := s.mediaRepo.CreateMedia(ctx, data)
	if err != nil {
		return dto.MediaDTO{}, err
	}

	return dto.FromMediaModel(media), nil
}

// UploadFromURL uploads a file from a URL to ImageKit and creates a media record
func (s *MediaService) UploadFromURL(ctx context.Context, url string, fileName string) (dto.MediaDTO, error) {
	// Upload file from URL to ImageKit
	result, err := s.imageKitService.UploadFromURL(url, fileName)
	if err != nil {
		return dto.MediaDTO{}, err
	}

	// If fileName is empty, use the one from the result
	if fileName == "" {
		fileName = result["fileName"].(string)
	}

	// Create media record
	data := map[string]interface{}{
		"file_id":     result["fileId"],
		"file_name":   fileName,
		"url":         result["url"],
		"thumb_url":   result["thumbUrl"],
		"file_type":   filepath.Ext(fileName),
		"size":        result["size"],
		"description": "",
	}

	media, err := s.mediaRepo.CreateMedia(ctx, data)
	if err != nil {
		return dto.MediaDTO{}, err
	}

	return dto.FromMediaModel(media), nil
}