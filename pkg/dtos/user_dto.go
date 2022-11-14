package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type UserCreateRequest struct {
	FullName    string    `json:"full_name" binding:"required"`
	Username    string    `json:"username" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Password    string    `json:"password" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Gender      string    `json:"gender" binding:"required"`
	Birthdate   time.Time `json:"birthdate" binding:"required"`
	Photo       *string   `json:"photo"`
	GoogleID    *string   `json:"google_id"`
	Role        string    `json:"role" binding:"required"`
}

func (source *UserCreateRequest) ToModel() (target *models.User) {
	target = &models.User{
		FullName:    source.FullName,
		Username:    source.Username,
		PhoneNumber: source.PhoneNumber,
		Gender:      source.Gender,
		Birthdate:   source.Birthdate,
		Photo:       source.Photo,
		GoogleID:    source.GoogleID,
		Role:        source.Role,
	}
	return
}

type UserCreateResponse struct {
	ID          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Gender      string    `json:"gender"`
	Birthdate   time.Time `json:"birthdate"`
	Photo       *string   `json:"photo,omitempty"`
	GoogleID    *string   `json:"google_id,omitempty"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewUserCreateResponse(source *models.User) (target *UserCreateResponse) {
	target = &UserCreateResponse{
		ID:          source.ID,
		FullName:    source.FullName,
		Username:    source.Username,
		PhoneNumber: source.PhoneNumber,
		Gender:      source.Gender,
		Birthdate:   source.Birthdate,
		Photo:       source.Photo,
		GoogleID:    source.GoogleID,
		Role:        source.Role,
		CreatedAt:   source.CreatedAt,
		UpdatedAt:   source.UpdatedAt,
	}
	return
}
