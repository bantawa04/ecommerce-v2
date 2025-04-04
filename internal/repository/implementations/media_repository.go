package implementations

import (
	"context"

	"beautyessentials.com/internal/models"
	"beautyessentials.com/internal/repository/interfaces"
	"gorm.io/gorm"
)

// MediaRepository implements the MediaRepository interface
type MediaRepository struct {
	db       *gorm.DB
	paginate int
}

// NewMediaRepository creates a new instance of MediaRepository
func NewMediaRepository(db *gorm.DB) interfaces.MediaRepository {
	return &MediaRepository{
		db:       db,
		paginate: 15, // Default pagination size
	}
}

// GetAllMedia retrieves all media from the database with filtering and pagination
func (r *MediaRepository) GetAllMedia(ctx context.Context, filters map[string]interface{}, appends map[string]interface{}) (interface{}, error) {
	var media []models.Media

	// Start with base query
	query := r.db.WithContext(ctx).Model(&models.Media{})

	// Apply filters
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("file_name ILIKE ?", "%"+search+"%")
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

		result := query.Limit(perPage).Offset(offset).Find(&media)
		if result.Error != nil {
			return nil, result.Error
		}

		// Return paginated result
		return map[string]interface{}{
			"data":         media,
			"total":        total,
			"per_page":     perPage,
			"current_page": page,
			"last_page":    (int(total) + perPage - 1) / perPage,
		}, nil
	}

	// Return all results without pagination
	result := query.Find(&media)
	if result.Error != nil {
		return nil, result.Error
	}

	return media, nil
}

// FindMedia finds a media by ID
func (r *MediaRepository) FindMedia(ctx context.Context, id string) (models.Media, error) {
	var media models.Media
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&media)
	if result.Error != nil {
		return models.Media{}, result.Error
	}
	return media, nil
}

// CreateMedia creates a new media
func (r *MediaRepository) CreateMedia(ctx context.Context, data map[string]interface{}) (models.Media, error) {
	// Create a new media instance
	media := models.Media{
		FileID:      data["file_id"].(string),
		FileName:    data["file_name"].(string),
		URL:         data["url"].(string),
		ThumbURL:    data["thumb_url"].(string),
		FileType:    data["file_type"].(string),
		Size:        data["size"].(int64),
		Description: data["description"].(string),
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return models.Media{}, tx.Error
	}

	// Defer a rollback in case anything fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the media within the transaction
	if err := tx.Create(&media).Error; err != nil {
		tx.Rollback()
		return models.Media{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return models.Media{}, err
	}

	// Return the created media
	return media, nil
}

// UpdateMedia updates an existing media
func (r *MediaRepository) UpdateMedia(ctx context.Context, data map[string]interface{}, id string) (models.Media, error) {
	// Find the media first
	media, err := r.FindMedia(ctx, id)
	if err != nil {
		return models.Media{}, err
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return models.Media{}, tx.Error
	}

	// Defer a rollback in case anything fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update the media within the transaction
	if err := tx.Model(&media).Updates(data).Error; err != nil {
		tx.Rollback()
		return models.Media{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return models.Media{}, err
	}

	// Refresh the media data
	return r.FindMedia(ctx, id)
}

// DeleteMedia soft deletes a media
func (r *MediaRepository) DeleteMedia(ctx context.Context, id string) error {
	// Find the media first
	media, err := r.FindMedia(ctx, id)
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

	// Delete the media within the transaction
	if err := tx.Delete(&media).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// FindMediaByFileID finds a media by file ID
func (r *MediaRepository) FindMediaByFileID(ctx context.Context, fileID string) (models.Media, error) {
	var media models.Media
	result := r.db.WithContext(ctx).Where("file_id = ?", fileID).First(&media)
	if result.Error != nil {
		return models.Media{}, result.Error
	}
	return media, nil
}