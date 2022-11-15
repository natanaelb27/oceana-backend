package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type AuthService interface {
	Register(userCreateRequest *dtos.UserCreateRequest) (*dtos.UserCreateResponse, error)
	Login(userLoginRequest *dtos.UserLoginRequest) (*dtos.UserLoginResponse, error)
}

type authService struct {
	userRepository         r.UserRepository
	emailRepository        r.EmailRepository
	walletRepository       r.WalletRepository
	refreshTokenRepository r.RefreshTokenRepository
}

type AuthServiceConfig struct {
	UserRepository         r.UserRepository
	EmailRepository        r.EmailRepository
	WalletRepository       r.WalletRepository
	RefreshTokenRepository r.RefreshTokenRepository
}

func NewAuthService(c *AuthServiceConfig) AuthService {
	return &authService{
		userRepository:         c.UserRepository,
		emailRepository:        c.EmailRepository,
		walletRepository:       c.WalletRepository,
		refreshTokenRepository: c.RefreshTokenRepository,
	}
}

func (u *authService) Register(userCreateRequest *dtos.UserCreateRequest) (*dtos.UserCreateResponse, error) {
	_, err := u.emailRepository.GetByEmail(userCreateRequest.Email)
	if err != e.ErrDataNotFound {
		return nil, e.ErrEmailAddressAlreadyExist
	}
	_, err = u.userRepository.GetByUsername(userCreateRequest.Username)
	if err != e.ErrDataNotFound {
		return nil, e.ErrUsernameAlreadyExist
	}
	_, err = u.userRepository.GetByPhoneNumber(userCreateRequest.PhoneNumber)
	if err != e.ErrDataNotFound {
		return nil, e.ErrPhoneNumberAlreadyExist
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
	hashedPassword, err := utils.Hash(userCreateRequest.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	createdUser, err := u.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	wallet := &models.Wallet{
		UserID:          createdUser.ID,
		Balance:         0,
		Pin:             "",
		IsActive:        false,
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

func (u *authService) Login(userLoginRequest *dtos.UserLoginRequest) (*dtos.UserLoginResponse, error) {
	user, err := u.userRepository.GetByUsername(userLoginRequest.LoginCredentials)
	if err == nil {
		userLoginResponse, err := utils.AuthenticateUser(user, userLoginRequest.Password)
		if err != nil {
			return nil, err
		}
		refreshToken := utils.NewRefreshToken(user)
		_, err = u.refreshTokenRepository.Create(refreshToken)
		if err != nil {
			return nil, err
		}
		userLoginResponse.RefreshToken = refreshToken.Code
		return userLoginResponse, nil
	}
	user, err = u.userRepository.GetByPhoneNumber(userLoginRequest.LoginCredentials)
	if err == nil {
		userLoginResponse, err := utils.AuthenticateUser(user, userLoginRequest.Password)
		if err != nil {
			return nil, err
		}
		refreshToken := utils.NewRefreshToken(user)
		_, err = u.refreshTokenRepository.Create(refreshToken)
		if err != nil {
			return nil, err
		}
		userLoginResponse.RefreshToken = refreshToken.Code
		return userLoginResponse, nil
	}
	email, err := u.emailRepository.GetByEmail(userLoginRequest.LoginCredentials)
	if err == nil {
		user, err = u.userRepository.GetByEmailID(email.ID)
		if err != nil {
			return nil, err
		}
		userLoginResponse, err := utils.AuthenticateUser(user, userLoginRequest.Password)
		if err != nil {
			return nil, err
		}
		refreshToken := utils.NewRefreshToken(user)
		_, err = u.refreshTokenRepository.Create(refreshToken)
		if err != nil {
			return nil, err
		}
		userLoginResponse.RefreshToken = refreshToken.Code
		return userLoginResponse, nil
	}
	return nil, e.ErrUserNotFound
}
