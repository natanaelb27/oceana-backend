package dtos

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"github.com/gin-gonic/gin"
)

type CategoryGetResponse struct {
	ID       uint             `json:"id" gorm:"primary_key"`
	Image    string           `json:"image"`
	Name     string           `json:"name"`
	Level    int              `json:"level"`
	Category []*dtos.Category `json:"category"`
}

func NewCategoryGetResponse(source *dtos.Category) (target *CategoryGetResponse) {
	target = &CategoryGetResponse{
		ID:       source.ID,
		Image:    source.Image,
		Name:     source.Name,
		Level:    source.Level,
		Category: source.Category,
	}
	return
}

type CategoryQuery struct {
	CategoryID uint `json:"category-id"`
	LevelDepth int  `json:"level-depth"`
}

func NewCategoryQuery(c *gin.Context) *CategoryQuery {
	var query CategoryQuery

	CategoryId := c.Query("category-id")
	LevelDepth := c.Query("level-depth")

	CategoryIdInt, err := strconv.Atoi(CategoryId)
	if err != nil || CategoryIdInt <= 0 {
		CategoryIdInt = 0
	}
	query.CategoryID = uint(CategoryIdInt)

	LevelDepthInt, err := strconv.Atoi(LevelDepth)
	if err != nil || LevelDepthInt <= 0 {
		LevelDepthInt = -1
	}
	query.LevelDepth = LevelDepthInt

	return &query
}
