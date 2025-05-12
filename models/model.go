package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	*gorm.Model
	Context string `json:"content"`
}

type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
