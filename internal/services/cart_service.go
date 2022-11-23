package services

import (
	in_cart_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/cart_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type CartService interface {
	Update(request *in_cart_dtos.CartRequest) (*models.Cart, error)
}

type cartService struct {
	cartRepository r.CartRepository
}

type CartServiceConfig struct {
	CartRepository r.CartRepository
}

func NewCartService(c *CartServiceConfig) CartService {
	return &cartService{
		cartRepository: c.CartRepository,
	}
}

func (p *cartService) Update(request *in_cart_dtos.CartRequest) (result *models.Cart, err error) {
	cart, err := p.cartRepository.GetByUserIdProductItemId(request.UserID, request.ProductItemID)
	if err != nil && err != e.ErrDataNotFound {
		return
	}

	if *request.Quantity == 0 {
		return p.cartRepository.Delete(cart.ID)
	}

	if err == e.ErrDataNotFound {
		return p.cartRepository.Create(request.ToModel())
	}

	cart.Quantity = *request.Quantity

	return p.cartRepository.Update(cart.ID, cart)
}
