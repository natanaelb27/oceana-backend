package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
)

type AuthService interface {
	Create(userCreateRequest *dtos.UserCreateRequest) (*dtos.UserCreateResponse, error)
}

type authService struct {
	userRepository r.UserRepository
	emailRepository r.EmailRepository
}

type AuthServiceConfig struct {
	UserRepository r.UserRepository
	EmailRepository r.EmailRepository
}

func NewAuthService(c *AuthServiceConfig) AuthService {
	return &authService{
		userRepository: c.UserRepository,
		emailRepository: c.EmailRepository,
	}
}

func (u *authService) Create(userCreateRequest *dtos.UserCreateRequest) (*dtos.UserCreateResponse, error) {
	email := &models.Email{
		Address: userCreateRequest.Email,
	}
	email, err := u.emailRepository.Create(email)
	if err != nil {
		return nil, err
	}
	user := userCreateRequest.ToModel()
	user.EmailID = email.ID.String()
	hashedPassword, err := utils.HashPassword(userCreateRequest.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	createdUser, err := u.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	responseData := dtos.NewUserCreateResponse(createdUser)
	responseData.Email = email.Address
	return responseData, nil
}
