package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
)

type CityService interface {
	GetAll() ([]*dtos.CityResponse, error)
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

func (s *cityService) GetAll() ([]*dtos.CityResponse, error) {
	cities, err := s.cityRepository.CitiesGetAll()
	if err != nil {
		return nil, err
	}

	return dtos.NewCitiesResponse(cities), nil
}
