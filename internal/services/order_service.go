package services

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"

	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/google/uuid"
)

type OrderService interface {
	OrderListGetAll(id uuid.UUID, query *dtos.OrderListQuery) ([]*dtos.OrderListGetAllResponse, *utils.Pagination, error)
	OrderListGetAllByMerchantID(userID uuid.UUID, query *dtos.OrderListQuery) ([]*dtos.OrderListGetAllResponse, *utils.Pagination, error)
	OrderListDetailsGetByID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error)
	OrderListDetailsGetByMerchantID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error)
	UpdateProcessDelivery(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error)
	UpdateProcessOrder(orderID uuid.UUID, userID uuid.UUID) (*models.Order, error)
	CronPaymentExpiration() error
}

type orderService struct {
	orderRepository           r.OrderRepository
	orderDetailRepository     r.OrderDetailRepository
	addressSnapshotRepository r.AddressSnapshotRepository
	orderTrackingRepository   r.OrderTrackingRepository
	merchantRepository        r.MerchantRepository
}

type OrderServiceConfig struct {
	OrderRepository           r.OrderRepository
	OrderDetailRepository     r.OrderDetailRepository
	AddressSnapshotRepository r.AddressSnapshotRepository
	OrderTrackingRepository   r.OrderTrackingRepository
	MerchantRepository        r.MerchantRepository
}

func NewOrderService(c *OrderServiceConfig) OrderService {
	return &orderService{
		orderRepository:           c.OrderRepository,
		orderDetailRepository:     c.OrderDetailRepository,
		addressSnapshotRepository: c.AddressSnapshotRepository,
		orderTrackingRepository:   c.OrderTrackingRepository,
		merchantRepository:        c.MerchantRepository,
	}
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
			UserID:                               order.UserID,
			UserFullName:                         order.UserFullName,
			CourierServiceID:                     order.CourierServiceID,
			CourierDescription:                   order.CourierDescription,
			OrderDetails:                         orderDetails,
			TransactionDatetime:                  order.TransactionDatetime,
			OrderStatus:                          order.OrderStatus,
			TotalPrice:                           order.TotalPrice,
			RefundID:                             order.RefundID,
			ClaimStatus:                          order.ClaimStatus,
			SellerConfirmationExpirationDatetime: order.SellerConfirmationExpirationDatetime,
			OnProcessExpirationDatetime:          order.OnProcessExpirationDatetime,
			DeliveredExpirationDatetime:          order.DeliveredExpirationDatetime,
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
			UserID:                               order.UserID,
			UserFullName:                         order.UserFullName,
			CourierServiceID:                     order.CourierServiceID,
			CourierDescription:                   order.CourierDescription,
			OrderDetails:                         orderDetails,
			TransactionDatetime:                  order.TransactionDatetime,
			OrderStatus:                          order.OrderStatus,
			TotalPrice:                           order.TotalPrice,
			RefundID:                             order.RefundID,
			ClaimStatus:                          order.ClaimStatus,
			SellerConfirmationExpirationDatetime: order.SellerConfirmationExpirationDatetime,
			OnProcessExpirationDatetime:          order.OnProcessExpirationDatetime,
			DeliveredExpirationDatetime:          order.DeliveredExpirationDatetime,
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
			Title:           constants.OrderShippedTitle,
			Description:     constants.OrderShippedDescription,
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

func (s *orderService) CronPaymentExpiration() error {
	pendingPayment, err := s.orderRepository.PendingPaymentGetAll()

	if err != nil {
		return err
	}

	newOrderTrackings := []*models.OrderTracking{}

	for _, v := range pendingPayment {
		if !time.Now().After(v.Transaction.TransactionDatetime.Add(time.Minute * time.Duration(constants.PaymentExpirationInterval))) {
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
