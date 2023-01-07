package dtos

import (
	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/internal/models"
)

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

type CitiesQuery struct {
	Search string `json:"search"`
}

func NewCitiesQuery(c *gin.Context) *CitiesQuery {
	var query CitiesQuery

	search := c.Query("search")

	query.Search = search

	return &query
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
