package implementations

import (
	"context"
	"strings"
	"time"

	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/repository/interfaces"
	"gorm.io/gorm"
)

// BrandRepository implements the BrandRepository interface
type BrandRepository struct {
	db       *gorm.DB
	paginate int
}

// NewBrandRepository creates a new instance of BrandRepository
func NewBrandRepository(db *gorm.DB) interfaces.BrandRepository {
	return &BrandRepository{
		db:       db,
		paginate: 15, // Default pagination size
	}
}

// GetAllBrands retrieves all brands from the database with filtering and pagination
func (r *BrandRepository) GetAllBrands(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error) {
	var brands []models.Brand

	// Start with base query
	query := r.db.WithContext(ctx).Model(&models.Brand{})

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

		result := query.Limit(perPage).Offset(offset).Find(&brands)
		if result.Error != nil {
			return nil, result.Error
		}

		// Return paginated result
		return map[string]interface{}{
			"data":         brands,
			"total":        total,
			"per_page":     perPage,
			"current_page": page,
			"last_page":    (int(total) + perPage - 1) / perPage,
		}, nil
	}

	// Return all results without pagination
	result := query.Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	return brands, nil
}

// FindBrand finds a brand by ID
func (r *BrandRepository) FindBrand(ctx context.Context, id string) (models.Brand, error) {
	var brand models.Brand
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&brand)
	if result.Error != nil {
		return models.Brand{}, result.Error
	}
	return brand, nil
}

// CreateBrand creates a new brand
func (r *BrandRepository) CreateBrand(ctx context.Context, data map[string]interface{}) (models.Brand, error) {
	// Create a new brand instance
	brand := models.Brand{
		Name: data["name"].(string),
	}

	// If status is provided, set it
	if status, ok := data["status"].(string); ok {
		brand.Status = models.StatusEnum(status)
	} else {
		brand.Status = models.StatusActive // Default status
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return models.Brand{}, tx.Error
	}

	// Defer a rollback in case anything fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the brand within the transaction
	if err := tx.Create(&brand).Error; err != nil {
		tx.Rollback()
		return models.Brand{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return models.Brand{}, err
	}

	// Return the created brand
	return brand, nil
}

// UpdateBrand updates an existing brand
func (r *BrandRepository) UpdateBrand(ctx context.Context, data map[string]interface{}, id string) (models.Brand, error) {
	// Find the brand first
	brand, err := r.FindBrand(ctx, id)
	if err != nil {
		return models.Brand{}, err
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return models.Brand{}, tx.Error
	}

	// Defer a rollback in case anything fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update the brand within the transaction
	if err := tx.Model(&brand).Updates(data).Error; err != nil {
		tx.Rollback()
		return models.Brand{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return models.Brand{}, err
	}

	// Refresh the brand data
	return r.FindBrand(ctx, id)
}

// DeleteBrand soft deletes a brand
func (r *BrandRepository) DeleteBrand(ctx context.Context, id string) error {
	// Find the brand first
	brand, err := r.FindBrand(ctx, id)
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

	// Use UpdateColumn instead of Update to bypass soft delete callbacks
	now := time.Now()
	if err := tx.Model(&brand).UpdateColumn("deleted_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// GetActiveBrands retrieves all active brands
func (r *BrandRepository) GetActiveBrands(ctx context.Context) ([]models.Brand, error) {
	var brands []models.Brand
	result := r.db.WithContext(ctx).Where("status = ?", models.StatusActive).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	return brands, nil
}

// GetGroupedBrands retrieves active brands grouped by first letter
func (r *BrandRepository) GetGroupedBrands(ctx context.Context) (map[string][]models.Brand, error) {
	var brands []models.Brand

	// Select only needed fields
	result := r.db.WithContext(ctx).
		Where("status = ?", models.StatusActive).
		Select("id, name, slug").
		Find(&brands)

	if result.Error != nil {
		return nil, result.Error
	}

	// Group brands by first letter
	groupedBrands := make(map[string][]models.Brand)
	for _, brand := range brands {
		if len(brand.Name) > 0 {
			firstLetter := strings.ToUpper(string(brand.Name[0]))
			groupedBrands[firstLetter] = append(groupedBrands[firstLetter], brand)
		}
	}

	return groupedBrands, nil
}
