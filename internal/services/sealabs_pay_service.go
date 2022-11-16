package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type SealabsPayService interface {
	Create(id uuid.UUID, request *dtos.SealabsPayCreateRequest) (*dtos.SealabsPayCreateResponse, error)
}

type sealabsPayService struct {
	sealabsPayRepository r.SealabsPayRepository
}

type SealabsPayServiceConfig struct {
	SealabsPayRepository r.SealabsPayRepository
}

func NewSealabsPayService(c *SealabsPayServiceConfig) SealabsPayService {
	return &sealabsPayService{
		sealabsPayRepository: c.SealabsPayRepository,
	}
}

func (p *sealabsPayService) Create(id uuid.UUID, request *dtos.SealabsPayCreateRequest) (*dtos.SealabsPayCreateResponse, error) {
	model := request.ToModel(id)

	if !utils.ValidateSealabsPayAccountNumber(model.AccountNumber) {
		return nil, e.ErrInvalidSealabsPayAccountNumber
	}

	newModel, err := p.sealabsPayRepository.Create(model)

	if err != nil {
		return nil, err
	}

	response := dtos.NewSealabsPayCreateResponse(newModel)

	return response, nil
}
