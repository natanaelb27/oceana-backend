package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"github.com/google/uuid"
)

type OrderTrackingService interface {
	CronOrderTrackingUpdater() error
}

type orderTrackingService struct {
	orderTrackingRepository r.OrderTrackingRepository
	orderRepository         r.OrderRepository
}

type OrderTrackingServiceConfig struct {
	OrderTrackingRepository r.OrderTrackingRepository
	OrderRepository         r.OrderRepository
}

func NewOrderTrackingService(c *OrderTrackingServiceConfig) OrderTrackingService {
	return &orderTrackingService{
		orderTrackingRepository: c.OrderTrackingRepository,
		orderRepository:         c.OrderRepository,
	}
}

func (p *orderTrackingService) CronOrderTrackingUpdater() error {
	orderTrackings, err := p.orderTrackingRepository.GetByOrderStatus(constants.OutForDeliveryStatus)
	if err != nil {
		return err
	}

	newOrderTrackings := make([]*models.OrderTracking, 0)

	for _, v := range orderTrackings {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		interval, err := strconv.Atoi(os.Getenv("CRON_ORDER_TRACKING_UPDATE_INTERVAL"))
		if err != nil {
			interval = constants.CronOrderTrackingUpdateInterval
		}
		fmt.Println(v.OrderID, !time.Now().After(v.HistoryDatetime.Add(time.Minute*time.Duration(interval))))
		fmt.Println(v.Title)
		// if time after the last history date time hasnt pass 15 minutes then continue to next iteration
		if !time.Now().After(v.HistoryDatetime.Add(time.Minute * time.Duration(interval))) {
			continue
		}

		newOrderTracking := &models.OrderTracking{}
		newOrder := &models.Order{}

		switch v.Title {
		// case when ... then update order trackings to next status
		case constants.CronOrderTrackingTitleSentFromOriginWarehouse:
			newOrderTracking = &models.OrderTracking{
				Title:       constants.CronOrderTrackingTitleWithCourier,
				Description: constants.CronOrderTrackingDescriptionWithCourier,
			}
		// case handled by courier then update order trackings and order to be delivered
		case constants.CronOrderTrackingTitleWithCourier:
			newOrderTracking = &models.OrderTracking{
				Title:       constants.CronOrderTrackingTitleDelivered,
				Description: constants.CronOrderTrackingDescriptionDelivered,
			}
			newOrder = &models.Order{
				ID:     v.OrderID,
				Status: constants.DeliveredStatus,
			}
		}

		// if there is an update
		if newOrderTracking.Title != "" {
			newOrderTracking.OrderID = v.OrderID
			newOrderTracking.HistoryDatetime = time.Now()
			newOrderTrackings = append(newOrderTrackings, newOrderTracking)
		}
		if newOrder.ID != uuid.Nil {
			_, err := p.orderRepository.Update(newOrder.ID, newOrder)
			if err != nil {
				return err
			}
		}
	}

	if len(newOrderTrackings) != 0 {
		_, err = p.orderTrackingRepository.Create(newOrderTrackings)
		if err != nil {
			return err
		}
	}

	return nil
}
