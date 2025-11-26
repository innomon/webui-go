package models

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Title     string         `gorm:"not null" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Messages  []Message      `gorm:"foreignKey:ChatID" json:"messages,omitempty"`
}

type Message struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ChatID    uint           `gorm:"not null" json:"chat_id"`
	Role      string         `gorm:"not null" json:"role"` // e.g., "user", "assistant"
	Content   string         `gorm:"not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
