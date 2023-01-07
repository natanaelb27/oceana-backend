package repositories

import (
	"log"

	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type RefundMediaRepository interface {
	WithTx(tx *gorm.DB) *refundMediaRepository
	Create(newModel []*models.RefundMedia) (model []*models.RefundMedia, err error)
}

type refundMediaRepository struct {
	db *gorm.DB
}

type RefundMediaRepositoryConfig struct {
	DB *gorm.DB
}

func NewRefundMediaRepository(c *RefundMediaRepositoryConfig) RefundMediaRepository {
	return &refundMediaRepository{
		db: c.DB,
	}
}

func (u refundMediaRepository) WithTx(tx *gorm.DB) *refundMediaRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *refundMediaRepository) Create(newModel []*models.RefundMedia) (model []*models.RefundMedia, err error) {
	queryResult := r.db.Create(&newModel)

	if queryResult.RowsAffected == 0 {
		err = e.ErrInsertFailed
		return
	}

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
