package models

import (
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
}

func (m *Address) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
