package models

import "github.com/google/uuid"

type Address struct {
	ID            uuid.UUID
	UserId        uuid.UUID
	User          User
	Name          string
	Detail        string
	SubDistrictId uint
	SubDistrict   SubDistrict
	IsDefault     bool
}
