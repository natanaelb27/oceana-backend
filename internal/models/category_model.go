package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        int            `json:"id" gorm:"primary_key"`
	Name      string         `json:"name"`
	ParentId  int            `json:"parent_id"`
	Category  *Category      `json:"category" gorm:"foreignKey:ParentId;references:ID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
