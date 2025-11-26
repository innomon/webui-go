package models

import (
	"time"

	"gorm.io/gorm"
)

// Model represents an AI model in the database
type Model struct {
	ID          string         `gorm:"primarykey" json:"id"`
	UserID      uint           `json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	BaseModelID string         `json:"base_model_id"`
	Meta        []byte         `gorm:"type:jsonb" json:"meta"`        // JSONB object for model metadata (e.g., tags, profile_image_url)
	Params      []byte         `gorm:"type:jsonb" json:"params"`      // JSONB object for model parameters
	AccessControl []byte        `gorm:"type:jsonb" json:"access_control"` // JSONB object for access control
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// ModelForm for creating and updating a model
type ModelForm struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	BaseModelID string `json:"base_model_id"`
	Meta        []byte `gorm:"type:jsonb" json:"meta"`
	Params      []byte `gorm:"type:jsonb" json:"params"`
	AccessControl []byte `gorm:"type:jsonb" json:"access_control"`
	IsActive    bool   `json:"is_active"`
}

// ModelResponse for returning model details
type ModelResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	BaseModelID string    `json:"base_model_id"`
	Meta        []byte    `json:"meta"`
	Params      []byte    `json:"params"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ModelListResponse for returning a list of models with total count
type ModelListResponse struct {
	Models []ModelResponse `json:"models"`
	Total  int64           `json:"total"`
}
