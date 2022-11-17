package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type UpdateProfilePasswordRequest struct {
	Password string `json:"password" binding:"required,secure_password"`
}

type UpdateProfilePasswordResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUpdateProfilePasswordResponse(source *models.User) (target *UpdateProfilePasswordResponse) {
	target = &UpdateProfilePasswordResponse{
		ID:        source.ID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
	return
}
