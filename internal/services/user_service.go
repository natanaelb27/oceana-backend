package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
)

type UserService interface {
	UpdateProfile(userProfileUpdateRequest *dtos.UserProfileUpdateRequest) (*dtos.UserProfileUpdateResponse, error)
}

type userService struct {
	userRepository r.UserRepository
}

type UserServiceConfig struct {
	UserRepository r.UserRepository
}

func NewUserService(c *UserServiceConfig) UserService {
	return &userService{
		userRepository: c.UserRepository,
	}
}

func (u *userService) UpdateProfile(userProfileUpdateRequest *dtos.UserProfileUpdateRequest) (*dtos.UserProfileUpdateResponse, error) {
	user := userProfileUpdateRequest.ToModel()
	if userProfileUpdateRequest.Photo != nil {
		photoFile, err := userProfileUpdateRequest.Photo.Open()
		if err != nil {
			return nil, err
		}
		imageUrl, err := utils.UploadImage(photoFile)
		if err != nil {
			return nil, err
		}
		user.Photo = &imageUrl
	}
	updatedUser, err := u.userRepository.Update(userProfileUpdateRequest.ID, user)
	if err != nil {
		return nil, err 
	}
	responseData := dtos.NewUserProfileUpdateResponse(updatedUser)
	return responseData, nil
}
