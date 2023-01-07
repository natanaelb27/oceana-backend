package repositories

import (
	"log"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/models"
	er "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type EmailRepository interface {
	WithTx(tx *gorm.DB) *emailRepository
	Create(email *models.Email) (*models.Email, error)
	GetById(id uuid.UUID) (model *models.Email, err error)
	GetByEmail(emailAddress string) (model *models.Email, err error)
}

type emailRepository struct {
	db *gorm.DB
}

type EmailRepositoryConfig struct {
	DB *gorm.DB
}

func NewEmailRepository(c *EmailRepositoryConfig) EmailRepository {
	return &emailRepository{
		db: c.DB,
	}
}

func (u emailRepository) WithTx(tx *gorm.DB) *emailRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (e *emailRepository) Create(newModel *models.Email) (model *models.Email, err error) {
	queryResult := e.db.Create(&newModel)
	if queryResult.RowsAffected == 0 {
		err = er.ErrInsertFailed
		return
	}
	if queryResult.Error != nil {
		err = er.ErrDatabaseOperationFailed
		return
	}
	model = newModel
	return
}

func (e *emailRepository) GetById(id uuid.UUID) (model *models.Email, err error) {
	result := e.db.First(&model, id)
	if result.RowsAffected == 0 {
		err = er.ErrDataNotFound
		return
	}
	return
}

func (e *emailRepository) GetByEmail(emailAddress string) (model *models.Email, err error) {
	result := e.db.Model(&models.Email{}).Where("address = ?", emailAddress).Find(&model)
	if result.RowsAffected == 0 {
		err = er.ErrDataNotFound
		return
	}
	return
}
