package services

import (
	internalLogisticDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/logistic_dtos"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"github.com/google/uuid"
)

type LogisticService interface {
	GetByMerchantId(id uuid.UUID) (*internalLogisticDtos.Logistics, error)
}

type logisticService struct {
	courierRepository                r.CourierRepository
	courierServiceRepository         r.CourierServiceRepository
	merchantCourierServiceRepository r.MerchantCourierServiceRepository
	userRepository                   r.UserRepository
	merchantRepository               r.MerchantRepository
}

type LogisticServiceConfig struct {
	CourierRepository                r.CourierRepository
	CourierServiceRepository         r.CourierServiceRepository
	MerchantCourierServiceRepository r.MerchantCourierServiceRepository
	UserRepository                   r.UserRepository
	MerchantRepository               r.MerchantRepository
}

func NewLogisticService(c *LogisticServiceConfig) LogisticService {
	return &logisticService{
		courierRepository:                c.CourierRepository,
		courierServiceRepository:         c.CourierServiceRepository,
		merchantCourierServiceRepository: c.MerchantCourierServiceRepository,
		userRepository:                   c.UserRepository,
		merchantRepository:               c.MerchantRepository,
	}
}

func (p *logisticService) GetByMerchantId(id uuid.UUID) (*internalLogisticDtos.Logistics, error) {
	merchant, err := p.merchantRepository.GetByUserID(id)
	if err != nil {
		return nil, err
	}

	newCouriers := make([]*internalLogisticDtos.Courier, 0)

	couriers, err := p.courierRepository.GetAll()
	for _, courier := range couriers {
		courierServices, err := p.merchantCourierServiceRepository.GetAllByCourierID(courier.ID, merchant.ID)
		if err != nil {
			return nil, err
		}

		// set courier[x].is_selected depending on if any of its courier_services is selected
		courierIsSelected := false
		for _, courierService := range courierServices {
			if courierService.IsSelected {
				courierIsSelected = true
			}
		}

		newCouriers = append(newCouriers, &internalLogisticDtos.Courier{
			ID:             courier.ID,
			Code:           courier.Code,
			Name:           courier.Name,
			IsSelected:     courierIsSelected,
			CourierService: courierServices,
		})
	}

	logistics := &internalLogisticDtos.Logistics{
		// Get PreorderDay
		// Get IsPrintLabel
		Courier: newCouriers,
	}

	return logistics, err
}
