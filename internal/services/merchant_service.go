package services

import (
	"github.com/google/uuid"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type MerchantService interface {
	WithTx(tx *gorm.DB) *merchantService
	Register(merchantCreateRequest *dtos.MerchantCreateRequest) (merchantCreateResponse *dtos.MerchantCreateResponse, err error)
	GetInformationByID(id uint) (merchantGetInformationByIDResponse *dtos.MerchantGetInformationByIDResponse, err error)
	MerchantGetProfile(id uuid.UUID) (merchantGetProfileResponse *dtos.MerchantGetProfileResponse, err error)
	MerchantUpdatePhoto(request *dtos.MerchantUpdatePhotoRequest) (merchantGetProfileResponse *dtos.MerchantGetProfileResponse, err error)
}

type merchantService struct {
	merchantRepository r.MerchantRepository
	userRepository     r.UserRepository
}

type MerchantServiceConfig struct {
	MerchantRepository r.MerchantRepository
	UserRepository     r.UserRepository
}

func NewMerchantService(c *MerchantServiceConfig) MerchantService {
	return &merchantService{
		merchantRepository: c.MerchantRepository,
		userRepository:     c.UserRepository,
	}
}

func (u merchantService) WithTx(tx *gorm.DB) *merchantService {
	u.merchantRepository = u.merchantRepository.WithTx(tx)
	u.userRepository = u.userRepository.WithTx(tx)

	return &u
}

func (m *merchantService) Register(merchantCreateRequest *dtos.MerchantCreateRequest) (merchantCreateResponse *dtos.MerchantCreateResponse, err error) {
	_, err = m.merchantRepository.GetByName(merchantCreateRequest.Name)
	if err != e.ErrDataNotFound {
		return nil, e.ErrMerchantNameAlreadyExist
	}
	merchant := merchantCreateRequest.ToModel()
	createdMerchant, err := m.merchantRepository.Create(merchant)
	if err != nil {
		return
	}
	user, err := m.userRepository.GetByID(merchantCreateRequest.UserID)
	if err != nil {
		return nil, err
	}
	user.Role = utils.ReplaceAtIndex(user.Role, '1', 1)
	_, err = m.userRepository.Update(merchantCreateRequest.UserID, user)
	if err != nil {
		return
	}
	merchantCreateResponse = dtos.NewMerchantCreateResponse(createdMerchant)
	return
}

func (s *merchantService) GetInformationByID(id uint) (merchantGetInformationByIDResponse *dtos.MerchantGetInformationByIDResponse, err error) {
	merchantInformation, err := s.merchantRepository.GetInformationByID(id)
	if err != nil {
		return
	}

	merchantGetInformationByIDResponse = dtos.NewMerchantGetInformationByIDResponse(merchantInformation)
	return
}

func (s *merchantService) MerchantGetProfile(id uuid.UUID) (merchantGetProfileResponse *dtos.MerchantGetProfileResponse, err error) {
	merchant, err := s.merchantRepository.GetByUserID(id)
	if err != nil {
		return
	}

	merchantGetProfileResponse = dtos.NewMerchantGetProfileResponse(merchant)
	return
}

func (s *merchantService) MerchantUpdatePhoto(request *dtos.MerchantUpdatePhotoRequest) (merchantGetProfileResponse *dtos.MerchantGetProfileResponse, err error) {
	if request.Photo == nil {
		return nil, e.ErrEmptyImage
	}

	_, err = s.merchantRepository.GetByUserID(request.ID)

	if err != nil {
		return
	}

	photoFile, err := request.Photo.Open()

	if err != nil {
		return nil, err
	}

	defer photoFile.Close()

	imageUrl, err := utils.UploadMedia(photoFile)

	if err != nil {
		return nil, err
	}

	merchant, err := s.merchantRepository.UpdateProfilePhoto(request.ID, imageUrl)

	merchantGetProfileResponse = dtos.NewMerchantGetProfileResponse(merchant)
	return
}
