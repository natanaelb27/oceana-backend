package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type RoleCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

func (source *RoleCreateRequest) ToModel() (target *models.Role) {
	target = &models.Role{
		Name: source.Name,
	}
	return
}

type RoleCreateResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewRoleCreateResponse(source *models.Role) (target *RoleCreateResponse) {
	target = &RoleCreateResponse{
		ID:        source.ID,
		Name:      source.Name,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
	return
}

type RoleGetResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewRoleGetResponse(source *models.Role) (target *RoleGetResponse) {
	target = &RoleGetResponse{
		ID:        source.ID,
		Name:      source.Name,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
	return
}

type RoleDeleteResponse struct {
	ID uuid.UUID `json:"id"`
}

func NewRoleDeleteResponse(source *models.Role) (target *RoleDeleteResponse) {
	target = &RoleDeleteResponse{
		ID: source.ID,
	}
	return
}
