package services

import (
	"math"
	"os"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"

	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/google/uuid"
)

type OrderService interface {
	WithTx(tx *gorm.DB) *orderService
	OrderListGetAll(id uuid.UUID, query *dtos.OrderListQuery) ([]*dtos.OrderListGetAllResponse, *utils.Pagination, error)
	OrderListGetAllByMerchantID(userID uuid.UUID, query *dtos.OrderListQuery) ([]*dtos.OrderListGetAllResponse, *utils.Pagination, error)
	OrderListDetailsGetByID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error)
	OrderListDetailsGetByMerchantID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error)
	UpdateProcessDelivery(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error)
	UpdateProcessOrder(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error)
	CronPaymentExpiration() error
	CronPendingSellerExpiration() error
	CronOrderProcessExpiration() error
	CronAutoCompleteDelivered() error
	SellerCancelOrder(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error)
	UpdateCompleteOrder(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error)
	SellerClaimFunds(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error)
}

type orderService struct {
	orderRepository                        r.OrderRepository
	orderDetailRepository                  r.OrderDetailRepository
	addressSnapshotRepository              r.AddressSnapshotRepository
	orderTrackingRepository                r.OrderTrackingRepository
	merchantRepository                     r.MerchantRepository
	refundRepository                       r.RefundRepository
	refundDetailRepository                 r.RefundDetailRepository
	transactionRepository                  r.TransactionRepository
	walletTransactionPaymentRepository     r.WalletTransactionPaymentRepository
	walletRepository                       r.WalletRepository
	sealabsPayTransactionPaymentRepository r.SealabsPayTransactionPaymentRepository
	sealabsPayService                      SealabsPayService
}

type OrderServiceConfig struct {
	OrderRepository                        r.OrderRepository
	OrderDetailRepository                  r.OrderDetailRepository
	AddressSnapshotRepository              r.AddressSnapshotRepository
	OrderTrackingRepository                r.OrderTrackingRepository
	MerchantRepository                     r.MerchantRepository
	RefundRepository                       r.RefundRepository
	RefundDetailRepository                 r.RefundDetailRepository
	TransactionRepository                  r.TransactionRepository
	WalletTransactionPaymentRepository     r.WalletTransactionPaymentRepository
	WalletRepository                       r.WalletRepository
	SealabsPayTransactionPaymentRepository r.SealabsPayTransactionPaymentRepository
	SealabsPayService                      SealabsPayService
}

func NewOrderService(c *OrderServiceConfig) OrderService {
	return &orderService{
		orderRepository:                        c.OrderRepository,
		orderDetailRepository:                  c.OrderDetailRepository,
		addressSnapshotRepository:              c.AddressSnapshotRepository,
		orderTrackingRepository:                c.OrderTrackingRepository,
		merchantRepository:                     c.MerchantRepository,
		refundRepository:                       c.RefundRepository,
		refundDetailRepository:                 c.RefundDetailRepository,
		transactionRepository:                  c.TransactionRepository,
		walletTransactionPaymentRepository:     c.WalletTransactionPaymentRepository,
		walletRepository:                       c.WalletRepository,
		sealabsPayTransactionPaymentRepository: c.SealabsPayTransactionPaymentRepository,
		sealabsPayService:                      c.SealabsPayService,
	}
}

func (u orderService) WithTx(tx *gorm.DB) *orderService {
	u.orderRepository = u.orderRepository.WithTx(tx)
	u.orderDetailRepository = u.orderDetailRepository.WithTx(tx)
	u.addressSnapshotRepository = u.addressSnapshotRepository.WithTx(tx)
	u.orderTrackingRepository = u.orderTrackingRepository.WithTx(tx)
	u.merchantRepository = u.merchantRepository.WithTx(tx)
	u.refundRepository = u.refundRepository.WithTx(tx)
	u.refundDetailRepository = u.refundDetailRepository.WithTx(tx)
	u.transactionRepository = u.transactionRepository.WithTx(tx)
	u.walletTransactionPaymentRepository = u.walletTransactionPaymentRepository.WithTx(tx)
	u.walletRepository = u.walletRepository.WithTx(tx)
	u.sealabsPayTransactionPaymentRepository = u.sealabsPayTransactionPaymentRepository.WithTx(tx)
	u.sealabsPayService = u.sealabsPayService.WithTx(tx)

	return &u
}

