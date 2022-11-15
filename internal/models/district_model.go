package models

import "time"

type District struct {
	ID        uint
	Name      string
	CityId    uint
	City      City
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
