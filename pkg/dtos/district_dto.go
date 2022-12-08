package dtos

import "git.garena.com/azmi-adhani/group-project-backend/internal/models"

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