func (s *orderService) OrderListGetAll(id uuid.UUID, query *dtos.OrderListQuery) ([]*dtos.OrderListGetAllResponse, *utils.Pagination, error) {
	pagination := &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	orderData, pagination, err := s.orderRepository.OrderListGetAll(id, query, pagination)

	if err != nil {
		return nil, nil, err
	}

	orderList := []*internalDtos.OrderList{}

	for _, order := range orderData {
		orderDetails, err := s.orderDetailRepository.OrderDetailListGetByOrderID(order.OrderID)

		if err != nil {
			return nil, nil, err
		}

		orderList = append(orderList, &internalDtos.OrderList{
			TransactionID:                        order.TransactionID,
			OrderID:                              order.OrderID,
			MerchantID:                           order.MerchantID,
			InvoiceNumber:                        order.InvoiceNumber,
			ShopName:                             order.ShopName,
			ShopPhoto:                            order.ShopPhoto,
			UserID:                               order.UserID,
			UserFullName:                         order.UserFullName,
			CourierServiceID:                     order.CourierServiceID,
			CourierDescription:                   order.CourierDescription,
			OrderDetails:                         orderDetails,
			TransactionDatetime:                  order.TransactionDatetime,
			OrderStatus:                          order.OrderStatus,
			TotalPrice:                           order.TotalPrice,
			RefundID:                             order.RefundID,
			OngoingRefundID:                      order.OngoingRefundID,
			ApprovedRefundID:                     order.ApprovedRefundID,
			ClaimStatus:                          order.ClaimStatus,
			IsAllReviewed:                        order.IsAllReviewed,
			SellerConfirmationExpirationDatetime: order.SellerConfirmationExpirationDatetime,
			OnProcessExpirationDatetime:          order.OnProcessExpirationDatetime,
			DeliveredExpirationDatetime:          order.DeliveredExpirationDatetime,
			RefundSellerConfirmationExpirationDatetime:        order.RefundSellerConfirmationExpirationDatetime,
			RefundRejectedAutoCompleteOrderExpirationDatetime: order.RefundRejectedAutoCompleteOrderExpirationDatetime,
		})
	}

	return dtos.NewOrderListGetAllResponse(orderList), pagination, nil
}

func (s *orderService) OrderListGetAllByMerchantID(userID uuid.UUID, query *dtos.OrderListQuery) ([]*dtos.OrderListGetAllResponse, *utils.Pagination, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	pagination := &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	orderData, pagination, err := s.orderRepository.OrderListGetAllByMerchantID(merchant.ID, query, pagination)

	if err != nil {
		return nil, nil, err
	}

	orderList := []*internalDtos.OrderList{}

	for _, order := range orderData {
		orderDetails, err := s.orderDetailRepository.OrderDetailListGetByOrderID(order.OrderID)

		if err != nil {
			return nil, nil, err
		}

		orderList = append(orderList, &internalDtos.OrderList{
			TransactionID:                        order.TransactionID,
			OrderID:                              order.OrderID,
			MerchantID:                           order.MerchantID,
			InvoiceNumber:                        order.InvoiceNumber,
			ShopName:                             order.ShopName,
			ShopPhoto:                            order.ShopPhoto,
			UserID:                               order.UserID,
			UserFullName:                         order.UserFullName,
			CourierServiceID:                     order.CourierServiceID,
			CourierDescription:                   order.CourierDescription,
			OrderDetails:                         orderDetails,
			TransactionDatetime:                  order.TransactionDatetime,
			OrderStatus:                          order.OrderStatus,
			TotalPrice:                           order.TotalPrice,
			RefundID:                             order.RefundID,
			OngoingRefundID:                      order.OngoingRefundID,
			ApprovedRefundID:                     order.ApprovedRefundID,
			ClaimStatus:                          order.ClaimStatus,
			IsAllReviewed:                        order.IsAllReviewed,
			SellerConfirmationExpirationDatetime: order.SellerConfirmationExpirationDatetime,
			OnProcessExpirationDatetime:          order.OnProcessExpirationDatetime,
			DeliveredExpirationDatetime:          order.DeliveredExpirationDatetime,
			RefundSellerConfirmationExpirationDatetime:        order.RefundSellerConfirmationExpirationDatetime,
			RefundRejectedAutoCompleteOrderExpirationDatetime: order.RefundRejectedAutoCompleteOrderExpirationDatetime,
		})
	}

	return dtos.NewOrderListGetAllResponse(orderList), pagination, nil
}

