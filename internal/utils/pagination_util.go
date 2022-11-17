package utils

import (
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      uint `json:"limit"`
	Page       uint `json:"page"`
	TotalRows  uint `json:"total_rows"`
	TotalPages uint `json:"total_pages"`
}

func (p *Pagination) GetLimit() uint {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() uint {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetOffset() uint {
	return (p.GetPage() - 1) * p.GetLimit()
}

func Paginate(pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Count(&totalRows)

	pagination.TotalRows = uint(totalRows)
	pagination.TotalPages = uint(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(pagination.GetOffset())).Limit(int(pagination.GetLimit()))
	}
}
