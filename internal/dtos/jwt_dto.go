package dtos

import (
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/models"
)

type JwtUser struct {
	ID uuid.UUID `json:"id"`
}

func NewJwtUser(source *models.User) *JwtUser {
	return &JwtUser{
		ID: source.ID,
	}
}