func (s *orderService) OrderListDetailsGetByID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error) {
	orderData, err := s.orderRepository.OrderListDetailsGetByID(orderID)

	if err != nil {
		return nil, err
	}

	if orderData.UserID != userID {
		return nil, e.ErrForbiddenOrder
	}

	addressSnapshot, err := s.addressSnapshotRepository.GetByOrderID(orderID)

	if err != nil {
		return nil, err
	}

	orderData.Address = addressSnapshot

	orderTrackingHistory, err := s.orderTrackingRepository.GetByOrderID(orderID)

	if err != nil {
		return nil, err
	}

	orderData.OrderTrackingHistory = orderTrackingHistory

	orderDetails, err := s.orderDetailRepository.OrderDetailListGetByOrderID(orderData.OrderID)

	if err != nil {
		return nil, err
	}

	orderData.OrderDetails = orderDetails

	return dtos.NewOrderListDetailsGetByIDResponse(orderData), nil
}

func (s *orderService) OrderListDetailsGetByMerchantID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)

	if err != nil {
		return nil, err
	}

	orderData, err := s.orderRepository.OrderListDetailsGetByID(orderID)

	if err != nil {
		return nil, err
	}

	if orderData.MerchantID != merchant.ID {
		return nil, e.ErrForbiddenMerchant
	}

	addressSnapshot, err := s.addressSnapshotRepository.GetByOrderID(orderID)

	if err != nil {
		return nil, err
	}

	orderData.Address = addressSnapshot

	orderTrackingHistory, err := s.orderTrackingRepository.GetByOrderID(orderID)

	if err != nil {
		return nil, err
	}

	orderData.OrderTrackingHistory = orderTrackingHistory

	orderDetails, err := s.orderDetailRepository.OrderDetailListGetByOrderID(orderData.OrderID)

	if err != nil {
		return nil, err
	}

	orderData.OrderDetails = orderDetails

	return dtos.NewOrderListDetailsGetByIDResponse(orderData), nil
}

