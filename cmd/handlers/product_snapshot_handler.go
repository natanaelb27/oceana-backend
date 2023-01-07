package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func (h *Handler) ProductSnapshotDetailsGetByOrderDetailID(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	orderDetailID, err := uuid.Parse(c.Param("order-detail-id"))
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidOrderDetailID)
		return
	}

	response, err := h.productSnapshotService.GetByOrderDetailID(id, orderDetailID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, response, nil)
}
