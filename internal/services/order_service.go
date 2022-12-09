package services

import (
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
	OrderListGetAll(id uuid.UUID, query *dtos.OrderListQuery) (*dtos.OrderListGetAllResponse, *utils.Pagination, error)
	OrderListGetAllByMerchantID(userID uuid.UUID, query *dtos.OrderListQuery) (*dtos.OrderListGetAllResponse, *utils.Pagination, error)
	OrderListDetailsGetByID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error)
	OrderListDetailsGetByIDByMerchantID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error)
	Update(request *internalDtos.OrderProcessOrder) (*models.Order, error)
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

func (s *orderService) OrderListGetAll(id uuid.UUID, query *dtos.OrderListQuery) (*dtos.OrderListGetAllResponse, *utils.Pagination, error) {
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
			TransactionID:       order.TransactionID,
			OrderID:             order.OrderID,
			MerchantID:          order.MerchantID,
			InvoiceNumber:       order.InvoiceNumber,
			ShopName:            order.ShopName,
			OrderDetails:        orderDetails,
			TransactionDatetime: order.TransactionDatetime,
			OrderStatus:         order.OrderStatus,
			TotalPrice:          order.TotalPrice,
		})

	}

	return dtos.NewOrderListGetAllResponse(orderList), pagination, nil
}

func (s *orderService) OrderListGetAllByMerchantID(userID uuid.UUID, query *dtos.OrderListQuery) (*dtos.OrderListGetAllResponse, *utils.Pagination, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	query.FilterMerchantID = merchant.ID

	return s.OrderListGetAll(userID, query)
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

func (s *orderService) OrderListDetailsGetByIDByMerchantID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error) {
	merchant, err := s.merchantRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	order, err := s.OrderListDetailsGetByID(userID, orderID)
	if err != nil {
		return nil, err
	}

	if merchant.ID != order.MerchantID {
		return nil, e.ErrUnauthorized
	}

	return order, err
}

func (s *orderService) Update(request *internalDtos.OrderProcessOrder) (*models.Order, error) {
	order, err := s.orderRepository.GetByID(request.OrderID)
	if err != nil {
		return nil, err
	}

	if order.Status != constants.PendingSellerStatus {
		return nil, e.ErrStatusIsNotPendingSeller
	}

	order.CourierServiceID = request.CourierServiceID
	order.Status = constants.OutForDeliveryStatus

	return s.orderRepository.Update(request.OrderID, order)
}
