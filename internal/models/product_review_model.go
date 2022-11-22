package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductReview struct {
	ID            uint
	UserId        uint
	User          User
	OrderDetailId uuid.UUID
	OrderDetail   OrderDetail
	Rating        uint
	Review        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
