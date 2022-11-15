package models

import (
	"time"

	"gorm.io/gorm"
)

type SubDistrict struct {
	ID         uint
	Name       string
	DistrictId uint
	District   District
	ZipCode    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
