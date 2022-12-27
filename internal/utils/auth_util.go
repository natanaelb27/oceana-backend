package utils

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
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
