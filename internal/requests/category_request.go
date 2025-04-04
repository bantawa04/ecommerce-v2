package requests

// CategoryCreateRequest represents the request to create a category
type CategoryCreateRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"omitempty"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive"`
	MediaID     string `json:"media_id" validate:"omitempty,ulid"`
}

// CategoryUpdateRequest represents the request to update a category
type CategoryUpdateRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=255"`
	Description string `json:"description" validate:"omitempty"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive"`
	MediaID     string `json:"media_id" validate:"omitempty,ulid"`
}