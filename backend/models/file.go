package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	FolderID  *uint          `json:"folder_id"`
	Name      string         `gorm:"not null" json:"name"`
	Path      string         `gorm:"not null" json:"path"` // Stored path on the server
	MimeType  string         `json:"mime_type"`
	Size      int64          `json:"size"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Folder Folder `gorm:"foreignKey:FolderID" json:"-"`
}

type Folder struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	ParentID  *uint          `json:"parent_id"`
	Name      string         `gorm:"not null" json:"name"`	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Parent *Folder `gorm:"foreignKey:ParentID" json:"-"`
	Files  []File  `gorm:"foreignKey:FolderID" json:"files,omitempty"`
	Subfolders []Folder `gorm:"foreignKey:ParentID" json:"subfolders,omitempty"`
}
