package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func (h *Handler) MerchantVoucherCreate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.MerchantVoucherCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	merchantVouherCreateResponse, err := h.merchantVoucherService.Create(id, request)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVouherCreateResponse, nil)
}

func (h *Handler) MerchantVoucherUpdate(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.MerchantVoucherUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	merchantVoucherIDUint64, err := strconv.ParseUint(c.Param("merchant-voucher-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantVoucherID := uint(merchantVoucherIDUint64)

	merchantVoucherUpdateResponse, err := h.merchantVoucherService.Update(id, merchantVoucherID, request)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVoucherUpdateResponse, nil)
}

func (h *Handler) MerchantVoucherDelete(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantVoucherIDUint64, err := strconv.ParseUint(c.Param("merchant-voucher-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantVoucherID := uint(merchantVoucherIDUint64)

	merchantVouherDeleteResponse, err := h.merchantVoucherService.Delete(id, merchantVoucherID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVouherDeleteResponse, nil)
}

func (h *Handler) MerchantVoucherGetAll(c *gin.Context) {
	query := dtos.NewMerchantVoucherQuery(c)
	merchantVouherGetAllResponse, pagination, err := h.merchantVoucherService.GetAll(query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVouherGetAllResponse, pagination)
}

func (h *Handler) MerchantVoucherGetByID(c *gin.Context) {
	merchantVoucherIDUint64, err := strconv.ParseUint(c.Param("merchant-voucher-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantVoucherID := uint(merchantVoucherIDUint64)

	merchantVouherGetByIDResponse, err := h.merchantVoucherService.GetByID(merchantVoucherID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVouherGetByIDResponse, nil)
}

func (h *Handler) MerchantVoucherGetAllByBuyer(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewBuyerMerchantVoucherQuery(c)
	merchantVouherGetAllResponse, pagination, err := h.merchantVoucherService.GetAllByBuyer(id, query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVouherGetAllResponse, pagination)
}
