package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	User          User
	Name          string
	Detail        string
	SubDistrictID uint
	SubDistrict   SubDistrict
	IsDefault     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
