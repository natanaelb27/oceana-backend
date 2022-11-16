package models

import "github.com/google/uuid"

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
