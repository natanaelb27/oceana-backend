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

type ProfileCheckPasswordVerificationCodeRequest struct {
	VerificationCode string `json:"verification_code" binding:"required"`
}

type ProfileCheckPasswordVerificationCodeResponse struct {
	Status bool `json:"status"`
}

func NewProfileCheckPasswordVerificationCodeResponse(status bool) (target *ProfileCheckPasswordVerificationCodeResponse) {
	target = &ProfileCheckPasswordVerificationCodeResponse{
		Status: status,
	}
	return
}

func NewPasswordVerificationCodeKey(userId string) string {
	return "password-verification-code:" + userId
}

type ProfileRequestChangePasswordVerificationCodeRequest struct {
	VerificationCode string `json:"verification_code" binding:"required"`
}

type ProfileRequestChangePasswordVerificationCodeResponse struct {
	Status bool `json:"status"`
}

func NewProfileRequestChangePasswordVerificationCodeResponse(status bool) (target *ProfileRequestChangePasswordVerificationCodeResponse) {
	target = &ProfileRequestChangePasswordVerificationCodeResponse{
		Status: status,
	}
	return
}
