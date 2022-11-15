package dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
)

type CategoryGetResponse struct {
	ID       uint             `json:"id" gorm:"primary_key"`
	Name     string           `json:"name"`
	Level    int              `json:"level"`
	Category []*dtos.Category `json:"category"`
}

func NewCategoryGetResponse(source *dtos.Category) (target *CategoryGetResponse) {
	target = &CategoryGetResponse{
		ID:       source.ID,
		Name:     source.Name,
		Level:    source.Level,
		Category: source.Category,
	}
	return
}
