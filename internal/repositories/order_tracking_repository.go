package repositories

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderTrackingRepository interface {
	GetByOrderID(id uuid.UUID) ([]*internalDtos.OrderTrackingHistory, error)
	Create(newModel *models.OrderTracking) (model *models.OrderTracking, err error)
	GetByOrderStatus(orderStatus string) (result []*internalDtos.COTUOrderTracking, err error)
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

func (r *orderTrackingRepository) GetByOrderID(id uuid.UUID) ([]*internalDtos.OrderTrackingHistory, error) {
	var orderTrackingHistory []*internalDtos.OrderTrackingHistory

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

func (r *orderTrackingRepository) GetByOrderStatus(orderStatus string) (result []*internalDtos.COTUOrderTracking, err error) {
	sqOrderTrackings := r.db.
		Table("order_trackings ot").
		Select(`
			ROW_NUMBER() OVER (PARTITION BY ot.order_id ORDER BY ot.history_datetime DESC) AS row_number,
			ot.id AS id,
			ot.order_id AS order_id,
			ot.history_datetime AS history_datetime,
			ot.title AS title,
			ot.description AS description, 
			o.id AS order_id,
			o.status AS order_status
		`).
		Joins("JOIN orders o ON ot.order_id = o.id").
		Where(`o.status = ?`, orderStatus).
		Order("ot.history_datetime DESC")
	queryResult := r.db.
		Select(`
			*
		`).
		Table("(?) o", sqOrderTrackings).
		Where("o.row_number = ?", 1).
		Find(&result)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return result, nil
}

func (r *orderTrackingRepository) Create(newModel *models.OrderTracking) (model *models.OrderTracking, err error) {
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
