package dtos

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type SubDistrictResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	ZipCode string `json:"zip_code"`
}

func NewSubDistrictResponse(source *models.SubDistrict) (target *SubDistrictResponse) {
	target = &SubDistrictResponse{
		ID:      source.ID,
		Name:    source.Name,
		ZipCode: source.ZipCode,
	}
	return
}

type GetSubDistrictsResponse struct {
	Provinces    []*ProvinceResponse    `json:"provinces"`
	Cities       []*CityResponse        `json:"cities"`
	Districts    []*DistrictResponse    `json:"districts"`
	SubDistricts []*SubDistrictResponse `json:"sub_districts"`
}

func NewGetSubDistrictsResponse(
	provinceModels []*models.Province,
	cityModels []*models.City,
	districtModels []*models.District,
	subDistrictModels []*models.SubDistrict,
) (target *GetSubDistrictsResponse) {
	provinces := []*ProvinceResponse{}
	cities := []*CityResponse{}
	districts := []*DistrictResponse{}
	subDistricts := []*SubDistrictResponse{}

	for _, model := range provinceModels {
		provinces = append(provinces, NewProvinceResponse(model))
	}

	for _, model := range cityModels {
		cities = append(cities, NewCityResponse(model))
	}

	for _, model := range districtModels {
		districts = append(districts, NewDistrictResponse(model))
	}

	for _, model := range subDistrictModels {
		subDistricts = append(subDistricts, NewSubDistrictResponse(model))
	}

	target = &GetSubDistrictsResponse{
		Provinces:    provinces,
		Cities:       cities,
		Districts:    districts,
		SubDistricts: subDistricts,
	}
	return
}

type SubDistrictQuery struct {
	Province uint `json:"province"`
	City     uint `json:"city"`
	District uint `json:"cistrict"`
}

func NewSubDistrictQuery(c *gin.Context) *SubDistrictQuery {
	var query SubDistrictQuery

	province := c.Query("province")
	city := c.Query("city")
	district := c.Query("district")

	provinceInt, err := strconv.Atoi(province)
	if err != nil || provinceInt <= 0 {
		provinceInt = 0
	}

	cityInt, err := strconv.Atoi(city)
	if err != nil || cityInt <= 0 {
		cityInt = 0
	}

	districtInt, err := strconv.Atoi(district)
	if err != nil || districtInt <= 0 {
		districtInt = 0
	}

	query.Province = uint(provinceInt)
	query.City = uint(cityInt)
	query.District = uint(districtInt)

	return &query
}
