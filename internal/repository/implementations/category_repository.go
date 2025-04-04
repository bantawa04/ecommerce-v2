package implementations

import (
	"context"

	"beautyessentials.com/internal/constant"
	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/repository/interfaces"
	"gorm.io/gorm"
)

// CategoryRepository implements the CategoryRepository interface
type CategoryRepository struct {
	db       *gorm.DB
	paginate int
}

// NewCategoryRepository creates a new instance of CategoryRepository
func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepository {
	return &CategoryRepository{
		db:       db,
		paginate: 15, // Default pagination size
	}
}

// GetAllCategories retrieves all categories from the database with filtering and pagination
func (r *CategoryRepository) GetAllCategories(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error) {
	var categories []models.Category

	// Start with base query
	query := r.db.WithContext(ctx).Model(&models.Category{})

	// Apply filters
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// Handle trashed (soft deleted) records
	if trashed, ok := filters["trashed"].(bool); ok && trashed {
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	}

	// Apply sorting
	sortBy := "created_at"
	sortDirection := "desc"

	if sb, ok := filters["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}

	if sd, ok := filters["sort_direction"].(string); ok && sd != "" {
		sortDirection = sd
	}

	query = query.Order(sortBy + " " + sortDirection)

	// Handle pagination
	shouldPaginate := true
	if pag, ok := appends["paginate"].(string); ok {
		shouldPaginate = pag == "true"
	}

	perPage := r.paginate
	if pp, ok := appends["per_page"].(int); ok {
		perPage = pp
	}

	// Execute query with or without pagination
	if shouldPaginate {
		var total int64
		query.Count(&total)

		page := 1
		if p, ok := appends["page"].(int); ok {
			page = p
		}

		offset := (page - 1) * perPage

		result := query.Limit(perPage).Offset(offset).Find(&categories)
		if result.Error != nil {
			return nil, result.Error
		}

		// Return paginated result
		return map[string]interface{}{
			"data":         categories,
			"total":        total,
			"per_page":     perPage,
			"current_page": page,
			"last_page":    (int(total) + perPage - 1) / perPage,
		}, nil
	}

	// Return all results without pagination
	result := query.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

// FindCategory finds a category by ID
func (r *CategoryRepository) FindCategory(ctx context.Context, id string) (models.Category, error) {
	var category models.Category
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&category)
	if result.Error != nil {
		return models.Category{}, result.Error
	}
	return category, nil
}

// CreateCategory creates a new category
func (r *CategoryRepository) CreateCategory(ctx context.Context, data map[string]interface{}) (models.Category, error) {
	// Create a new category instance
	category := models.Category{
		Name:        data["name"].(string),
	}

	// If status is provided, set it
	if status, ok := data["status"].(string); ok {
		category.Status = constant.StatusEnum(status)
	} else {
		category.Status = constant.StatusActive // Default status
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return models.Category{}, tx.Error
	}

	// Defer a rollback in case anything fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the category within the transaction
	if err := tx.Create(&category).Error; err != nil {
		tx.Rollback()
		return models.Category{}, err
	}

	// Handle media attachment if provided
	if mediaID, ok := data["media_id"].(string); ok && mediaID != "" {
		if err := tx.Exec("INSERT INTO mediables (id, mediable_id, mediable_type, media_id) VALUES (?, ?, ?, ?)",
			data["mediable_id"], category.ID, "App\\Model\\Category", mediaID).Error; err != nil {
			tx.Rollback()
			return models.Category{}, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return models.Category{}, err
	}

	// Return the created category
	return category, nil
}

// UpdateCategory updates an existing category
func (r *CategoryRepository) UpdateCategory(ctx context.Context, data map[string]interface{}, id string) (models.Category, error) {
	// Find the category first
	category, err := r.FindCategory(ctx, id)
	if err != nil {
		return models.Category{}, err
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return models.Category{}, tx.Error
	}

	// Defer a rollback in case anything fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update the category within the transaction
	if err := tx.Model(&category).Updates(data).Error; err != nil {
		tx.Rollback()
		return models.Category{}, err
	}

	// Handle media sync if provided
	if mediaID, ok := data["media_id"].(string); ok && mediaID != "" {
		// First delete existing media relationships
		if err := tx.Exec("DELETE FROM mediables WHERE mediable_id = ? AND mediable_type = ?",
			category.ID, "App\\Model\\Category").Error; err != nil {
			tx.Rollback()
			return models.Category{}, err
		}

		// Then add the new one
		if err := tx.Exec("INSERT INTO mediables (id, mediable_id, mediable_type, media_id) VALUES (?, ?, ?, ?)",
			data["mediable_id"], category.ID, "App\\Model\\Category", mediaID).Error; err != nil {
			tx.Rollback()
			return models.Category{}, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return models.Category{}, err
	}

	// Refresh the category data
	return r.FindCategory(ctx, id)
}

// DeleteCategory soft deletes a category
func (r *CategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	// Find the category first
	category, err := r.FindCategory(ctx, id)
	if err != nil {
		return err
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Defer a rollback in case anything fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete the category within the transaction
	if err := tx.Delete(&category).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// GetActiveCategories retrieves all active categories
func (r *CategoryRepository) GetActiveCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	result := r.db.WithContext(ctx).Where("status = ?", constant.StatusActive).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

// FindCategoryBySlug finds categories by slug
func (r *CategoryRepository) FindCategoryBySlug(ctx context.Context, slug string) ([]models.Category, error) {
	var categories []models.Category
	result := r.db.WithContext(ctx).Where("slug = ?", slug).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}
