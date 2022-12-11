package cronservice

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"github.com/go-co-op/gocron"
)

type Cron struct {
	scheduler            gocron.Scheduler
	roleService          services.RoleService
	orderTrackingService services.OrderTrackingService
	refundService        services.RefundService
}

type CronConfig struct {
	Scheduler            gocron.Scheduler
	RoleService          services.RoleService
	OrderTrackingService services.OrderTrackingService
	RefundService        services.RefundService
}

func New(c *CronConfig) *Cron {
	return &Cron{
		scheduler:            c.Scheduler,
		roleService:          c.RoleService,
		orderTrackingService: c.OrderTrackingService,
		refundService:        c.RefundService,
	}
}
