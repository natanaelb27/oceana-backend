package cronservice

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"github.com/go-co-op/gocron"
)

func Init() {
	roleRepository := repositories.NewRoleRepository(&repositories.RoleRepositoryConfig{
		DB: db.Get(),
	})
	roleService := services.NewRoleService(&services.RoleServiceConfig{
		RoleRepository: roleRepository,
	})

	orderRepository := repositories.NewOrderRepository(&repositories.OrderRepositoryConfig{
		DB: db.Get(),
	})
	orderTrackingRepository := repositories.NewOrderTrackingRepository(&repositories.OrderTrackingRepositoryConfig{
		DB: db.Get(),
	})
	orderTrackingService := services.NewOrderTrackingService(&services.OrderTrackingServiceConfig{
		OrderTrackingRepository: orderTrackingRepository,
		OrderRepository:         orderRepository,
	})

	refundRepository := repositories.NewRefundRepository(&repositories.RefundRepositoryConfig{
		DB: db.Get(),
	})
	refundService := services.NewRefundService(&services.RefundServiceConfig{
		RefundRepository: refundRepository,
	})

	c := New(&CronConfig{
		Scheduler:            *gocron.NewScheduler(time.UTC),
		RoleService:          roleService,
		OrderTrackingService: orderTrackingService,
		RefundService:        refundService,
	})

	// cron that runs per 15 second
	c.Per15Seconds()

	c.scheduler.StartAsync()
}
