package dtos

type AddressDetails struct {
	Detail      string `json:"detail"`
	SubDistrict string `json:"sub_district"`
	ZipCode     string `json:"zip_code"`
	District    string `json:"district"`
	City        string `json:"city"`
	Province    string `json:"province"`
}
