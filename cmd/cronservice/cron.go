package cronservice

import (
	"github.com/go-co-op/gocron"
	"github.com/natanaelb27/oceana-backend/internal/services"
)

type Cron struct {
	scheduler            gocron.Scheduler
	roleService          services.RoleService
	orderService         services.OrderService
	orderTrackingService services.OrderTrackingService
	refundService        services.RefundService
}

type CronConfig struct {
	Scheduler            gocron.Scheduler
	RoleService          services.RoleService
	OrderService         services.OrderService
	OrderTrackingService services.OrderTrackingService
	RefundService        services.RefundService
}

func New(c *CronConfig) *Cron {
	return &Cron{
		scheduler:            c.Scheduler,
		roleService:          c.RoleService,
		orderService:         c.OrderService,
		orderTrackingService: c.OrderTrackingService,
		refundService:        c.RefundService,
	}
}
