package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	merchantVoucherIDUint64, err := strconv.ParseUint(c.Param("merchantVoucherId"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidParams)
		return
	}

	merchantVoucherID := uint(merchantVoucherIDUint64)

	merchantVouherUpdateResponse, err := h.merchantVoucherService.Update(id, merchantVoucherID, request)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVouherUpdateResponse, nil)
}

func (h *Handler) MerchantVoucherDelete(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	merchantVoucherIDUint64, err := strconv.ParseUint(c.Param("merchantVoucherId"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidParams)
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
	merchantVouherGetAllResponse, err := h.merchantVoucherService.GetAll(query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, merchantVouherGetAllResponse, nil)
}

func (h *Handler) MerchantVoucherGetByID(c *gin.Context) {
	merchantVoucherIDUint64, err := strconv.ParseUint(c.Param("merchantVoucherId"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidParams)
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
