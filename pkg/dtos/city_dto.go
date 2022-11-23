package dtos

type CityResponse struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	ProvinceID uint             `json:"province_id"`
	Province   ProvinceResponse `json:"province"`
}
