package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) WalletActivate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.WalletActivateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	requestModel := request.ToModel()

	result, err := h.walletService.Activate(id, requestModel)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewWalletActivateResponse(result)
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) WalletChangePin(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.WalletChangePinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	result, err := h.walletService.ChangePin(c, id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewWalletChangePinResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) WalletGetDetails(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	responseData, err := h.walletService.GetWalletDetails(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) WalletTopup(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.TopupWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	err := h.walletService.CreateWalletTopup(c, id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, nil, nil)
}
