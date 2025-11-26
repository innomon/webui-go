package models

import (
	"time"

	"gorm.io/gorm"
)

// Tool represents a tool in the database
type Tool struct {
	ID        string         `gorm:"primarykey" json:"id"`
	UserID    uint           `json:"user_id"`
	Name      string         `gorm:"not null" json:"name"`
	Content   string         `gorm:"not null" json:"content"` // Python code for the tool
	Specs     []byte         `gorm:"type:jsonb" json:"specs"`     // JSONB object for tool specifications
	Meta      []byte         `gorm:"type:jsonb" json:"meta"`      // JSONB object for tool metadata
	AccessControl []byte        `gorm:"type:jsonb" json:"access_control"` // JSONB object for access control
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ToolForm for creating and updating a tool
type ToolForm struct {
	ID      string `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
	Meta    []byte `gorm:"type:jsonb" json:"meta"`
	AccessControl []byte `gorm:"type:jsonb" json:"access_control"`
}

// ToolResponse for returning tool details
type ToolResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Specs     []byte    `json:"specs"`
	Meta      []byte    `json:"meta"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToolUserResponse for returning tool details with user-specific info
type ToolUserResponse struct {
	ToolResponse
	HasUserValves bool `json:"has_user_valves"`
}
