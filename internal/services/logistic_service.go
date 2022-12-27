package services

import (
	"fmt"

	internalLogisticDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/logistic_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type LogisticService interface {
	GetByMerchantId(id uuid.UUID) (*internalLogisticDtos.Logistic, error)
	Update(id uuid.UUID, request *internalLogisticDtos.Logistic) (*internalLogisticDtos.Logistic, error)
	GetByMerchantIdSelected(userID uuid.UUID) (*internalLogisticDtos.Logistic, error)
	ProductGetByMerchantID(merchantID uint) (*internalLogisticDtos.Logistic, error)
}

type logisticService struct {
	courierRepository                r.CourierRepository
	courierServiceRepository         r.CourierServiceRepository
	merchantCourierServiceRepository r.MerchantCourierServiceRepository
	userRepository                   r.UserRepository
	merchantRepository               r.MerchantRepository
	merchantConfigurationRepository  r.MerchantConfigurationRepository
	productRepository                r.ProductRepository
}

type LogisticServiceConfig struct {
	CourierRepository                r.CourierRepository
	CourierServiceRepository         r.CourierServiceRepository
	MerchantCourierServiceRepository r.MerchantCourierServiceRepository
	UserRepository                   r.UserRepository
	MerchantRepository               r.MerchantRepository
	MerchantConfigurationRepository  r.MerchantConfigurationRepository
	ProductRepository                r.ProductRepository
}

func NewLogisticService(c *LogisticServiceConfig) LogisticService {
	return &logisticService{
		courierRepository:                c.CourierRepository,
		courierServiceRepository:         c.CourierServiceRepository,
		merchantCourierServiceRepository: c.MerchantCourierServiceRepository,
		userRepository:                   c.UserRepository,
		merchantRepository:               c.MerchantRepository,
		merchantConfigurationRepository:  c.MerchantConfigurationRepository,
		productRepository:                c.ProductRepository,
	}
}

func (s *logisticService) GetByMerchantId(userID uuid.UUID) (*internalLogisticDtos.Logistic, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	newCouriers := make([]*internalLogisticDtos.Courier, 0)

	couriers, err := s.courierRepository.GetAll()
	if err != nil {
		return nil, err
	}
	for _, courier := range couriers {
		courierServices, err := s.merchantCourierServiceRepository.GetAllByCourierID(courier.ID, merchant.ID)
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
			LogoURL:        courier.LogoURL,
			IsSelected:     courierIsSelected,
			CourierService: courierServices,
		})
	}
	merchantConfiguration, err := s.merchantConfigurationRepository.GetByMerchantID(merchant.ID)
	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}

	logistics := &internalLogisticDtos.Logistic{
		PreorderDay:  merchantConfiguration.PreorderDay,
		IsPrintLabel: merchantConfiguration.IsPrintLabel,
		Courier:      newCouriers,
	}

	return logistics, nil
}

func (s *logisticService) GetByMerchantIdSelected(userID uuid.UUID) (*internalLogisticDtos.Logistic, error) {
	logistics, err := s.GetByMerchantId(userID)
	if err != nil {
		return nil, err
	}

	newCouriers := make([]*internalLogisticDtos.Courier, 0)

	for _, vCourier := range logistics.Courier {
		newCourierServices := make([]*internalLogisticDtos.CourierService, 0)
		for _, v := range vCourier.CourierService {
			if v.IsSelected {
				newCourierServices = append(newCourierServices, v)
			}
		}
		vCourier.CourierService = newCourierServices
		newCouriers = append(newCouriers, vCourier)
	}

	newLogistics := &internalLogisticDtos.Logistic{
		PreorderDay:  logistics.PreorderDay,
		IsPrintLabel: logistics.IsPrintLabel,
		Courier:      newCouriers,
	}

	return newLogistics, nil
}

func (s *logisticService) ProductGetByMerchantID(merchantID uint) (*internalLogisticDtos.Logistic, error) {
	merchantUser, err := s.merchantRepository.GetByMerchantIDJoinUser(merchantID)
	fmt.Println(merchantUser.ID)
	if err != nil {
		return nil, err
	}
	return s.GetByMerchantIdSelected(merchantUser.UserID)
}

