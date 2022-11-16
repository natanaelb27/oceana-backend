package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/google/uuid"
)

type UserService interface {
	UpdateProfile(userProfileUpdateRequest *dtos.UserProfileUpdateRequest) (*dtos.UserProfileUpdateResponse, error)
	GetProfile(id uuid.UUID) (*dtos.UserProfileGetResponse, error)
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

func (u *userService) GetProfile(id uuid.UUID) (*dtos.UserProfileGetResponse, error) {
	user, err := u.userRepository.GetByID(id)
	if err != nil {
		return nil, err 
	}
	responseData := dtos.NewUserProfileGetResponse(user) 
	return responseData, nil
}
