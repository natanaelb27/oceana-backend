package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
)

type CityService interface {
	GetAll(query *dtos.CitiesQuery) ([]*dtos.CityResponse, error)
}

type cityService struct {
	cityRepository r.CityRepository
}

type CityServiceConfig struct {
	CityRepository r.CityRepository
}

func NewCityService(c *CityServiceConfig) CityService {
	return &cityService{
		cityRepository: c.CityRepository,
	}
}

func (s *cityService) GetAll(query *dtos.CitiesQuery) ([]*dtos.CityResponse, error) {
	cities, err := s.cityRepository.CitiesGetAll(query)
	if err != nil {
		return nil, err
	}

	return dtos.NewCitiesResponse(cities), nil
}
