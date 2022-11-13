package dtos

import "git.garena.com/azmi-adhani/group-project-backend/internal/models"

type JwtUser struct {
	ID string `json:"id"`
}

func NewJwtUser(source *models.User) *JwtUser {
	return &JwtUser{
		ID: source.ID,
	}
}
