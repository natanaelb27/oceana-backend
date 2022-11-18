package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type MerchantService interface {
	Register(merchantCreateRequest *dtos.MerchantCreateRequest) (merchantCreateResponse *dtos.MerchantCreateResponse, err error)
	GetInformationByID(id uint) (merchantGetInformationByIDResponse *dtos.MerchantGetInformationByIDResponse, err error)
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

func (m *merchantService) GetInformationByID(id uint) (merchantGetInformationByIDResponse *dtos.MerchantGetInformationByIDResponse, err error) {
	merchantInformation, err := m.merchantRepository.GetInformationByID(id)
	if err != nil {
		return
	}

	merchantGetInformationByIDResponse = dtos.NewMerchantGetInformationByIDResponse(merchantInformation)
	return
}
