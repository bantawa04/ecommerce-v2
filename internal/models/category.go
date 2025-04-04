package models

import (
	"time"

	"beautyessentials.com/internal/constant"
	"beautyessentials.com/internal/utils"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// StatusEnum represents the status of a brand

// Brand represents a brand in the system
type Category struct {
	ID        string              `json:"id" gorm:"primaryKey;type:char(26)"`
	Name      string              `json:"name" gorm:"type:varchar(255);not null"`
	Slug      string              `json:"slug" gorm:"type:varchar(255);not null;uniqueIndex"`
	Status    constant.StatusEnum `json:"status" gorm:"type:enum('active','inactive');default:active"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	DeletedAt gorm.DeletedAt      `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeCreate will set a ULID rather than numeric ID and generate a slug
func (b *Category) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		// Generate a new ULID
		id := ulid.Make()
		b.ID = id.String()
	}

	// Generate slug from name if not provided
	if b.Slug == "" && b.Name != "" {
		b.Slug = utils.GenerateSlug(b.Name)
	}

	return nil
}

// TableName specifies the table name for the Brand model
func (Category) TableName() string {
	return "categories"
}
