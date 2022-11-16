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
	jwtWrapper, err := NewJwtWrapper()
	if err != nil {
		return nil, err
	}
	jwtUser := internalDtos.NewJwtUser(user)
	token, err := jwtWrapper.GenerateToken(jwtUser)
	if err != nil {
		return nil, err
	}
	userLoginResponse := &internalDtos.UserLoginData{
		AccessToken: token,
		Role:        user.Role,
	}
	return userLoginResponse, nil
}
