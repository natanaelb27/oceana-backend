package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ProductDetailsGetByID(c *gin.Context) {
	productIDUint64, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	productID := uint(productIDUint64)

	productDetailsGetByIDResponse, err := h.productService.ProductDetailsGetByID(productID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, productDetailsGetByIDResponse, nil)
}

func (h *Handler) GetPurchaseQuantitiesByID(c *gin.Context) {
	id := utils.GinGetUUID(c, "id")
	if id == uuid.Nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}

	query := dtos.NewPurchaseQuantitiesQuery(c)

	getPurchaseQuantitiesByIDResponse, err := h.productService.GetPurchaseQuantitiesByID(id, query)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, getPurchaseQuantitiesByIDResponse, nil)
}
