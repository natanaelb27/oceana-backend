package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ProductItemCreate(c *gin.Context) {
	productIDUint64, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidParams)
		return
	}

	productID := uint(productIDUint64)

	var request *dtos.ProductItemCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, err)
		return
	}

	model := request.ToModel(productID)

	result, err := h.productItemService.CreateProductItem(model)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	responseData := dtos.NewProductItemCreateResponse(result)

	utils.ResponseSuccess(c, responseData, nil)

}
