package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) SealabsPayCreate(c *gin.Context) {
	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.SealabsPayCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	_, err := h.sealabsPayService.Create(c, userID, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, nil, nil)
}

func (h *Handler) SealabsPayGetAll(c *gin.Context) {
	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	results, err := h.sealabsPayService.GetAll(userID)
	if err != nil && err != e.ErrDataNotFound {
		utils.ResponseError(c, err)
		return
	}

	var responseData []*dtos.SealabsPayCreateResponse = make([]*dtos.SealabsPayCreateResponse, 0)

	for _, result := range results {
		responseData = append(responseData, dtos.NewSealabsPayCreateResponse(result))
	}

	utils.ResponseSuccess(c, responseData, nil)

}

func (h *Handler) SealabsPayGetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("sealabs-pay-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.sealabsPayService.GetById(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewSealabsPayCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) SealabsPayDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("sealabs-pay-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.sealabsPayService.Delete(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewSealabsPayDeleteResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) SealabsPaySetDefault(c *gin.Context) {
	userID := utils.GinGetUUID(c, "id")
	if userID == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	id, err := uuid.Parse(c.Param("sealabs-pay-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	err = h.sealabsPayService.SetDefault(id, userID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, nil, nil)
}

func (h *Handler) SealabsPayTopupOTP(c *gin.Context) {
	var inboundRequest *dtos.SealabsPayPaymentInboundRequest
	err := c.ShouldBindJSON(&inboundRequest)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	otpURL, err := h.sealabsPayService.TopupWalletOTP(inboundRequest)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, *otpURL, nil)
}

func (h *Handler) SealabsPayRegisterOTP(c *gin.Context) {
	var inboundRequest *dtos.SealabsPayPaymentInboundRequest
	err := c.ShouldBindJSON(&inboundRequest)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	otpURL, err := h.sealabsPayService.RegisterOTP(inboundRequest)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, *otpURL, nil)
}

func (h *Handler) SealabsPayCallback(c *gin.Context) {
	var request *dtos.SealabsPayCallbackRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	err = h.sealabsPayService.StoreSignature(c, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, request, nil)
}
