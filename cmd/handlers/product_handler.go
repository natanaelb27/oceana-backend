package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ProductGetAll(c *gin.Context) {
	query := product_dtos.NewProductQuery(c)

	results, pagination, err := h.productService.GetAll(query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, pagination)
}

func (h *Handler) ShopProductGetAll(c *gin.Context) {
	query, err := product_dtos.NewShopProductQuery(c)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	results, pagination, err := h.productService.GetAll(query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, results, pagination)
}
