package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uint
	UserID    uuid.UUID
	User      User
	Code      string
	ExpiresAt time.Time
}
