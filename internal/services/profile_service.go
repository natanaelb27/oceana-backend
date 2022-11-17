package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileService interface {
	UpdateProfilePassword(c *gin.Context, id uuid.UUID, request *externalDtos.UpdateProfilePasswordRequest) (model *models.User, err error)
}

type profileService struct {
	userRepository  r.UserRepository
	redisRepository r.RedisRepository
}

type ProfileServiceConfig struct {
	UserRepository  r.UserRepository
	RedisRepository r.RedisRepository
}

func NewProfileService(c *ProfileServiceConfig) ProfileService {
	return &profileService{
		userRepository:  c.UserRepository,
		redisRepository: c.RedisRepository,
	}
}

func (u *profileService) UpdateProfilePassword(c *gin.Context, id uuid.UUID, request *externalDtos.UpdateProfilePasswordRequest) (model *models.User, err error) {
	user, err := u.userRepository.GetByID(id)
	if err != nil {
		return
	}

	if !utils.IsPasswordContainsUsername(request.Password, user.Username) {
		err = e.ErrPasswordCannotContainsUsername
		return
	}

	err = utils.CheckHash(user.Password, request.Password)
	if err == nil {
		err = e.ErrPasswordCannotBeTheSameAsPreviousPassword
		return
	}

	hashedPassword, err := utils.Hash(request.Password)
	if err != nil {
		return
	}

	user.Password = hashedPassword

	// blacklist all refreshToken and accessToken here

	return u.userRepository.Update(id, user)
}
