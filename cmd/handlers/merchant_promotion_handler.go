package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

func (h *Handler) MerchantPromotionCreate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var mpRequest *dtos.MerchantPromotionRequest
	if err := c.ShouldBindJSON(&mpRequest); err != nil {
		utils.ResponseError(c, err)
		return
	}

	tx := c.MustGet("tx").(*gorm.DB)

	responses, err := h.merchantPromotionService.WithTx(tx).CreateMyShopPromotion(id, mpRequest)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, responses, nil)
}

func (h *Handler) MerchantPromotionGetAll(c *gin.Context) {
	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewMerchantPromotionQuery(c)

	results, pagination, err := h.merchantPromotionService.GetAll(userID, query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	// var responseData []*dtos.MerchantPromotionResponse = make([]*dtos.MerchantPromotionResponse, 0)

	// for _, result := range results {
	// 	responseData = append(responseData, dtos.NewMerchantPromotionResponse(result, nil))
	// }

	utils.ResponseSuccess(c, results, pagination)
}

func (h *Handler) MerchantPromotionGetByID(c *gin.Context) {
	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantPromotionsIDUint64, err := strconv.ParseUint(c.Param("merchant-promotion-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantPromotionsID := uint(merchantPromotionsIDUint64)

	responseData, err := h.merchantPromotionService.GetDetailsByID(merchantPromotionsID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) MerchantPromotionUpdate(c *gin.Context) {
	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.MerchantPromotionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	merchantPromotionIDUint64, err := strconv.ParseUint(c.Param("merchant-promotion-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	merchantPromotionID := uint(merchantPromotionIDUint64)

	tx := c.MustGet("tx").(*gorm.DB)

	responses, err := h.merchantPromotionService.WithTx(tx).Update(userID, merchantPromotionID, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, responses, nil)
}

func (h *Handler) MerchantPromotionDelete(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantPromotionIDUint64, err := strconv.ParseUint(c.Param("merchant-promotion-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantPromotionID := uint(merchantPromotionIDUint64)

	merchantPromotionDeleteResponse, err := h.merchantPromotionService.Delete(id, merchantPromotionID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantPromotionDeleteResponse, nil)
}