func (s *orderService) UpdateProcessDelivery(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepository.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	if merchant.ID != order.MerchantID {
		return nil, e.ErrForbiddenMerchant
	}

	if order.Status != constants.OnProcessStatus {
		return nil, e.ErrStatusIsNotOnProcess
	}

	order.Status = constants.OutForDeliveryStatus

	updatedOrder, err := s.orderRepository.Update(orderID, order)

	if err != nil {
		return nil, err
	}

	orderTracking := []*models.OrderTracking{
		{
			OrderID:         orderID,
			HistoryDatetime: time.Now(),
			Title:           constants.CronOrderTrackingTitleSentFromOriginWarehouse,
			Description:     constants.CronOrderTrackingDescriptionReceivedOnDestinationWarehouse,
		},
	}

	_, err = s.orderTrackingRepository.Create(orderTracking)

	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (s *orderService) UpdateProcessOrder(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepository.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	if merchant.ID != order.MerchantID {
		return nil, e.ErrForbiddenMerchant
	}

	if order.Status != constants.PendingSellerStatus {
		return nil, e.ErrStatusIsNotPendingSeller
	}

	order.Status = constants.OnProcessStatus

	updatedOrder, err := s.orderRepository.Update(orderID, order)

	if err != nil {
		return nil, err
	}

	orderTracking := []*models.OrderTracking{
		{
			OrderID:         orderID,
			HistoryDatetime: time.Now(),
			Title:           constants.OrderProcessedTitle,
			Description:     constants.OrderProcessedDescription,
		},
	}

	_, err = s.orderTrackingRepository.Create(orderTracking)

	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (s *orderService) RefundCanceledOrder(orderID uuid.UUID) error {
	orderDetails, err := s.orderDetailRepository.OrderDetailListGetByOrderID(orderID)

	if err != nil {
		return err
	}

	timeNow := time.Now()

	refund, err := s.refundRepository.Create(&models.Refund{
		OrderID:               orderID,
		RefundRequestDatetime: timeNow,
		Status:                constants.RefundStatusCompleted,
		CompletedDatetime:     &timeNow,
	})

	if err != nil {
		return err
	}

	var refundDetailList []*models.RefundDetail

	var sumUserRefundAmount float64

	for _, od := range orderDetails {
		refundDetail, err := s.refundDetailRepository.CalculateRefundDetails(od.OrderDetailID, od.Quantity)

		if err != nil {
			return err
		}

		refundDetailList = append(refundDetailList, &models.RefundDetail{
			RefundID:           refund.ID,
			OrderDetailID:      od.OrderDetailID,
			Quantity:           od.Quantity,
			AmountUserGet:      math.Round(refundDetail.AmountUserGet),
			AmountMerchantLost: math.Round(refundDetail.AmountMerchantLost),
		})

		sumUserRefundAmount += math.Round(refundDetail.AmountUserGet)
	}

	order, err := s.orderRepository.GetByID(orderID)

	if err != nil {
		return err
	}

	sumUserRefundAmount += float64(order.ShippingFee)

	_, err = s.refundDetailRepository.Create(refundDetailList)

	if err != nil {
		return err
	}

	transaction, err := s.transactionRepository.GetByOrderID(orderID)

	if err != nil {
		return err
	}

	if transaction.PaymentMethodID == constants.WalletID {
		walletTransactionPayment, err := s.walletTransactionPaymentRepository.GetByTransactionID(transaction.ID)

		if err != nil {
			return err
		}

		newWallet := &models.Wallet{
			Balance: walletTransactionPayment.Wallet.Balance + uint(sumUserRefundAmount),
		}

		_, err = s.walletRepository.Update(walletTransactionPayment.WalletID, newWallet)

		if err != nil {
			return err
		}

	} else if transaction.PaymentMethodID == constants.SealabsPayID {
		sealabsPayTransactionPayment, err := s.sealabsPayTransactionPaymentRepository.GetByTransactionID(transaction.ID)

		if err != nil {
			return err
		}

		txIDInt, err := strconv.Atoi(sealabsPayTransactionPayment.TxID)

		if err != nil {
			return err
		}

		err = s.sealabsPayService.Refund(&dtos.SealabsPayRefundRequest{
			TxnID:  txIDInt,
			Amount: int(sumUserRefundAmount),
			Reason: "",
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *orderService) CronPaymentExpiration() error {
	pendingPayment, err := s.orderRepository.PendingPaymentGetAll()

	if err != nil {
		return err
	}

	newOrderTrackings := []*models.OrderTracking{}

	for _, v := range pendingPayment {
		interval, err := strconv.Atoi(os.Getenv("CRON_PAYMENT_EXPIRATION_INTERVAL"))
		if err != nil {
			interval = constants.CronPaymentExpirationInterval
		}

		if !time.Now().After(v.Transaction.TransactionDatetime.Add(time.Minute * time.Duration(interval))) {
			continue
		}

		_, err = s.orderRepository.Update(v.ID, &models.Order{
			Status: constants.SystemCanceledStatus,
		})

		if err != nil {
			return err
		}

		newOrderTrackings = append(newOrderTrackings, &models.OrderTracking{
			OrderID:         v.ID,
			HistoryDatetime: time.Now(),
			Title:           constants.CanceledBySystemPaymentTitle,
			Description:     constants.CanceledBySystemPaymentDescription,
		})
	}

	if len(newOrderTrackings) != 0 {
		_, err = s.orderTrackingRepository.Create(newOrderTrackings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *orderService) CronPendingSellerExpiration() error {
	orderTrackings, err := s.orderTrackingRepository.GetByOrderStatus(constants.PendingSellerStatus)

	if err != nil {
		return err
	}

	var newOrderTrackings []*models.OrderTracking

	for _, v := range orderTrackings {
		interval, err := strconv.Atoi(os.Getenv("CRON_SELLER_CONFIRMATION_EXPIRATION_INTERVAL"))
		if err != nil {
			interval = constants.CronSellerConfirmationExpirationInterval
		}

		if !time.Now().After(v.HistoryDatetime.Add(time.Minute * time.Duration(interval))) {
			continue
		}

		_, err = s.orderRepository.Update(v.OrderID, &models.Order{
			Status: constants.SystemCanceledStatus,
		})

		if err != nil {
			return err
		}

		newOrderTrackings = append(newOrderTrackings, &models.OrderTracking{
			OrderID:         v.OrderID,
			HistoryDatetime: time.Now(),
			Title:           constants.CanceledBySystemPendingSellerTitle,
			Description:     constants.CanceledBySystemPendingSellerDescription,
		})

		err = s.RefundCanceledOrder(v.OrderID)

		if err != nil {
			return err
		}
	}

	if len(newOrderTrackings) != 0 {
		_, err = s.orderTrackingRepository.Create(newOrderTrackings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *orderService) CronOrderProcessExpiration() error {
	orderTrackings, err := s.orderTrackingRepository.GetByOrderStatus(constants.OnProcessStatus)

	if err != nil {
		return err
	}

	var newOrderTrackings []*models.OrderTracking

	for _, v := range orderTrackings {
		maxProcessedDays, err := s.orderRepository.GetMaxProcessedDaysByOrderID(v.OrderID)

		if err != nil {
			return err
		}

		if !time.Now().After(v.HistoryDatetime.Add(time.Minute * time.Duration(60*24*maxProcessedDays))) {
			continue
		}

		_, err = s.orderRepository.Update(v.OrderID, &models.Order{
			Status: constants.SystemCanceledStatus,
		})

		if err != nil {
			return err
		}

		newOrderTrackings = append(newOrderTrackings, &models.OrderTracking{
			OrderID:         v.OrderID,
			HistoryDatetime: time.Now(),
			Title:           constants.CanceledBySystemOrderProcessTitle,
			Description:     constants.CanceledBySystemOrderProcessDescription,
		})

		err = s.RefundCanceledOrder(v.OrderID)

		if err != nil {
			return err
		}
	}

	if len(newOrderTrackings) != 0 {
		_, err = s.orderTrackingRepository.Create(newOrderTrackings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *orderService) CronAutoCompleteDelivered() error {
	orderTrackings, err := s.orderTrackingRepository.GetByOrderStatus(constants.DeliveredStatus)

	if err != nil {
		return err
	}

	var newOrderTrackings []*models.OrderTracking

	for _, v := range orderTrackings {
		interval, err := strconv.Atoi(os.Getenv("CRON_AUTO_COMPLETE_DELIVERED_EXPIRATION_INTERVAL"))
		if err != nil {
			interval = constants.CronAutoCompleteDeliveredExpirationInterval
		}

		if !time.Now().After(v.HistoryDatetime.Add(time.Minute * time.Duration(interval))) {
			continue
		}

		_, err = s.orderRepository.Update(v.OrderID, &models.Order{
			Status: constants.CompletedStatus,
		})

		if err != nil {
			return err
		}

		now := time.Now()

		newOrderTrackings = append(newOrderTrackings, &models.OrderTracking{
			OrderID:         v.OrderID,
			HistoryDatetime: now,
			Title:           constants.ConfirmedTitle,
			Description:     constants.ConfirmedDescription,
		})

		newOrderTrackings = append(newOrderTrackings, &models.OrderTracking{
			OrderID:         v.OrderID,
			HistoryDatetime: now,
			Title:           constants.CompletedTitle,
			Description:     constants.CompletedDescription,
		})
	}

	if len(newOrderTrackings) != 0 {
		_, err = s.orderTrackingRepository.Create(newOrderTrackings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *orderService) SellerCancelOrder(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)

	if err != nil {
		return nil, err
	}

	order, err := s.orderRepository.GetByID(orderID)

	if err != nil {
		return nil, err
	}

	if order.MerchantID != merchant.ID {
		return nil, e.ErrForbiddenMerchant
	}

	if order.Status != constants.PendingSellerStatus {
		return nil, e.ErrStatusIsNotPendingSeller
	}

	order.Status = constants.SellerCanceledStatus

	updatedOrder, err := s.orderRepository.Update(orderID, order)

	if err != nil {
		return nil, err
	}

	orderTracking := []*models.OrderTracking{
		{
			OrderID:         orderID,
			HistoryDatetime: time.Now(),
			Title:           constants.CanceledBySellerTitle,
			Description:     constants.CanceledBySellerDescription,
		},
	}

	_, err = s.orderTrackingRepository.Create(orderTracking)

	if err != nil {
		return nil, err
	}

	err = s.RefundCanceledOrder(orderID)

	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (s *orderService) UpdateCompleteOrder(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepository.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	if userID != order.Transaction.UserID {
		return nil, e.ErrForbiddenOrder
	}

	if order.Status != constants.DeliveredStatus {
		return nil, e.ErrStatusIsNotDelivered
	}

	order.Status = constants.CompletedStatus

	updatedOrder, err := s.orderRepository.Update(orderID, order)

	if err != nil {
		return nil, err
	}

	orderTracking := []*models.OrderTracking{
		{
			OrderID:         orderID,
			HistoryDatetime: time.Now(),
			Title:           constants.ConfirmedTitle,
			Description:     constants.ConfirmedDescription,
		},
		{
			OrderID:         orderID,
			HistoryDatetime: time.Now(),
			Title:           constants.CompletedTitle,
			Description:     constants.CompletedDescription,
		},
	}

	_, err = s.orderTrackingRepository.Create(orderTracking)

	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (s *orderService) SellerClaimFunds(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)

	if err != nil {
		return nil, err
	}

	order, err := s.orderRepository.GetByID(orderID)

	if err != nil {
		return nil, err
	}

	if order.MerchantID != merchant.ID {
		return nil, e.ErrForbiddenMerchant
	}

	if order.Status != constants.CompletedStatus {
		return nil, e.ErrStatusIsNotCompleted
	}

	if order.ClaimStatus {
		return nil, e.ErrOrderFundsAlreadyClaimed
	}

	order.ClaimStatus = true

	timeNow := time.Now()
	order.ClaimDatetime = &timeNow

	updatedOrder, err := s.orderRepository.Update(orderID, order)

	if err != nil {
		return nil, err
	}

	wallet, err := s.walletRepository.GetByUserId(userID)

	if err != nil {
		return nil, err
	}

	completedRefundDetails, err := s.refundRepository.GetDetailsByOrderIDAndStatus(orderID, constants.RefundStatusCompleted)

	if err != nil {
		return nil, err
	}

	var totalRefundMerchantLost float64 = 0
	for _, refund := range completedRefundDetails {
		totalRefundMerchantLost += refund.TotalMerchantLostAmount
	}

	totalOrderFunds := order.TotalPriceAfterDiscount - totalRefundMerchantLost

	newWallet := &models.Wallet{
		Balance: wallet.Balance + uint(totalOrderFunds),
	}

	_, err = s.walletRepository.Update(wallet.ID, newWallet)

	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}
