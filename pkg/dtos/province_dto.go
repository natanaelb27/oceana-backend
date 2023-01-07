package dtos

import "github.com/natanaelb27/oceana-backend/internal/models"

type ProvinceResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewProvinceResponse(source *models.Province) (target *ProvinceResponse) {
	target = &ProvinceResponse{
		ID:   source.ID,
		Name: source.Name,
	}
	return
}
