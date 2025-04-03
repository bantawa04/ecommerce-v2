package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// StatusEnum represents the status of a brand
type StatusEnum string

const (
	StatusActive   StatusEnum = "active"
	StatusInactive StatusEnum = "inactive"
)

// Brand represents a brand in the system
type Brand struct {
	ID        string     `json:"id" gorm:"primaryKey;type:char(26)"`
	Name      string     `json:"name" gorm:"type:varchar(255);not null"`
	Slug      string     `json:"slug" gorm:"type:varchar(255);not null;uniqueIndex"`
	Status    StatusEnum `json:"status" gorm:"type:enum('active','inactive');default:active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeCreate will set a ULID rather than numeric ID
func (b *Brand) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		// Generate a new ULID
		id := ulid.Make()
		b.ID = id.String()
	}
	return nil
}

// TableName specifies the table name for the Brand model
func (Brand) TableName() string {
	return "brands"
}