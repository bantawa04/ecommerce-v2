package dto

import (
	"time"

	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/utils/transformer"
	"gorm.io/gorm"
)

// MediaDTO represents the data transfer object for Media
type MediaDTO struct {
	ID          string         `json:"id"`
	FileID      string         `json:"file_id"`
	FileName    string         `json:"file_name"`
	URL         string         `json:"url"`
	ThumbURL    string         `json:"thumb_url"`
	FileType    string         `json:"file_type"`
	Size        int64          `json:"size"`
	Description string         `json:"description"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// FromMediaModel converts a Media model to a MediaDTO
func FromMediaModel(media models.Media) MediaDTO {
	return MediaDTO{
		ID:          media.ID,
		FileID:      media.FileID,
		FileName:    media.FileName,
		URL:         media.URL,
		ThumbURL:    media.ThumbURL,
		FileType:    media.FileType,
		Size:        media.Size,
		Description: media.Description,
		CreatedAt:   &media.CreatedAt,
		UpdatedAt:   &media.UpdatedAt,
		DeletedAt:   media.DeletedAt,
	}
}

// ToMediaModel converts a MediaDTO to a Media model
func (dto MediaDTO) ToMediaModel() models.Media {

	return models.Media{
		ID:          dto.ID,
		FileID:      dto.FileID,
		FileName:    dto.FileName,
		URL:         dto.URL,
		ThumbURL:    dto.ThumbURL,
		FileType:    dto.FileType,
		Size:        dto.Size,
		Description: dto.Description,
		CreatedAt:   *dto.CreatedAt,
		UpdatedAt:   *dto.UpdatedAt,
		DeletedAt:   dto.DeletedAt,
	}
}

// TransformMediaCollection transforms a slice of Media models to a slice of MediaDTOs
func TransformMediaCollection(media []models.Media) []MediaDTO {
	return transformer.TransformCollection(media, FromMediaModel)
}

// TransformMediaPagination transforms a paginated result of Media models to a paginated result of MediaDTOs
func TransformMediaPagination(paginatedResult map[string]interface{}) map[string]interface{} {
	return transformer.TransformPagination(paginatedResult, FromMediaModel)
}
