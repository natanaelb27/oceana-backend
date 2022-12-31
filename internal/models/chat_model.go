package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID         uint
	FromUserID uuid.UUID
	FromUser   User
	ToUserID   uuid.UUID
	ToUser     User
	Message    string
	MediaURL   string
	// Read       bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
