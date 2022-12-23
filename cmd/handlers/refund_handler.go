package handlers

import (
	"encoding/json"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateRefund(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	orderID, err := uuid.Parse(c.Param("order-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	var request *dtos.CreateRefundRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	orderDetailRefund := []*dtos.OrderDetailRefund{}

	err = json.Unmarshal([]byte(request.OrderDetailForm), &orderDetailRefund)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	results, err := h.refundService.Create(id, orderID, request)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}

func (h *Handler) GetAllRefundDetailsByBuyer(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewGetAllRefundQuery(c)

	results, pagination, err := h.refundService.GetAllDetailsByBuyer(id, query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, pagination)
}

func (h *Handler) GetDetailsByID(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	refundID, err := uuid.Parse(c.Param("refund-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	results, err := h.refundService.GetDetailsByID(id, refundID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}

func (h *Handler) GetAllRefundDetailsByMerchant(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewGetAllRefundQuery(c)

	results, pagination, err := h.refundService.GetAllDetailsByMerchant(id, query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, pagination)
}

func (h *Handler) GetAllRefundDetailsByAdmin(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewGetAllRefundQuery(c)

	results, pagination, err := h.refundService.GetAllDetailsByAdmin(id, query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, pagination)
}

func (h *Handler) UpdateApproveRefund(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	refundID, err := uuid.Parse(c.Param("refund-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.refundService.UpdateApproveRefund(refundID, id)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewUpdateRefundResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UpdateRejectRefund(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	refundID, err := uuid.Parse(c.Param("refund-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.refundService.UpdateRejectRefund(refundID, id)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewUpdateRefundResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) UpdateCompleteRefund(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	refundID, err := uuid.Parse(c.Param("refund-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.refundService.UpdateCompleteRefund(refundID, id)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewUpdateRefundResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
