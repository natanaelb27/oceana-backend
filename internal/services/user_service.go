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
	userRepository  r.UserRepository
	emailRepository r.EmailRepository
}

type UserServiceConfig struct {
	UserRepository  r.UserRepository
	EmailRepository r.EmailRepository
}

func NewUserService(c *UserServiceConfig) UserService {
	return &userService{
		userRepository:  c.UserRepository,
		emailRepository: c.EmailRepository,
	}
}

func (u *userService) UpdateProfile(userProfileUpdateRequest *dtos.UserProfileUpdateRequest) (*dtos.UserProfileUpdateResponse, error) {
	user, err := u.userRepository.GetByID(userProfileUpdateRequest.ID)
	if err != nil {
		return nil, err
	}
	if userProfileUpdateRequest.FullName != nil {
		user.FullName = *userProfileUpdateRequest.FullName
	}
	if userProfileUpdateRequest.Gender != nil {
		user.Gender = userProfileUpdateRequest.Gender
	}
	if userProfileUpdateRequest.Birthdate != nil {
		user.Birthdate = userProfileUpdateRequest.Birthdate
	}
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
	email, err := u.emailRepository.GetById(user.EmailID)
	if err != nil {
		return nil, err
	}
	responseData.Email = email.Address
	return responseData, nil
}
