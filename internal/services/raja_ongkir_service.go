package services

import (
	"fmt"

	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type RajaOngkirService interface {
	GetRajaOngkirShippingCost(request *dtos.GetRajaOngkirShippingCostRequest) (*dtos.GetRajaOngkirShippingCostResponse, error)
}

type rajaOngkirService struct {
	rajaOngkirRepository            r.RajaOngkirRepository
	productCourierServiceRepository r.ProductCourierServiceRepository
}

type RajaOngkirServiceConfig struct {
	RajaOngkirRepository            r.RajaOngkirRepository
	ProductCourierServiceRepository r.ProductCourierServiceRepository
}

func NewRajaOngkirService(c *RajaOngkirServiceConfig) RajaOngkirService {
	return &rajaOngkirService{
		rajaOngkirRepository:            c.RajaOngkirRepository,
		productCourierServiceRepository: c.ProductCourierServiceRepository,
	}
}

func (s *rajaOngkirService) GetRajaOngkirShippingCost(request *dtos.GetRajaOngkirShippingCostRequest) (*dtos.GetRajaOngkirShippingCostResponse, error) {
	productCourierServices, err := s.productCourierServiceRepository.GetAllCourierAndCourierServiceFromProductIDs(request.ProductIDs)
	if err != nil {
		return nil, err
	}

	// find courier in product courier services
	serviceCodes := map[string]bool{}
	for _, vPCS := range productCourierServices {
		if vPCS.Code == request.Courier {
			for _, vSC := range vPCS.ServiceCodes {
				serviceCodes[vPCS.Code+vSC] = true
			}
		}
	}
	if len(serviceCodes) == 0 {
		return nil, e.ErrCourierServicesNotFound
	}

	result, err := s.rajaOngkirRepository.GetShippingCost(request)
	if err != nil {
		return nil, err
	}

	rajaOngkirShipping := dtos.NewGetRajaOngkirShippingCostResponse(result)
	newRajaOngkirShipping := &dtos.GetRajaOngkirShippingCostResponse{
		Code:  rajaOngkirShipping.Code,
		Name:  rajaOngkirShipping.Name,
		Costs: make([]dtos.RajaOngkirShippingService, 0),
	}

	// check serviceCodes[ROS.code+service] exist then append to a new raja ongkir shipping for user to see
	for _, vROSC := range rajaOngkirShipping.Costs {
		fmt.Println(rajaOngkirShipping.Code + vROSC.Service)
		if exists := serviceCodes[rajaOngkirShipping.Code+vROSC.Service]; !exists {
			continue
		}
		newRajaOngkirShipping.Costs = append(newRajaOngkirShipping.Costs, vROSC)
	}
	if len(newRajaOngkirShipping.Costs) == 0 {
		return nil, e.ErrCourierServicesNotFound
	}

	return newRajaOngkirShipping, nil
}
