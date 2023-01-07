package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func (h *Handler) CategoryGetAll(c *gin.Context) {
	query := dtos.NewCategoryQuery(c)

	results, err := h.categoryService.GetAll(query)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	var responseData []*dtos.CategoryGetResponse = make([]*dtos.CategoryGetResponse, 0)

	for _, result := range results {
		responseData = append(responseData, dtos.NewCategoryGetResponse(result))
	}

	utils.ResponseSuccess(c, responseData, nil)
}

func (h *Handler) CategoryAncestorsGetByID(c *gin.Context) {
	categoryIDUint64, err := strconv.ParseUint(c.Param("category-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	categoryID := uint(categoryIDUint64)

	result, err := h.categoryService.GetAncestorsByID(categoryID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, result, nil)
}

func (h *Handler) GetMerchantCategories(c *gin.Context) {
	merchantIDUint64, err := strconv.ParseUint(c.Param("merchant-id"), 10, 64)
	if err != nil {
		utils.ResponseError(c, e.ErrInvalidIdentifier)
		return
	}
	merchantID := uint(merchantIDUint64)

	result, err := h.categoryService.GetMerchantCategories(merchantID)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, result, nil)
}
