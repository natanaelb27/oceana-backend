package handlers

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
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
