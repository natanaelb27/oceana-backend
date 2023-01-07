package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
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

	query := dtos.NewWalletDetailsQuery(c)

	responseData, pagination, err := h.walletService.GetWalletDetails(id, query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, responseData, pagination)
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

	tx := c.MustGet("tx").(*gorm.DB)

	err := h.walletService.WithTx(tx).CreateWalletTopup(c, id, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, nil, nil)
}
