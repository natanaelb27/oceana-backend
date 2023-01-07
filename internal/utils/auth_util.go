package utils

import (
	internalDtos "github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func AuthenticateUser(user *models.User, providedPassword string) (*internalDtos.UserLoginData, error) {
	err := CheckHash(user.Password, providedPassword)
	if err != nil {
		return nil, e.ErrInvalidPassword
	}
	token, err := GenerateUserTokenString(user)
	if err != nil {
		return nil, err
	}
	userLoginData := &internalDtos.UserLoginData{
		AccessToken: token,
		Role:        user.Role,
	}
	return userLoginData, nil
}

func GenerateUserTokenString(user *models.User) (string, error) {
	jwtWrapper, err := NewJwtWrapper()
	if err != nil {
		return "", err
	}
	jwtUser := internalDtos.NewJwtUser(user)
	token, err := jwtWrapper.GenerateToken(jwtUser, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateRefreshTokenString(user *models.User) (string, error) {
	jwtWrapper, err := NewRefreshTokenJwtWrapper()
	if err != nil {
		return "", err 
	}
	jwtUser := internalDtos.NewJwtUser(user)
	token, err := jwtWrapper.GenerateToken(jwtUser, jwtWrapper)
	if err != nil {
		return "", err
	}
	return token, nil
}
