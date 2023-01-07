package dtos

import "github.com/natanaelb27/oceana-backend/internal/models"

type DistrictResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewDistrictResponse(source *models.District) (target *DistrictResponse) {
	target = &DistrictResponse{
		ID:   source.ID,
		Name: source.Name,
	}
	return
}
