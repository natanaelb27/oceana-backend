package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductSnapshotRepository interface {
	WithTx(tx *gorm.DB) *productSnapshotRepository
	CreateProductSnapshot(model *models.ProductSnapshot) (newModel *models.ProductSnapshot, err error)
}

type productSnapshotRepository struct {
	db *gorm.DB
}

type ProductSnapshotRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductSnapshotRepository(c *ProductSnapshotRepositoryConfig) ProductSnapshotRepository {
	return &productSnapshotRepository{
		db: c.DB,
	}
}

func (u productSnapshotRepository) WithTx(tx *gorm.DB) *productSnapshotRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *productSnapshotRepository) CreateProductSnapshot(newModel *models.ProductSnapshot) (model *models.ProductSnapshot, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
