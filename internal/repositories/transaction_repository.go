package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(newModel *models.Transaction) (model *models.Transaction, err error)
}

type transactionRepository struct {
	db *gorm.DB
}

type TransactionRepositoryConfig struct {
	DB *gorm.DB
}

func NewTransactionRepository(c *TransactionRepositoryConfig) TransactionRepository {
	return &transactionRepository{
		db: c.DB,
	}
}

func (r *transactionRepository) Create(newModel *models.Transaction) (model *models.Transaction, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
