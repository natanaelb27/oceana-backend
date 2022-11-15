package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type SealabsPayService interface {
	Create(sealabsPay *models.SealabsPay) (*models.SealabsPay, error)
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

func (p *sealabsPayService) Create(model *models.SealabsPay) (*models.SealabsPay, error) {
	if !utils.ValidateSealabsPayAccountNumber(model.AccountNumber) {
		return nil, e.ErrInvalidSealabsPayAccountNumber
	}

	return p.sealabsPayRepository.Create(model)
}
