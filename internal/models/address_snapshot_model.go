package models

import (
	"time"

	"gorm.io/gorm"
)

type AddressSnapshot struct {
	ID            uint
	Address       string `json:"address"`
	SubDistrictID uint   `json:"sub_district_id"`
	SubDistrict   SubDistrict
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
