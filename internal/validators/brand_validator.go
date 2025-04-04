package validators

// BrandCreateRequest represents the request structure for brand creation
type BrandCreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// BrandUpdateRequest represents the request structure for brand updates
type BrandUpdateRequest struct {
	Name string `json:"name" validate:"omitempty,min=2,max=100"`
}
