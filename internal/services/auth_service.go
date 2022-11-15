package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type AuthService interface {
	Create(userCreateRequest *dtos.UserCreateRequest) (*dtos.UserCreateResponse, error)
}

type authService struct {
	userRepository  r.UserRepository
	emailRepository r.EmailRepository
	walletRepository r.WalletRepository
}

type AuthServiceConfig struct {
	UserRepository  r.UserRepository
	EmailRepository r.EmailRepository
	WalletRepository r.WalletRepository
}

func NewAuthService(c *AuthServiceConfig) AuthService {
	return &authService{
		userRepository:  c.UserRepository,
		emailRepository: c.EmailRepository,
		walletRepository: c.WalletRepository,
	}
}

func (u *authService) Create(userCreateRequest *dtos.UserCreateRequest) (*dtos.UserCreateResponse, error) {
	_, err := u.emailRepository.GetByEmail(userCreateRequest.Email)
	if err != e.ErrDataNotFound {
		return nil, e.ErrEmailAddressAlreadyExist
	}
	email := &models.Email{
		Address: userCreateRequest.Email,
	}
	createdEmail, err := u.emailRepository.Create(email)
	if err != nil {
		return nil, err
	}
	user := userCreateRequest.ToModel()
	user.EmailID = email.ID
	if userCreateRequest.Photo != nil {
		photoFile, err := userCreateRequest.Photo.Open()
		if err != nil {
			return nil, err
		}
		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}
		user.Photo = &imageUrl
	}
	hashedPassword, err := utils.HashPassword(userCreateRequest.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	createdUser, err := u.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	wallet := &models.Wallet{
		UserID: createdUser.ID,
		Balance: 0,
		Pin: "",
		IsActive: false,
		InvalidPinCount: 0,
	}
	_, err = u.walletRepository.Create(wallet)
	if err != nil {
		return nil, err
	}
	responseData := dtos.NewUserCreateResponse(createdUser)
	responseData.Email = createdEmail.Address
	return responseData, nil
}
