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

	refundRepository := repositories.NewRefundRepository(&repositories.RefundRepositoryConfig{
		DB: db.Get(),
	})

	refundDetailRepository := repositories.NewRefundDetailRepository(&repositories.RefundDetailRepositoryConfig{
		DB: db.Get(),
	})

	transactionRepository := repositories.NewTransactionRepository(&repositories.TransactionRepositoryConfig{
		DB: db.Get(),
	})

	orderDetailRepository := repositories.NewOrderDetailRepository(&repositories.OrderDetailRepositoryConfig{
		DB: db.Get(),
	})

	orderRepository := repositories.NewOrderRepository(&repositories.OrderRepositoryConfig{
		DB: db.Get(),
	})

	orderTrackingRepository := repositories.NewOrderTrackingRepository(&repositories.OrderTrackingRepositoryConfig{
		DB: db.Get(),
	})

	walletRepository := repositories.NewWalletRepository(&repositories.WalletRepositoryConfig{
		DB: db.Get(),
	})

	walletTransactionPaymentRepository := repositories.NewWalletTransactionPaymentRepository(&repositories.WalletTransactionPaymentRepositoryConfig{
		DB: db.Get(),
	})

	sealabsPayTransactionPaymentRepository := repositories.NewSealabsPayTransactionPaymentRepository(&repositories.SealabsPayTransactionPaymentRepositoryConfig{
		DB: db.Get(),
	})

	orderTrackingService := services.NewOrderTrackingService(&services.OrderTrackingServiceConfig{
		OrderTrackingRepository: orderTrackingRepository,
		OrderRepository:         orderRepository,
	})

	roleService := services.NewRoleService(&services.RoleServiceConfig{
		RoleRepository: roleRepository,
	})

	orderService := services.NewOrderService(&services.OrderServiceConfig{
		OrderDetailRepository:                  orderDetailRepository,
		OrderRepository:                        orderRepository,
		OrderTrackingRepository:                orderTrackingRepository,
		RefundRepository:                       refundRepository,
		RefundDetailRepository:                 refundDetailRepository,
		TransactionRepository:                  transactionRepository,
		WalletRepository:                       walletRepository,
		WalletTransactionPaymentRepository:     walletTransactionPaymentRepository,
		SealabsPayTransactionPaymentRepository: sealabsPayTransactionPaymentRepository,
	})

	refundService := services.NewRefundService(&services.RefundServiceConfig{
		RefundRepository:        refundRepository,
		OrderRepository:         orderRepository,
		OrderTrackingRepository: orderTrackingRepository,
	})

	c := New(&CronConfig{
		Scheduler:            *gocron.NewScheduler(time.UTC),
		RoleService:          roleService,
		OrderService:         orderService,
		OrderTrackingService: orderTrackingService,
		RefundService:        refundService,
	})

	// cron that runs per 15 second
	c.Per15Seconds()

	c.scheduler.StartAsync()
}
