package implementations

import (
	"context"

	"beautyessentials.com/internal/dto"
	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/repository/interfaces"
	serviceInterfaces "beautyessentials.com/internal/service/interfaces"
	"beautyessentials.com/internal/requests"
	"github.com/oklog/ulid/v2"
)

// CategoryService implements the CategoryService interface
type CategoryService struct {
	categoryRepo interfaces.CategoryRepository
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(categoryRepo interfaces.CategoryRepository) serviceInterfaces.CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// GetAllCategories retrieves all categories with filtering and pagination
func (s *CategoryService) GetAllCategories(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error) {
	result, err := s.categoryRepo.GetAllCategories(ctx, filters, appends)
	if err != nil {
		return nil, err
	}

	// Check if the result is paginated
	if paginatedResult, ok := result.(map[string]interface{}); ok {
		return dto.TransformCategoryPagination(paginatedResult), nil
	}

	// If not paginated, it should be a slice of Category models
	if categories, ok := result.([]models.Category); ok {
		return dto.TransformCategoryCollection(categories), nil
	}

	// Return as is if we can't transform
	return result, nil
}

// FindCategory finds a category by ID
func (s *CategoryService) FindCategory(ctx context.Context, id string) (dto.CategoryDTO, error) {
	category, err := s.categoryRepo.FindCategory(ctx, id)
	if err != nil {
		return dto.CategoryDTO{}, err
	}
	return dto.FromCategoryModel(category), nil
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(ctx context.Context, request requests.CategoryCreateRequest) (dto.CategoryDTO, error) {
	// Convert request to data map
	data := map[string]interface{}{
		"name":        request.Name,
		"description": request.Description,
		"status":      request.Status,
	}
	
	// Add media ID if provided
	if request.MediaID != "" {
		data["media_id"] = request.MediaID
		data["mediable_id"] = ulid.Make().String() // Generate a new ULID for the mediable relation
	}
	
	category, err := s.categoryRepo.CreateCategory(ctx, data)
	if err != nil {
		return dto.CategoryDTO{}, err
	}
	
	return dto.FromCategoryModel(category), nil
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(ctx context.Context, data map[string]interface{}, id string) (dto.CategoryDTO, error) {
	// Add mediable_id if media_id is provided
	if _, ok := data["media_id"]; ok {
		data["mediable_id"] = ulid.Make().String() // Generate a new ULID for the mediable relation
	}
	
	category, err := s.categoryRepo.UpdateCategory(ctx, data, id)
	if err != nil {
		return dto.CategoryDTO{}, err
	}
	
	return dto.FromCategoryModel(category), nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.categoryRepo.DeleteCategory(ctx, id)
}

// GetActiveCategories retrieves all active categories
func (s *CategoryService) GetActiveCategories(ctx context.Context) ([]dto.CategoryDTO, error) {
	categories, err := s.categoryRepo.GetActiveCategories(ctx)
	if err != nil {
		return nil, err
	}
	return dto.TransformCategoryCollection(categories), nil
}

// FindCategoryBySlug finds categories by slug
func (s *CategoryService) FindCategoryBySlug(ctx context.Context, slug string) ([]dto.CategoryDTO, error) {
	categories, err := s.categoryRepo.FindCategoryBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return dto.TransformCategoryCollection(categories), nil
}