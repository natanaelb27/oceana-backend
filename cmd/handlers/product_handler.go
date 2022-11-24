package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ProductDetailsGetByID(c *gin.Context) {
	productIDUint64, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	productID := uint(productIDUint64)

	orderListDetailsGetAllResponse, err := h.productService.ProductDetailsGetByID(productID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, orderListDetailsGetAllResponse, nil)
}
