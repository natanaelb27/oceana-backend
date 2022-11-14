package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        int `gorm:"primary_key"`
	Name      string
	ParentId  int
	Category  *Category `gorm:"foreignKey:ParentId;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
