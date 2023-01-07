package services

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/models"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	externalDtos "github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

type ProfileService interface {
	UpdateProfilePassword(c *gin.Context, id uuid.UUID, request *externalDtos.UpdateProfilePasswordRequest) (model *models.User, err error)
	CheckPasswordVerificationCode(c *gin.Context, id uuid.UUID, request *externalDtos.ProfileCheckPasswordVerificationCodeRequest) (status bool, err error)
	RequestChangePasswordVerificationCode(c *gin.Context, id uuid.UUID) (status bool, err error)
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

func (s *profileService) UpdateProfilePassword(c *gin.Context, id uuid.UUID, request *externalDtos.UpdateProfilePasswordRequest) (model *models.User, err error) {
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return
	}

	redisKey := externalDtos.NewPasswordVerificationCodeKey(id.String())
	verificationCode, err := s.redisRepository.Get(c, redisKey)
	if err != nil {
		err = e.ErrVerificationCodeNotFound
		return
	}
	if verificationCode != request.VerificationCode {
		err = e.ErrInvalidVerificationCode
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

	model, err = s.userRepository.Update(id, user)
	if err != nil {
		return
	}

	err = s.redisRepository.Delete(c, redisKey)
	if err != nil {
		return
	}

	// blacklist all refreshToken and accessToken here

	return
}

func (p *profileService) CheckPasswordVerificationCode(c *gin.Context, id uuid.UUID, request *externalDtos.ProfileCheckPasswordVerificationCodeRequest) (status bool, err error) {
	verificationCode, err := p.redisRepository.Get(c, externalDtos.NewPasswordVerificationCodeKey(id.String()))
	if err != nil {
		err = e.ErrVerificationCodeNotFound
		return
	}
	if verificationCode != request.VerificationCode {
		err = e.ErrInvalidVerificationCode
		return
	}

	status = true

	return
}

func (s *profileService) RequestChangePasswordVerificationCode(c *gin.Context, id uuid.UUID) (status bool, err error) {
	user, err := s.userRepository.GetByIDJoinsEmail(id)
	if err != nil {
		return
	}

	verificationCode := utils.RandomVerificationCode()
	err = s.redisRepository.Set(c, externalDtos.NewPasswordVerificationCodeKey(id.String()), verificationCode, 10*time.Minute)
	if err != nil {
		err = e.ErrCacheOperationFailed
		return
	}

	to := []string{user.Email.Address}
	subject := "Change Password Verification Code "
	message := utils.WriteChangePasswordEmailMessage(verificationCode)

	err = utils.SendMail(to, subject, message)
	if err != nil {
		err = e.ErrSendingEmailFailed
		return
	}

	status = true

	return
}
