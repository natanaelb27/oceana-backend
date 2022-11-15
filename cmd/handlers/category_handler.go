package handlers

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CategoryGetAll(c *gin.Context) {
	results, err := h.categoryService.GetAll()
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
