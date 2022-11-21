package dtos

import (
	"mime/multipart"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type UserCreateRequest struct {
	FullName    string                `form:"full_name" json:"full_name" binding:"required"`
	Username    string                `form:"username" json:"username" binding:"required"`
	Email       string                `form:"email" json:"email" binding:"required,email"`
	Password    string                `form:"password" json:"password" binding:"required"`
	PhoneNumber string                `form:"phone_number" json:"phone_number" binding:"required"`
	Gender      *string               `form:"gender" json:"gender"`
	Birthdate   *time.Time            `form:"birthdate" json:"birthdate"`
	Photo       *multipart.FileHeader `form:"photo" json:"photo"`
	GoogleID    *string               `form:"google_id" json:"google_id"`
	Role        string                `form:"role" json:"role" binding:"required"`
}

func (source *UserCreateRequest) ToModel() (target *models.User) {
	target = &models.User{
		FullName:    source.FullName,
		Username:    source.Username,
		PhoneNumber: source.PhoneNumber,
		Gender:      source.Gender,
		Birthdate:   source.Birthdate,
		GoogleID:    source.GoogleID,
		Role:        source.Role,
	}
	return
}

type UserCreateResponse struct {
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

type UserLoginRequest struct {
	LoginCredentials string `json:"login_credentials" binding:"required"`
	Password         string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Role string `json:"role"`
}

type AuthStepUpPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type AuthStepUpPasswordResponse struct {
	UpgradeToken string `json:"upgrade_token"`
}

func NewAuthStepUpPasswordResponse(token string) (target *AuthStepUpPasswordResponse) {
	target = &AuthStepUpPasswordResponse{
		UpgradeToken: token,
	}
	return
}

type UserVerifyResetPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type UserVerifyResetPasswordResponse struct {
	Status bool `json:"status"`
}

func NewUserVerifyResetPasswordResponse(status bool) (target *UserVerifyResetPasswordResponse) {
	target = &UserVerifyResetPasswordResponse{
		Status: status,
	}
	return
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type ResetPasswordResponse struct {
	Status bool `json:"status"`
}

func NewResetPasswordResponse(status bool) (target *ResetPasswordResponse) {
	target = &ResetPasswordResponse{
		Status: status,
	}
	return
}