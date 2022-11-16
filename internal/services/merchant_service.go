package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type MerchantService interface {
	Register(merchantCreateRequest *dtos.MerchantCreateRequest) (merchantCreateResponse *dtos.MerchantCreateResponse, err error)
}

type merchantService struct {
	merchantRepository r.MerchantRepository
}

type MerchantServiceConfig struct {
	MerchantRepository r.MerchantRepository
}

func NewMerchantService(c *MerchantServiceConfig) MerchantService {
	return &merchantService{
		merchantRepository: c.MerchantRepository,
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
		return nil, err
	}
	merchantCreateResponse = dtos.NewMerchantCreateResponse(createdMerchant)
	return
}
