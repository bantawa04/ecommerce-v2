package dto

import (
	"time"

	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/utils/transformer"
)

// MediaDTO represents the data transfer object for Media
type MediaDTO struct {
	ID        string     `json:"id"`
	FileID    string     `json:"file_id"`
	URL       string     `json:"url"`
	ThumbURL  string     `json:"thumb_url"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// FromMediaModel converts a Media model to a MediaDTO
func FromMediaModel(media models.Media) MediaDTO {
	return MediaDTO{
		ID:        media.ID,
		FileID:    media.FileID,
		URL:       media.URL,
		ThumbURL:  media.ThumbURL,
		CreatedAt: &media.CreatedAt,
		UpdatedAt: &media.UpdatedAt,
	}
}

// ToMediaModel converts a MediaDTO to a Media model
func (dto MediaDTO) ToMediaModel() models.Media {

	return models.Media{
		ID:        dto.ID,
		FileID:    dto.FileID,
		URL:       dto.URL,
		ThumbURL:  dto.ThumbURL,
		CreatedAt: *dto.CreatedAt,
		UpdatedAt: *dto.UpdatedAt,
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
