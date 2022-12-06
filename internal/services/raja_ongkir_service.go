package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
)

type RajaOngkirService interface {
	GetRajaOngkirShippingCost(request *dtos.GetRajaOngkirShippingCostRequest) (*dtos.GetRajaOngkirShippingCostResponse, error)
}

type rajaOngkirService struct {
	rajaOngkirRepository r.RajaOngkirRepository
}

type RajaOngkirServiceConfig struct {
	RajaOngkirRepository r.RajaOngkirRepository
}

func NewRajaOngkirService(c *RajaOngkirServiceConfig) RajaOngkirService {
	return &rajaOngkirService{
		rajaOngkirRepository: c.RajaOngkirRepository,
	}
}

func (s *rajaOngkirService) GetRajaOngkirShippingCost(request *dtos.GetRajaOngkirShippingCostRequest) (*dtos.GetRajaOngkirShippingCostResponse, error) {
	result, err := s.rajaOngkirRepository.GetShippingCost(request)

	if err != nil {
		return nil, err
	}

	return dtos.NewGetRajaOngkirShippingCostResponse(result), nil
}
