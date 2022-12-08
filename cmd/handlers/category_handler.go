package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
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
