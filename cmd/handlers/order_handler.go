package handlers

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) OrderListGetAll(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewOrderListQuery(c)
	query.FilterMerchantID = 0

	orderListGetAllResponse, pagination, err := h.orderService.OrderListGetAll(id, query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, orderListGetAllResponse, pagination)
}

func (h *Handler) OrderListGetAllByMerchantID(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewOrderListQuery(c)

	orderListGetAllResponse, pagination, err := h.orderService.OrderListGetAllByMerchantID(id, query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, orderListGetAllResponse, pagination)
}

func (h *Handler) OrderListDetailsGetByID(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	orderListDetailsGetAllResponse, err := h.orderService.OrderListDetailsGetByID(id, orderID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, orderListDetailsGetAllResponse, nil)
}

func (h *Handler) OrderListDetailsGetByIDByMerchantID(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	orderListDetailsGetAllResponse, err := h.orderService.OrderListDetailsGetByIDByMerchantID(id, orderID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, orderListDetailsGetAllResponse, nil)
}

func (h *Handler) OrderProcessOrder(c *gin.Context) {
	var request *dtos.OrderProcessOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	orderID, err := uuid.Parse(c.Param("order-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	result, err := h.orderService.Update(&internalDtos.OrderProcessOrder{
		CourierServiceID: request.CourierServiceID,
		OrderID:          orderID,
	})
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewOrderProcessOrderResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
