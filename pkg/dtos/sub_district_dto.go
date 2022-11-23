package dtos

type SubDistrictResponse struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	DistrictID uint             `json:"district_id"`
	District   DistrictResponse `json:"district"`
	ZipCode    string           `json:"zip_code"`
}
