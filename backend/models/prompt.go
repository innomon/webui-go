package models

import (
	"time"

	"gorm.io/gorm"
)

// Prompt represents a prompt in the database
type Prompt struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"not null" json:"content"`
	Command   string         `gorm:"uniqueIndex;not null" json:"command"` // e.g., "/summarize"
	AccessControl []byte        `gorm:"type:jsonb" json:"access_control"` // JSONB object for access control
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PromptForm for creating and updating a prompt
type PromptForm struct {
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content" binding:"required"`
	Command       string `json:"command" binding:"required"`
	AccessControl []byte `gorm:"type:jsonb" json:"access_control"`
}

// PromptUserResponse for returning prompt details with user-specific info
type PromptUserResponse struct {
	ID        uint           `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Command   string         `json:"command"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
