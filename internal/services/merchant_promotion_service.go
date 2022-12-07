package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
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
	err := utils.ValidateMpCreateDatetime(mpRequest.StartDatetime, mpRequest.EndDatetime)
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

	// get all PMPs related to the promotion
	merchantPromotionModels, err := mp.productMerchantPromotionRepository.GetByMerchantPromotionID(merchantPromotionID)
	if err != nil {
		return nil, err
	}
	for _, merchantPromotionModel := range merchantPromotionModels {
		pmpResponse := dtos.NewProductMerchantPromotionResponse(merchantPromotionModel)
		response.ProductMerchantPromotion = append(response.ProductMerchantPromotion, pmpResponse)
	}

	return response, nil
}

func (mp *merchantPromotionService) Update(userID uuid.UUID, merchantPromotionID uint, mpRequest *dtos.MerchantPromotionRequest) (*dtos.MerchantPromotionResponse, error) {
	merchant, err := mp.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	currentMp, err := mp.merchantPromotionRepository.GetByID(merchantPromotionID)
	if err != nil {
		return nil, err
	}
	currentMpStatus := utils.PromotionOrVoucherStatus(currentMp.StartDatetime, currentMp.EndDatetime)
	// MP validations
	// name
	if mpRequest.Name == "" {
		return nil, e.ErrEmptyName
	}
	// datetime
	err = utils.ValidateMpUpdateDatetime(mpRequest.StartDatetime, mpRequest.EndDatetime, currentMp.StartDatetime, currentMp.EndDatetime, currentMpStatus)
	if err != nil {
		return nil, err
	}
	// match product ID list with PMP
	// convert ProductIDList []int64 -> []int
	var mpRequestProductIDList []int
	for _, productID := range mpRequest.ProductIDList {
		mpRequestProductIDList = append(mpRequestProductIDList, int(productID))
	}
	// convert PMP request to product ID list
	var pmpRequestProductIDList []int
	for _, pmp := range mpRequest.ProductMerchantPromotion {
		pmpRequestProductIDList = append(pmpRequestProductIDList, int(pmp.ProductID))
	}
	err = utils.ValidateMpProductIDList(mpRequestProductIDList, pmpRequestProductIDList)
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
	// currently existing PMPs
	currentPmpList, _ := mp.productMerchantPromotionRepository.GetByMerchantPromotionID(merchantPromotionID)
	var currentPmpProductIDList []int
	for _, pmp := range currentPmpList {
		currentPmpProductIDList = append(currentPmpProductIDList, int(pmp.ProductID))
	}

	for _, pmpRequest := range mpRequest.ProductMerchantPromotion {
		// PMP validation
		// ongoing promotion, no updates allowed to PMP
		if currentMpStatus == constants.OngoingPromotionOrVoucher {
			return nil, e.ErrPromotionIsOngoing
		}
		// validate discount and type
		err := utils.ValidatePmpDiscountAndType(pmpRequest.DiscountType, pmpRequest.Discount)
		if err != nil {
			return nil, err
		}

		if pmpRequest.DiscountType == constants.Nominal {
			pmpRequest.MaximumDiscountAmount = pmpRequest.Discount
		}

		// search PMP based on productIDList
		existsInCurrentPmp := utils.ContainsInt(currentPmpProductIDList, int(pmpRequest.ProductID))
		var updatedPmpModel *models.ProductMerchantPromotion
		if !existsInCurrentPmp {
			// CREATE
			var newPmpModelSlice []*models.ProductMerchantPromotion
			newPmpModel := pmpRequest.ToModel(updatedMpModel.ID)
			newPmpModelSlice = append(newPmpModelSlice, newPmpModel)
			createdPmpModel, err := mp.productMerchantPromotionRepository.CreateProductMerchantPromotion(newPmpModelSlice)
			if err != nil {
				return nil, err
			}
			updatedPmpModel = createdPmpModel[0]
		} else {
			// UPDATE
			// validate quota before update
			currentPmpModel, err := mp.productMerchantPromotionRepository.GetByMpAndProductID(merchantPromotionID, pmpRequest.ProductID)
			if err != nil {
				return nil, err
			}
			err = utils.ValidatePmpQuota(currentPmpModel.CurrentQuota, pmpRequest.MaximumQuota, pmpRequest.MaximumUserLimit)
			if err != nil {
				return nil, err
			}
			currentPmpModel.DiscountType = pmpRequest.DiscountType
			currentPmpModel.Discount = pmpRequest.Discount
			currentPmpModel.MaximumDiscountAmount = pmpRequest.MaximumDiscountAmount
			currentPmpModel.MaximumQuota = pmpRequest.MaximumQuota
			currentPmpModel.MaximumUserLimit = pmpRequest.MaximumUserLimit
			updatedPmpModel, err = mp.productMerchantPromotionRepository.Update(currentPmpModel)
			if err != nil {
				return nil, err
			}
		}
		updatedPmpResponse = append(updatedPmpResponse, dtos.NewProductMerchantPromotionResponse(updatedPmpModel))
	}

	// DELETE
	for _, productID := range currentPmpProductIDList {
		existsInPmpRequest := utils.ContainsInt(mpRequestProductIDList, productID)
		if !existsInPmpRequest {
			_, err := mp.productMerchantPromotionRepository.DeleteByProductID(uint(productID))
			if err != nil {
				return nil, err
			}
		}
	}

	response := dtos.NewMerchantPromotionResponse(mpModel, updatedPmpResponse)
	return response, nil
}
