package repositories

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SealabsPayTransactionPaymentRepository interface {
	WithTx(tx *gorm.DB) *sealabsPayTransactionPaymentRepository
	Create(newModel *models.SealabsPayTransactionPayment) (model *models.SealabsPayTransactionPayment, err error)
	GetByTransactionID(id uuid.UUID) (model *models.SealabsPayTransactionPayment, err error)
	Update(model *models.SealabsPayTransactionPayment) (err error)
}

type sealabsPayTransactionPaymentRepository struct {
	db *gorm.DB
}

type SealabsPayTransactionPaymentRepositoryConfig struct {
	DB *gorm.DB
}

func NewSealabsPayTransactionPaymentRepository(c *SealabsPayTransactionPaymentRepositoryConfig) SealabsPayTransactionPaymentRepository {
	return &sealabsPayTransactionPaymentRepository{
		db: c.DB,
	}
}

func (u sealabsPayTransactionPaymentRepository) WithTx(tx *gorm.DB) *sealabsPayTransactionPaymentRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *sealabsPayTransactionPaymentRepository) Create(newModel *models.SealabsPayTransactionPayment) (model *models.SealabsPayTransactionPayment, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}

func (r *sealabsPayTransactionPaymentRepository) GetByTransactionID(id uuid.UUID) (model *models.SealabsPayTransactionPayment, err error) {
	queryResult := r.db.Model(&models.SealabsPayTransactionPayment{}).
		Where("transaction_id = ?", id).
		First(&model)

	if queryResult.Error != nil {
		if queryResult.RowsAffected == 0 {
			err = e.ErrDataNotFound
			return
		}

		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *sealabsPayTransactionPaymentRepository) Update(model *models.SealabsPayTransactionPayment) (err error) {
	queryResult := r.db.Updates(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

