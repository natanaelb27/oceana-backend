package dtos

import (
	"mime/multipart"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type UserProfileUpdateRequest struct {
	ID        uuid.UUID
	FullName  *string               `form:"full_name" json:"full_name"`
	Photo     *multipart.FileHeader `form:"photo" json:"photo"`
	Gender    *string               `form:"gender" json:"gender"`
	Birthdate *time.Time            `form:"birthdate" json:"birthdate"`
}

type UserProfileUpdateResponse struct {
	ID        uuid.UUID  `json:"id"`
	FullName  *string    `json:"full_name,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
	Birthdate *time.Time `json:"birthdate,omitempty"`
	Photo     *string    `json:"photo,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewUserProfileUpdateResponse(source *models.User) (target *UserProfileUpdateResponse) {
	target = &UserProfileUpdateResponse{
		ID:        source.ID,
		FullName:  &source.FullName,
		Gender:    source.Gender,
		Birthdate: source.Birthdate,
		Photo:     source.Photo,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
	return
}

type UserProfileGetResponse struct {
	ID          uuid.UUID  `json:"id"`
	FullName    string     `json:"full_name"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Gender      *string    `json:"gender,omitempty"`
	Birthdate   *time.Time `json:"birthdate,omitempty"`
	Photo       *string    `json:"photo,omitempty"`
	GoogleID    *string    `json:"google_id,omitempty"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewUserProfileGetResponse(source *models.User) (target *UserProfileGetResponse) {
	target = &UserProfileGetResponse{
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

type UserProfileUpdateEmailRequest struct {
	ID    uuid.UUID
	Email string `form:"email" json:"email" binding:"required,email"`
}
