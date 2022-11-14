package dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type JwtUser struct {
	ID uuid.UUID `json:"id"`
}

func NewJwtUser(source *models.User) *JwtUser {
	return &JwtUser{
		ID: source.ID,
	}
}
