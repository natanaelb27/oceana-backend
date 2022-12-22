package dtos

import "github.com/lib/pq"

type CourierAndCourierServices struct {
	Code         string         `json:"codes"`
	ServiceCodes pq.StringArray `json:"service_codes" gorm:"type:text[]"`
}
