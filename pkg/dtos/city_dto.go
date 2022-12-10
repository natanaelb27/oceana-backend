package dtos

import "git.garena.com/azmi-adhani/group-project-backend/internal/models"

type CityResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewCityResponse(source *models.City) (target *CityResponse) {
	target = &CityResponse{
		ID:   source.ID,
		Name: source.Name,
	}
	return
}

func NewCitiesResponse(source []*models.City) (target []*CityResponse) {
	for _, result := range source {
		target = append(target, &CityResponse{
			ID:   result.ID,
			Name: result.Name,
		})
	}

	return
}
