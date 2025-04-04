package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Media represents a media file in the system
type Media struct {
	ID          string    `json:"id" gorm:"primaryKey;type:char(26)"`
	FileID      string    `json:"file_id" gorm:"type:varchar(255)"`	
	URL         string    `json:"url" gorm:"type:varchar(255)"`
	ThumbURL    string    `json:"thumb_url" gorm:"type:varchar(255)"`	
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate will set a ULID rather than numeric ID
func (m *Media) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		// Generate a new ULID
		id := ulid.Make()
		m.ID = id.String()
	}
	return nil
}

// TableName specifies the table name for the Media model
func (Media) TableName() string {
	return "medias"
}
