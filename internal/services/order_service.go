package services

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"

	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/google/uuid"
)

type OrderService interface {
	OrderListGetAll(id uuid.UUID, query *dtos.OrderListQuery) (*dtos.OrderListGetAllResponse, error)
	OrderListDetailsGetByID(userID uuid.UUID, orderID uuid.UUID) (*dtos.OrderListDetailsGetByIDResponse, error)
}

type orderService struct {
	orderRepository           r.OrderRepository
	orderDetailRepository     r.OrderDetailRepository
	addressSnapshotRepository r.AddressSnapshotRepository
	orderTrackingRepository   r.OrderTrackingRepository
}

type OrderServiceConfig struct {
	OrderRepository           r.OrderRepository
	OrderDetailRepository     r.OrderDetailRepository
	AddressSnapshotRepository r.AddressSnapshotRepository
	OrderTrackingRepository   r.OrderTrackingRepository
}

func NewOrderService(c *OrderServiceConfig) OrderService {
	return &orderService{
		orderRepository:           c.OrderRepository,
		orderDetailRepository:     c.OrderDetailRepository,
		addressSnapshotRepository: c.AddressSnapshotRepository,
		orderTrackingRepository:   c.OrderTrackingRepository,
	}
}

func (s *orderService) OrderListGetAll(id uuid.UUID, query *dtos.OrderListQuery) (*dtos.OrderListGetAllResponse, error) {
	pagination := &utils.Pagination{
		Limit: query.Limit,
		Page:  query.Page,
	}

	orderData, pagination, err := s.orderRepository.OrderListGetAll(id, query, pagination)

	if err != nil {
		return nil, err
	}

	orderList := []*internalDtos.OrderList{}

	for _, order := range orderData {
		orderDetails, err := s.orderDetailRepository.OrderDetailListGetByOrderID(order.OrderID)

		if err != nil {
			return nil, err
		}

		orderList = append(orderList, &internalDtos.OrderList{
			TransactionID:       order.TransactionID,
			OrderID:             order.OrderID,
			MerchantID:          order.MerchantID,
			ShopName:            order.ShopName,
			OrderDetails:        orderDetails,
			TransactionDatetime: order.TransactionDatetime,
			OrderStatus:         order.OrderStatus,
			TotalPrice:          order.TotalPrice,
		})

	}

	return dtos.NewOrderListGetAllResponse(orderList, pagination), nil
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
