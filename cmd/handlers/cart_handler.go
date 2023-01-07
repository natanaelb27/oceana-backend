package handlers

import (
	"github.com/gin-gonic/gin"
	in_cart_dtos "github.com/natanaelb27/oceana-backend/internal/dtos/cart_dtos"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	"gorm.io/gorm"
)

func (h *Handler) CartGetAllByUserId(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")

	results, err := h.cartService.GetAllByUserId(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	var responseData []*in_cart_dtos.Merchant = make([]*in_cart_dtos.Merchant, 0)

	responseData = append(responseData, results...)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) CartUpdate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")

	var request *dtos.CartUpsertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	internalRequest := request.CartUpsertRequestToInternalCartRequest()

	internalRequest.UserID = id

	result, err := h.cartService.Update(internalRequest)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewCartResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) PayCartWithWallet(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")

	var request *dtos.PayCartWithWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	request.UserID = id

	tx := c.MustGet("tx").(*gorm.DB)

	result, err := h.cartService.WithTx(tx).PayCartWithWallet(c, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, result, nil)
}

func (h *Handler) PayCartWithSealabsPay(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")

	var request *dtos.PayCartWithSealabsPayRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	request.UserID = id

	tx := c.MustGet("tx").(*gorm.DB)

	result, err := h.cartService.WithTx(tx).PayCartWithSealabsPay(c, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, result, nil)
}

func (h *Handler) PayCartWithSealabsPaySuccess(c *gin.Context) {
	var request *dtos.PayCartWithSealabsPaySuccessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	tx := c.MustGet("tx").(*gorm.DB)

	err := h.cartService.WithTx(tx).PayCartWithSealabsPaySuccess(c, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, nil, nil)
}