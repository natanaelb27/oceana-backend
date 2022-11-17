package services

import (
	"time"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(userCreateRequest *externalDtos.UserCreateRequest) (*externalDtos.UserCreateResponse, error)
	Login(userLoginRequest *externalDtos.UserLoginRequest) (*internalDtos.UserLoginData, error)
	StepUpPassword(c *gin.Context, id uuid.UUID, request *externalDtos.AuthStepUpPasswordRequest) (token string, err error)
}

type authService struct {
	userRepository         r.UserRepository
	emailRepository        r.EmailRepository
	walletRepository       r.WalletRepository
	refreshTokenRepository r.RefreshTokenRepository
	redisRepository        r.RedisRepository
}

type AuthServiceConfig struct {
	UserRepository         r.UserRepository
	EmailRepository        r.EmailRepository
	WalletRepository       r.WalletRepository
	RefreshTokenRepository r.RefreshTokenRepository
	RedisRepository        r.RedisRepository
}

func NewAuthService(c *AuthServiceConfig) AuthService {
	return &authService{
		userRepository:         c.UserRepository,
		emailRepository:        c.EmailRepository,
		walletRepository:       c.WalletRepository,
		refreshTokenRepository: c.RefreshTokenRepository,
		redisRepository:        c.RedisRepository,
	}
}

func (u *authService) Register(userCreateRequest *externalDtos.UserCreateRequest) (*externalDtos.UserCreateResponse, error) {
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
	responseData := externalDtos.NewUserCreateResponse(createdUser)
	responseData.Email = createdEmail.Address
	return responseData, nil
}

func (u *authService) Login(userLoginRequest *externalDtos.UserLoginRequest) (*internalDtos.UserLoginData, error) {
	user, err := u.userRepository.GetByUsername(userLoginRequest.LoginCredentials)
	if err == nil {
		userLoginData, err := utils.AuthenticateUser(user, userLoginRequest.Password)
		if err != nil {
			return nil, err
		}
		refreshToken := utils.NewRefreshToken(user)
		_, err = u.refreshTokenRepository.Create(refreshToken)
		if err != nil {
			return nil, err
		}
		userLoginData.RefreshToken = refreshToken.Code
		userLoginData.ExpiresAt = refreshToken.ExpiresAt
		return userLoginData, nil
	}
	user, err = u.userRepository.GetByPhoneNumber(userLoginRequest.LoginCredentials)
	if err == nil {
		userLoginData, err := utils.AuthenticateUser(user, userLoginRequest.Password)
		if err != nil {
			return nil, err
		}
		refreshToken := utils.NewRefreshToken(user)
		_, err = u.refreshTokenRepository.Create(refreshToken)
		if err != nil {
			return nil, err
		}
		userLoginData.RefreshToken = refreshToken.Code
		userLoginData.ExpiresAt = refreshToken.ExpiresAt
		return userLoginData, nil
	}
	email, err := u.emailRepository.GetByEmail(userLoginRequest.LoginCredentials)
	if err == nil {
		user, err = u.userRepository.GetByEmailID(email.ID)
		if err != nil {
			return nil, err
		}
		userLoginData, err := utils.AuthenticateUser(user, userLoginRequest.Password)
		if err != nil {
			return nil, err
		}
		refreshToken := utils.NewRefreshToken(user)
		_, err = u.refreshTokenRepository.Create(refreshToken)
		if err != nil {
			return nil, err
		}
		userLoginData.RefreshToken = refreshToken.Code
		userLoginData.ExpiresAt = refreshToken.ExpiresAt
		return userLoginData, nil
	}
	return nil, e.ErrUserNotFound
}

func (u *authService) StepUpPassword(c *gin.Context, id uuid.UUID, request *externalDtos.AuthStepUpPasswordRequest) (token string, err error) {
	user, err := u.userRepository.GetByID(id)
	if err != nil {
		return
	}
	err = utils.CheckHash(user.Password, request.Password)
	if err != nil {
		err = e.ErrInvalidPassword
		return
	}

	jwtWrapper, err := utils.NewJwtWrapper()
	if err != nil {
		return
	}
	jwtData := internalDtos.NewJwtStepUp("level1")
	token, err = jwtWrapper.JwtStepUpGenerateToken(jwtData)
	if err != nil {
		return
	}

	err = u.redisRepository.Set(c, "upgrade-token-password:"+id.String(), token, 10*time.Minute)
	if err != nil {
		return
	}

	return
}

func (u *authService) SendResetPasswordEmail(c *gin.Context, emailAddress string) error {
	email, err := u.emailRepository.GetByEmail(emailAddress)
	if err != nil {
		return e.ErrEmailNotFound
	}
	user, err := u.userRepository.GetByEmailID(email.ID)
	if err != nil {
		return e.ErrUserNotFound
	}
	err = u.redisRepository.Set(c, "reset-password:"+user.ID.String(), uuid.New().String(), 15*time.Minute)
	if err != nil {
		return err
	}
	return nil
}

