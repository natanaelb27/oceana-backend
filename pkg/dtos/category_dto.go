package dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
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
