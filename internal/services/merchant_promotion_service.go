package services

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type MerchantPromotionService interface {
	CreateMyShopPromotion(id uuid.UUID, mpRequest *dtos.MerchantPromotionCreateRequest, pmpRequests []*dtos.ProductMerchantPromotionCreateRequest) ([]*dtos.ProductMerchantPromotionCreateResponse, error)
	GetAll(userID uuid.UUID, query *dtos.MerchantPromotionQuery) (models []*models.MerchantPromotion, err error)
	GetByID(merchantPromotionID uint) (*dtos.MerchantPromotionDetails, error)
	Update(userID uuid.UUID, merchantPromotionID uint) (*dtos.MerchantPromotionUpdateResponse, error)
}

type merchantPromotionService struct {
	merchantRepository                 r.MerchantRepository
	merchantPromotionRepository        r.MerchantPromotionRepository
	productMerchantPromotionRepository r.ProductMerchantPromotionRepository
}

type MerchantPromotionServiceConfig struct {
	MerchantRepository                 r.MerchantRepository
	MerchantPromotionRepository        r.MerchantPromotionRepository
	ProductMerchantPromotionRepository r.ProductMerchantPromotionRepository
}

func NewMerchantPromotionService(c *MerchantPromotionServiceConfig) MerchantPromotionService {
	return &merchantPromotionService{
		merchantRepository:                 c.MerchantRepository,
		merchantPromotionRepository:        c.MerchantPromotionRepository,
		productMerchantPromotionRepository: c.ProductMerchantPromotionRepository,
	}
}

func (mp *merchantPromotionService) CreateMyShopPromotion(id uuid.UUID, mpRequest *dtos.MerchantPromotionCreateRequest, pmpRequests []*dtos.ProductMerchantPromotionCreateRequest) ([]*dtos.ProductMerchantPromotionCreateResponse, error) {
	if mpRequest.StartDatetime.Before(time.Now()) {
		return nil, e.ErrPromotionStartDatetimeEarlierThanNow
	}

	if mpRequest.EndDatetime.Before(mpRequest.StartDatetime) {
		return nil, e.ErrPromotionEndDatetimeEarlierThanVoucherStartDatetime
	}

	if int(time.Until(mpRequest.StartDatetime).Hours()/24) <= 0 {
		return nil, e.ErrPromotionStartDatetimeMinimumOneDayFromNow
	}

	if int(mpRequest.EndDatetime.Sub(mpRequest.StartDatetime).Hours()) <= 0 {
		return nil, e.ErrPromotionPeriodMinimumOneHour
	}

	if int(mpRequest.EndDatetime.Sub(mpRequest.StartDatetime).Hours()/24) > 180 {
		return nil, e.ErrPromotionPeriodExceedMaximumDays
	}

	merchant, err := mp.merchantRepository.GetByUserID(id)
	if err != nil {
		return nil, err
	}

	mpModel := mpRequest.ToModel(merchant.ID)

	// transaction needed

	newMpModel, err := mp.merchantPromotionRepository.CreateMyShopPromotion(mpModel, mpRequest.ProductIDList)
	if err != nil {
		return nil, err
	}

	var newPmpModels []*models.ProductMerchantPromotion
	var responses []*dtos.ProductMerchantPromotionCreateResponse

	for _, pmpRequest := range pmpRequests {
		if pmpRequest.DiscountType == "nominal" {
			pmpRequest.MaximumDiscountAmount = pmpRequest.Discount
		}

		if pmpRequest.DiscountType == "percentage" {
			if pmpRequest.Discount > 100 {
				return nil, e.ErrExceedDiscountPercentage
			}
		}

		newPmpModels = append(newPmpModels, pmpRequest.ToModel(newMpModel.ID))

	}

	newPmpModels, err = mp.productMerchantPromotionRepository.CreateProductMerchantPromotion(newPmpModels)
	if err != nil {
		return nil, err
	}

	for _, newPmpModel := range newPmpModels {
		responses = append(responses, dtos.NewProductMerchantPromotionCreateResponse(newPmpModel))
	}

	return responses, nil
}

func (mp *merchantPromotionService) GetAll(userID uuid.UUID, query *dtos.MerchantPromotionQuery) (models []*models.MerchantPromotion, err error) {
	merchant, err := mp.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	result, err := mp.merchantPromotionRepository.GetAll(merchant.ID, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (mp *merchantPromotionService) GetByID(merchantPromotionID uint) (*dtos.MerchantPromotionDetails, error) {
	result, err := mp.merchantPromotionRepository.GetByID(merchantPromotionID)
	if err != nil {
		return nil, err
	}

	response := dtos.NewMerchantPromotionDetails(result)

	return response, nil
}

func (mp *merchantPromotionService) Update(userID uuid.UUID, merchantPromotionID uint) (*dtos.MerchantPromotionUpdateResponse, error) {
	return nil, nil
}
