package requests

// MediaCreateRequest represents the request to create a media
type MediaCreateRequest struct {
	FileID   string `json:"file_id" validate:"required"`
	URL      string `json:"url" validate:"required,url"`
	ThumbURL string `json:"thumb_url" validate:"omitempty,url"`
}

// MediaUpdateRequest represents the request to update a media
type MediaUpdateRequest struct {
	FileID   string `json:"file_id" validate:"omitempty"`
	URL      string `json:"url" validate:"omitempty,url"`
	ThumbURL string `json:"thumb_url" validate:"omitempty,url"`
}

// MediaUploadURLRequest represents the request to upload a media from URL
type MediaUploadURLRequest struct {
	URL      string `json:"url" validate:"required,url"`
	FileName string `json:"file_name" validate:"omitempty"`
}
