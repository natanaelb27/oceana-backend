package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint `gorm:"primary_key"`
	Name      string
	ParentID  uint
	Category  *Category `gorm:"foreignKey:ParentID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
