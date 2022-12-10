package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderTrackingRepository interface {
	GetByOrderID(id uuid.UUID) ([]*dtos.OrderTrackingHistory, error)
	Create(newModel []*models.OrderTracking) (model []*models.OrderTracking, err error)
}

type orderTrackingRepository struct {
	db *gorm.DB
}

type OrderTrackingRepositoryConfig struct {
	DB *gorm.DB
}

func NewOrderTrackingRepository(c *OrderTrackingRepositoryConfig) OrderTrackingRepository {
	return &orderTrackingRepository{
		db: c.DB,
	}
}

func (r *orderTrackingRepository) GetByOrderID(id uuid.UUID) ([]*dtos.OrderTrackingHistory, error) {
	var orderTrackingHistory []*dtos.OrderTrackingHistory

	queryResult := r.db.Model(&models.OrderTracking{}).
		Select(`
			order_trackings.history_datetime AS datetime,
			order_trackings.title AS title,
			order_trackings.description AS description
		`).
		Joins("JOIN orders ON order_trackings.order_id = orders.id").
		Where("order_trackings.order_id = ?", id).
		Order("history_datetime DESC").
		Find(&orderTrackingHistory)

	if queryResult.Error != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return orderTrackingHistory, nil
}
func (r *orderTrackingRepository) Create(newModel []*models.OrderTracking) (model []*models.OrderTracking, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
