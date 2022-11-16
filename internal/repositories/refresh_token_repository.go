package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(refreshToken *models.RefreshToken) (*models.RefreshToken, error)
}

type refreshTokenRepository struct {
	db *gorm.DB
}

type RefreshTokenRepositoryConfig struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(c *RefreshTokenRepositoryConfig) RefreshTokenRepository {
	return &refreshTokenRepository{
		db: c.DB,
	}
}

func (r *refreshTokenRepository) Create(newModel *models.RefreshToken) (model *models.RefreshToken, err error) {
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
