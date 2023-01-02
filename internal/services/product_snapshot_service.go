package services

import (
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type ProductSnapshotService interface {
	GetByOrderDetailID(userID uuid.UUID, orderDetailID uuid.UUID) (*externalDtos.GetByOrderDetailIDResponse, error)
}

type productSnapshotService struct {
	productSnapshotRepository r.ProductSnapshotRepository
}

type ProductSnapshotServiceConfig struct {
	ProductSnapshotRepository r.ProductSnapshotRepository
}

func NewProductSnapshotService(c *ProductSnapshotServiceConfig) ProductSnapshotService {
	return &productSnapshotService{
		productSnapshotRepository: c.ProductSnapshotRepository,
	}
}

func (s *productSnapshotService) GetByOrderDetailID(userID uuid.UUID, orderDetailID uuid.UUID) (*externalDtos.GetByOrderDetailIDResponse, error) {
	productSnapshot, err := s.productSnapshotRepository.GetByOrderDetailID(orderDetailID)

	if err != nil {
		return nil, err
	}

	if productSnapshot.UserID != userID {
		return nil, e.ErrForbiddenOrder
	}

	return externalDtos.NewGetByOrderDetailIDResponse(productSnapshot), nil
}
