package dtos

type DistrictResponse struct {
	ID     uint         `json:"id"`
	Name   string       `json:"name"`
	CityID uint         `json:"city_id"`
	City   CityResponse `json:"city"`
}
