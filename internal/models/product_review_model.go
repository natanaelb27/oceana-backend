package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductReview struct {
	ID            uint
	UserId        uint
	User          User
	OrderDetailId int
	OrderDetail   OrderDetail
	Rating        uint
	Review        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
