package services

import (
	internalDtos "github.com/natanaelb27/oceana-backend/internal/dtos"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
)

type CourierService interface {
	GetAll() ([]*internalDtos.Courier, error)
}

type courierService struct {
	courierRepository        r.CourierRepository
	courierServiceRepository r.CourierServiceRepository
}

type CourierServiceConfig struct {
	CourierRepository        r.CourierRepository
	CourierServiceRepository r.CourierServiceRepository
}

func NewCourierService(c *CourierServiceConfig) CourierService {
	return &courierService{
		courierRepository:        c.CourierRepository,
		courierServiceRepository: c.CourierServiceRepository,
	}
}
func (p *courierService) GetAll() ([]*internalDtos.Courier, error) {
	newCouriers := make([]*internalDtos.Courier, 0)

	couriers, err := p.courierRepository.GetAll()
	for _, courier := range couriers {
		courierServices, err := p.courierServiceRepository.GetAllByCourierID(courier.ID)
		if err != nil {
			return nil, err
		}

		newCouriers = append(newCouriers, &internalDtos.Courier{
			ID:             courier.ID,
			Code:           courier.Code,
			Name:           courier.Name,
			CourierService: courierServices,
		})
	}

	return newCouriers, err
}
