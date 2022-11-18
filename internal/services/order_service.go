package services

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"

	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/google/uuid"
)

type OrderService interface {
	OrderListGetAll(id uuid.UUID, query *dtos.OrderListQuery) (*dtos.BuyerOrderListGetAllResponse, error)
}

type orderService struct {
	orderRepository       r.OrderRepository
	orderDetailRepository r.OrderDetailRepository
}

type OrderServiceConfig struct {
	OrderRepository       r.OrderRepository
	OrderDetailRepository r.OrderDetailRepository
}

func NewOrderService(c *OrderServiceConfig) OrderService {
	return &orderService{
		orderRepository:       c.OrderRepository,
		orderDetailRepository: c.OrderDetailRepository,
	}
}

func (s *orderService) OrderListGetAll(id uuid.UUID, query *dtos.OrderListQuery) (*dtos.BuyerOrderListGetAllResponse, error) {
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

	return dtos.NewBuyerOrderListGetAllResponse(orderList, pagination), nil
}
