package services

import (
	"sort"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type MerchantPromotionService interface {
	CreateMyShopPromotion(id uuid.UUID, mpRequest *dtos.MerchantPromotionRequest) ([]*dtos.ProductMerchantPromotionResponse, error)
	GetAll(userID uuid.UUID, query *dtos.MerchantPromotionQuery) (models []*models.MerchantPromotion, err error)
	GetByID(merchantPromotionID uint) (*dtos.MerchantPromotionDetails, error)
	Update(userID uuid.UUID, merchantPromotionID uint, mpRequest *dtos.MerchantPromotionRequest) (*dtos.MerchantPromotionResponse, error)
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

func (mp *merchantPromotionService) CreateMyShopPromotion(id uuid.UUID, mpRequest *dtos.MerchantPromotionRequest) ([]*dtos.ProductMerchantPromotionResponse, error) {
	err := utils.ValidateMerchantPromotionDatetime(mpRequest.StartDatetime, mpRequest.EndDatetime)
	if err != nil {
		return nil, err
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
	var responses []*dtos.ProductMerchantPromotionResponse

	for _, pmpRequest := range mpRequest.ProductMerchantPromotion {
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
		responses = append(responses, dtos.NewProductMerchantPromotionResponse(newPmpModel))
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

func (mp *merchantPromotionService) Update(userID uuid.UUID, merchantPromotionID uint, mpRequest *dtos.MerchantPromotionRequest) (*dtos.MerchantPromotionResponse, error) {
	err := utils.ValidateMerchantPromotionDatetime(mpRequest.StartDatetime, mpRequest.EndDatetime)
	if err != nil {
		return nil, err
	}

	merchant, err := mp.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	mpModel := mpRequest.ToModel(merchant.ID)
	mpModel.ID = merchantPromotionID

	updatedMpModel, err := mp.merchantPromotionRepository.Update(mpModel)
	if err != nil {
		return nil, err
	}

	var updatedPmpResponse []*dtos.ProductMerchantPromotionResponse
	// currenttly existing PMPs
	currentPmp, err := mp.productMerchantPromotionRepository.GetByMerchantPromotionID(merchantPromotionID)
	if err != nil {
		return nil, err
	}
	var currentPmpProductIDList []int
	for _, pmp := range currentPmp {
		currentPmpProductIDList = append(currentPmpProductIDList, int(pmp.ProductID))
	}

	// create/update PMP
	for _, pmpRequest := range mpRequest.ProductMerchantPromotion {
		if pmpRequest.DiscountType == "nominal" {
			pmpRequest.MaximumDiscountAmount = pmpRequest.Discount
		}

		if pmpRequest.DiscountType == "percentage" {
			if pmpRequest.Discount > 100 {
				return nil, e.ErrExceedDiscountPercentage
			}
		}

		// search PMP based on productIDList
		existsInCurrentPmp := sort.SearchInts(currentPmpProductIDList, int(pmpRequest.ProductID))
		pmpModel := pmpRequest.ToModel(updatedMpModel.ID)
		var updatedPmpModel *models.ProductMerchantPromotion
		switch existsInCurrentPmp {
		case len(currentPmpProductIDList):
			// create
			var pmpModelSlice []*models.ProductMerchantPromotion
			pmpModelSlice = append(pmpModelSlice, pmpModel)
			createdPmpModel, err := mp.productMerchantPromotionRepository.CreateProductMerchantPromotion(pmpModelSlice)
			if err != nil {
				return nil, err
			}
			updatedPmpModel = createdPmpModel[0]
		default:
			// update
			updatedPmpModel, err = mp.productMerchantPromotionRepository.Update(pmpModel)
			if err != nil {
				return nil, err
			}
		}
		updatedPmpResponse = append(updatedPmpResponse, dtos.NewProductMerchantPromotionResponse(updatedPmpModel))
	}

	// convert ProductIDList []int64 -> []int and sort ASC for SearchInts
	var productIDListIntSlice []int
	for _, productID := range mpRequest.ProductIDList {
		productIDListIntSlice = append(productIDListIntSlice, int(productID))
	}
	sort.SliceStable(productIDListIntSlice, func(i, j int) bool {
		return productIDListIntSlice[i] < productIDListIntSlice[j]
	})
	// delete
	for _, productID := range currentPmpProductIDList {
		existsInPmpRequest := sort.SearchInts(productIDListIntSlice, productID)
		switch existsInPmpRequest {
		case len(productIDListIntSlice):
			_, err := mp.productMerchantPromotionRepository.DeleteByProductID(uint(productID))
			if err != nil {
				return nil, err
			}
		}
	}

	response := dtos.NewMerchantPromotionResponse(mpModel, updatedPmpResponse)
	return response, nil
}
