package models

import (
	"time"

	"gorm.io/gorm"
)

// Knowledge represents a knowledge base in the database
type Knowledge struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	Name         string         `gorm:"uniqueIndex;not null" json:"name"`
	Description  string         `json:"description"`
	CollectionName string       `gorm:"not null" json:"collection_name"` // Name of the vector database collection
	FileIDs      []byte         `gorm:"type:jsonb" json:"file_ids"`      // JSONB array of file IDs
	AccessControl []byte        `gorm:"type:jsonb" json:"access_control"` // JSONB object for access control
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// KnowledgeForm for creating and updating a knowledge base
type KnowledgeForm struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// KnowledgeResponse for returning knowledge base details
type KnowledgeResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// KnowledgeUserResponse for returning knowledge base details with associated files
type KnowledgeUserResponse struct {
	KnowledgeResponse
	Files []File `json:"files"`
}
