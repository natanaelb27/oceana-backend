package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

func (h *Handler) OrderListGetAll(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewOrderListQuery(c)

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

func (h *Handler) MerchantOrderListDetailsGetByID(c *gin.Context) {
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

	orderListDetailsGetAllResponse, err := h.orderService.OrderListDetailsGetByMerchantID(id, orderID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, orderListDetailsGetAllResponse, nil)
}

func (h *Handler) OrderProcessDelivery(c *gin.Context) {
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

	tx := c.MustGet("tx").(*gorm.DB)

	result, err := h.orderService.WithTx(tx).UpdateProcessDelivery(orderID, id)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewOrderProcessDeliveryResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) OrderProcessOrder(c *gin.Context) {
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

	tx := c.MustGet("tx").(*gorm.DB)

	result, err := h.orderService.WithTx(tx).UpdateProcessOrder(orderID, id)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewUpdateOrderResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) SellerCancelOrder(c *gin.Context) {
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

	tx := c.MustGet("tx").(*gorm.DB)

	result, err := h.orderService.WithTx(tx).SellerCancelOrder(orderID, id)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewUpdateOrderResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) OrderCompleteOrder(c *gin.Context) {
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

	tx := c.MustGet("tx").(*gorm.DB)

	result, err := h.orderService.WithTx(tx).UpdateCompleteOrder(orderID, id)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewUpdateOrderResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) SellerClaimFunds(c *gin.Context) {
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

	tx := c.MustGet("tx").(*gorm.DB)

	result, err := h.orderService.WithTx(tx).SellerClaimFunds(orderID, id)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewUpdateOrderResponse(result)

	utils.ResponseSuccess(c, responseData, nil)
}
