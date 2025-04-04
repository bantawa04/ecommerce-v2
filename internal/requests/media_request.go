package requests

// MediaCreateRequest represents the request to create a media
type MediaCreateRequest struct {
	FileID   string `json:"file_id" validate:"required"`
	URL      string `json:"url" validate:"required,url"`
	ThumbURL string `json:"thumb_url" validate:"omitempty,url"`
}
