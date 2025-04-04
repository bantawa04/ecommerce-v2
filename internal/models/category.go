package models

import (
	"time"

	"beautyessentials.com/internal/constant"
	"beautyessentials.com/internal/utils"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Category represents a category in the system
type Category struct {
	ID        string              `json:"id" gorm:"primaryKey;type:char(26)"`
	Name      string              `json:"name" gorm:"type:varchar(255);not null"`
	Slug      string              `json:"slug" gorm:"type:varchar(255);not null;uniqueIndex"`
	Status    constant.StatusEnum `json:"status" gorm:"type:enum('active','inactive');default:active"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	DeletedAt gorm.DeletedAt      `json:"deleted_at,omitempty" gorm:"index"`
	// Media       []Media             `json:"media,omitempty" gorm:"many2many:mediables;foreignKey:ID;joinForeignKey:mediable_id;joinReferences:media_id;where:mediable_type = 'App\\Model\\Category'"`
}

// BeforeCreate will set a ULID rather than numeric ID and generate a slug
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		// Generate a new ULID
		id := ulid.Make()
		c.ID = id.String()
	}

	// Generate slug from name if not provided
	if c.Slug == "" && c.Name != "" {
		c.Slug = utils.GenerateSlug(c.Name)
	}

	return nil
}

// TableName specifies the table name for the Category model
func (Category) TableName() string {
	return "categories"
}
