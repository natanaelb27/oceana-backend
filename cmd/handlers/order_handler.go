package handlers

import (
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

	orderListGetAllResponse, err := h.orderService.OrderListGetAll(id, query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, orderListGetAllResponse, nil)
}

func (h *Handler) OrderListDetailsGetAll(c *gin.Context) {
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

	orderListDetailsGetAllResponse, err := h.orderService.OrderListDetailsGetAll(id, orderID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, orderListDetailsGetAllResponse, nil)
}
