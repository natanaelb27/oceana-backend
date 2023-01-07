package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func (h *Handler) PendingPaymentGetAllByUser(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	results, err := h.transactionService.PendingPaymentGetAllByUser(id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, nil)
}
