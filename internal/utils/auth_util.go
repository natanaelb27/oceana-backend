package utils

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

func AuthenticateUser(user *models.User, providedPassword string) (*dtos.UserLoginResponse, error) {
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
	userLoginResponse := &dtos.UserLoginResponse{
		AccessToken: token,
	}
	return userLoginResponse, nil
}
