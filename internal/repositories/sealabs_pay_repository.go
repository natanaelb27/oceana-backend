package repositories

import (
	"log"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SealabsPayRepository interface {
	WithTx(tx *gorm.DB) *sealabsPayRepository
	Create(sealabsPay *models.SealabsPay) (model *models.SealabsPay, err error)
	GetAll(userID uuid.UUID) (models []*models.SealabsPay, err error)
	GetById(id uuid.UUID) (model *models.SealabsPay, err error)
	GetByCardNumber(cardNumber string) (model *models.SealabsPay, err error)
	GetByUserIDAndCardNumber(userID uuid.UUID, cardNumber string) (model *models.SealabsPay, err error)
	GetDefaultByUserID(userID uuid.UUID) (model *models.SealabsPay, err error)
	Delete(id uuid.UUID) (model *models.SealabsPay, err error)
	Update(id uuid.UUID, model *models.SealabsPay) (err error)
	SetDefaultFalse(id uuid.UUID, model *models.SealabsPay) (err error)
}

type sealabsPayRepository struct {
	db *gorm.DB
}

type SealabsPayRepositoryConfig struct {
	DB *gorm.DB
}

func NewSealabsPayRepository(c *SealabsPayRepositoryConfig) SealabsPayRepository {
	return &sealabsPayRepository{
		db: c.DB,
	}
}

func (u sealabsPayRepository) WithTx(tx *gorm.DB) *sealabsPayRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *sealabsPayRepository) Create(newModel *models.SealabsPay) (model *models.SealabsPay, err error) {
	queryResult := r.db.Where("account_number = ?", newModel.AccountNumber).FirstOrCreate(&newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	if queryResult.RowsAffected == 0 {
		err = e.ErrSealabsPayCardNumberAlreadyExist
		return
	}

	model = newModel

	return
}

func (r *sealabsPayRepository) GetAll(userID uuid.UUID) (models []*models.SealabsPay, err error) {
	queryResult := r.db.Where("user_id = ?", userID).Order("is_default DESC").Find(&models)
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *sealabsPayRepository) GetById(id uuid.UUID) (model *models.SealabsPay, err error) {
	queryResult := r.db.First(&model, id)
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *sealabsPayRepository) GetByCardNumber(cardNumber string) (model *models.SealabsPay, err error) {
	queryResult := r.db.Where("account_number = ?", cardNumber).First(&model)
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	return
}

func (r *sealabsPayRepository) GetByUserIDAndCardNumber(userID uuid.UUID, cardNumber string) (model *models.SealabsPay, err error) {
	queryResult := r.db.Where("user_id = ? AND account_number = ?", userID, cardNumber).First(&model)
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	return
}

func (r *sealabsPayRepository) GetDefaultByUserID(userID uuid.UUID) (model *models.SealabsPay, err error) {
	queryResult := r.db.Where("user_id = ? AND is_default = true", userID).First(&model)
	if queryResult.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}
	if queryResult.Error != nil {
		err = queryResult.Error
		return
	}
	return
}

func (r *sealabsPayRepository) Delete(id uuid.UUID) (model *models.SealabsPay, err error) {
	result := r.db.Clauses(clause.Returning{}).Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *sealabsPayRepository) Update(id uuid.UUID, model *models.SealabsPay) (err error) {
	queryResult := r.db.Where("id = ?", id).Updates(&model)
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

func (r *sealabsPayRepository) SetDefaultFalse(id uuid.UUID, model *models.SealabsPay) (err error) {
	queryResult := r.db.Model(&model).Where("id = ?", id).Update("is_default", false)
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
