package repositories

import (
	"log"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type AddressSnapshotRepository interface {
	WithTx(tx *gorm.DB) *addressSnapshotRepository
	GetByOrderID(id uuid.UUID) (*dtos.AddressSnapshotDetails, error)
	Create(id uuid.UUID) (*models.AddressSnapshot, error)
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

func (u addressSnapshotRepository) WithTx(tx *gorm.DB) *addressSnapshotRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
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

func (r *addressSnapshotRepository) Create(id uuid.UUID) (*models.AddressSnapshot, error) {
	var newModel *models.AddressSnapshot

	queryResult := r.db.Model(&models.Address{}).
		Select(`
			detail AS address,
			sub_district_id
		`).
		Where("id = ?", id).
		First(&newModel)

	if queryResult.Error != nil {
		if queryResult.RowsAffected == 0 {
			return nil, e.ErrDataNotFound
		}

		return nil, e.ErrDatabaseOperationFailed
	}

	queryResult = r.db.Create(&newModel)

	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return newModel, nil
}
