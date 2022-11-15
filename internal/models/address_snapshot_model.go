package models

import (
	"time"

	"gorm.io/gorm"
)

type AddressSnapshot struct {
	ID            uint
	Address       string
	SubDistrictID uint
	SubDistrict   SubDistrict
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
