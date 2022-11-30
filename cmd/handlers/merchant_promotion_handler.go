package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) MerchantPromotionCreate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var mpRequest *dtos.MerchantPromotionCreateRequest
	if err := c.ShouldBindJSON(&mpRequest); err != nil {
		utils.ResponseError(c, err)
		return
	}

	responses, err := h.merchantPromotionService.CreateMyShopPromotion(id, mpRequest, mpRequest.ProductMerchantPromotion)

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

	results, err := h.merchantPromotionService.GetAll(userID, query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	var responseData []*dtos.MerchantPromotionCreateResponse = make([]*dtos.MerchantPromotionCreateResponse, 0)

	for _, result := range results {
		responseData = append(responseData, dtos.NewMerchantPromotionCreateResponse(result))
	}

	utils.ResponseSuccess(c, responseData, nil)
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

	responseData, err := h.merchantPromotionService.GetByID(merchantPromotionsID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, responseData, nil)
}