func (s *logisticService) Update(userID uuid.UUID, requestLogistics *internalLogisticDtos.Logistic) (*internalLogisticDtos.Logistic, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	logistics, err := s.GetByMerchantId(userID)
	if err != nil {
		return nil, err
	}

	// validate model by courier length
	if len(logistics.Courier) != len(requestLogistics.Courier) {
		return nil, e.ErrInvalidModel
	}

	// index courier
	var iC int
	for iC = range logistics.Courier {
		oriC := logistics.Courier[iC]        // original courier model
		reqC := requestLogistics.Courier[iC] // request courier model

		// validate model - courier id
		if oriC.ID != reqC.ID {
			return nil, e.ErrInvalidModel
		}

		// validate model by courier service length
		if len(oriC.CourierService) != len(reqC.CourierService) {
			return nil, e.ErrInvalidModel
		}

		// index courier service
		var iCS int
		for iCS = range oriC.CourierService {
			// alias
			oriCS := oriC.CourierService[iCS] // original courier service model
			reqCS := reqC.CourierService[iCS] // request courier service model

			// validate model - courier id
			if oriCS.ID != reqCS.ID {
				return nil, e.ErrInvalidModel
			}

			if oriCS.IsSelected == reqCS.IsSelected {
				// continue to next iteration
				continue
			}

			newMerchantCourierServiceMap := map[string]any{
				"merchant_id":        merchant.ID,
				"courier_service_id": oriCS.ID,
				"is_active":          reqCS.IsSelected,
			}
			// update or insert merchant courier service
			err = s.merchantCourierServiceRepository.UpsertUsingMap(newMerchantCourierServiceMap)
			if err != nil {
				return nil, err
			}
		}
	}

	// identify if there are any changes in merchant configuration
	merchantConfigurationChanges := false
	// initiate merchant configuration struct
	newMerchantConfiguration := &models.MerchantConfiguration{
		MerchantID:   merchant.ID,
		PreorderDay:  logistics.PreorderDay,
		IsPrintLabel: logistics.IsPrintLabel,
	}

	// if preorderDay changes
	if requestLogistics.PreorderDay != logistics.PreorderDay {
		merchantConfigurationChanges = true
		// fill merchan_configuration > preorder_day
		newMerchantConfiguration.PreorderDay = requestLogistics.PreorderDay
	}
	// if isPrintLabel changes
	if requestLogistics.IsPrintLabel != logistics.IsPrintLabel {
		merchantConfigurationChanges = true
		// fill merchant_configurations > IsPrintLabel
		newMerchantConfiguration.IsPrintLabel = requestLogistics.IsPrintLabel
	}

	// if there are any changes in merchant configuration
	if merchantConfigurationChanges {
		// get merchant_cnfiguration by merchant_id
		_, err = s.merchantConfigurationRepository.GetByMerchantID(merchant.ID)
		if err != nil && err != e.ErrDataNotFound {
			return nil, err
		}
		if err == e.ErrDataNotFound {
			// create merchant_configuration with its struct value
			_, err = s.merchantConfigurationRepository.Create(newMerchantConfiguration)
			if err != nil {
				return nil, err
			}
		} else {
			// update merchant_configuration with its struct
			newMerchantConfigurationMap := map[string]any{
				"preorder_day":   newMerchantConfiguration.PreorderDay,
				"is_print_label": newMerchantConfiguration.IsPrintLabel,
			}
			newMerchantConfigurationConditionMap := map[string]any{
				"merchant_id": merchant.ID,
			}
			err = s.merchantConfigurationRepository.UpdateByConditionMap(newMerchantConfigurationConditionMap, newMerchantConfigurationMap)
			if err != nil {
				return nil, err
			}
		}
	}

	// always update preorder day of all products by that merchant
	productMap := map[string]any{
		"preorder_day": requestLogistics.PreorderDay,
	}
	productMapCondition := map[string]any{
		"merchant_id": merchant.ID,
		"status":      "active",
	}
	err = s.productRepository.UpdateUsingMapByMapCondition(productMapCondition, productMap)
	if err != nil {
		return nil, err
	}

	// get latest model
	newLogistic, err := s.GetByMerchantId(userID)
	if err != nil {
		return nil, err
	}

	return newLogistic, err
}
