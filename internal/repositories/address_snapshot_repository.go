package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressSnapshotRepository interface {
	GetByOrderID(id uuid.UUID) (*dtos.AddressSnapshotDetails, error)
}

type addressSnapshotRepository struct {
	db *gorm.DB
}

type AddressSnapshotRepositoryConfig struct {
	DB *gorm.DB
}

func NewAddressSnapshotRepository(c *AddressSnapshotRepositoryConfig) AddressSnapshotRepository {
	return &addressSnapshotRepository{
		db: c.DB,
	}
}

func (r *addressSnapshotRepository) GetByOrderID(id uuid.UUID) (*dtos.AddressSnapshotDetails, error) {
	var addressSnapshotDetails *dtos.AddressSnapshotDetails

	queryResult := r.db.Model(&models.AddressSnapshot{}).
		Select(`
			address_snapshots.address AS detail,
			sub_districts.name AS sub_district,
			sub_districts.zip_code AS zip_code,
			districts.name AS district,
			cities.name AS city,
			provinces.name AS province
		`).
		Joins("JOIN orders ON address_snapshots.id = orders.address_snapshot_id").
		Joins("JOIN sub_districts ON address_snapshots.sub_district_id = sub_districts.id").
		Joins("JOIN districts ON sub_districts.district_id = districts.id").
		Joins("JOIN cities ON districts.city_id = cities.id").
		Joins("JOIN provinces ON cities.province_id = provinces.id").
		Where("orders.id = ?", id).
		Find(&addressSnapshotDetails)

	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return addressSnapshotDetails, nil
}
