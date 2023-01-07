package services

import (
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
)

type SubDistrictService interface {
	GetAll(query *dtos.SubDistrictQuery) (*dtos.GetSubDistrictsResponse, error)
}

type subDistrictService struct {
	provinceRepository    r.ProvinceRepository
	cityRepository        r.CityRepository
	districtRepository    r.DistrictRepository
	subDistrictRepository r.SubDistrictRepository
}

type SubDistrictServiceConfig struct {
	ProvinceRepository    r.ProvinceRepository
	CityRepository        r.CityRepository
	DistrictRepository    r.DistrictRepository
	SubDistrictRepository r.SubDistrictRepository
}

func NewSubDistrictService(c *SubDistrictServiceConfig) SubDistrictService {
	return &subDistrictService{
		provinceRepository:    c.ProvinceRepository,
		cityRepository:        c.CityRepository,
		districtRepository:    c.DistrictRepository,
		subDistrictRepository: c.SubDistrictRepository,
	}
}

func (s *subDistrictService) GetAll(query *dtos.SubDistrictQuery) (*dtos.GetSubDistrictsResponse, error) {
	provinces, err := s.provinceRepository.GetAll()
	if err != nil {
		return nil, err
	}

	cities, err := s.cityRepository.GetAll(query)
	if err != nil {
		return nil, err
	}

	districts, err := s.districtRepository.GetAll(query)
	if err != nil {
		return nil, err
	}

	subDistricts, err := s.subDistrictRepository.GetAll(query)
	if err != nil {
		return nil, err
	}

	return dtos.NewGetSubDistrictsResponse(provinces, cities, districts, subDistricts), nil
}
