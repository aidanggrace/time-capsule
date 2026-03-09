package models

import (
	"time"

	"gorm.io/gorm"
)

type Capsule struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	OwnerID        string         `json:"owner_id" gorm:"not_null"`
	RecipientEmail string         `json:"recipient_email"`
	Message        string         `json:"message"`
	MediaURL       string         `json:"media_url"`
	UnlockAt       time.Time      `json:"unlock_at" gorm:"not_null"`
	DeliveredAt    *time.Time     `json:"delivered_at"` // pointer bc it can be null
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
