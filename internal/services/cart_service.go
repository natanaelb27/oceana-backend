package services

import (
	in_cart_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/cart_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
)

type CartService interface {
	Update(request *in_cart_dtos.CartRequest) (*models.Cart, error)
	GetAllByUserId(userId uuid.UUID) ([]*in_cart_dtos.Merchant, error)
}

type cartService struct {
	cartRepository        r.CartRepository
	productItemRepository r.ProductItemRepository
}

type CartServiceConfig struct {
	CartRepository        r.CartRepository
	ProductItemRepository r.ProductItemRepository
}

func NewCartService(c *CartServiceConfig) CartService {
	return &cartService{
		cartRepository:        c.CartRepository,
		productItemRepository: c.ProductItemRepository,
	}
}

func (p *cartService) Update(request *in_cart_dtos.CartRequest) (result *models.Cart, err error) {
	cart, cartErr := p.cartRepository.GetByUserIdProductItemId(request.UserID, request.ProductItemID)
	if cartErr != nil && cartErr != e.ErrDataNotFound {
		err = cartErr
		return
	}

	// check if stock available
	productItem, err := p.productItemRepository.GetById(request.ProductItemID)
	if err != nil {
		return
	}
	if productItem.StockQuantity < *request.Quantity {
		err = e.ErrStockIsNotAvailable
		return
	}

	if cartErr == e.ErrDataNotFound {
		if *request.Quantity == 0 {
			err = e.ErrEmptyQuantity
			return
		}
		return p.cartRepository.Create(request.ToModel())
	}

	if *request.Quantity == 0 && err != e.ErrDataNotFound {
		return p.cartRepository.Delete(cart.ID)
	}

	cart.Quantity = *request.Quantity

	return p.cartRepository.Update(cart.ID, cart)
}

func (p *cartService) GetAllByUserId(userId uuid.UUID) ([]*in_cart_dtos.Merchant, error) {
	cartsMerchants, err := p.cartRepository.GetCartsMerchants(userId)
	for k, v := range cartsMerchants {
		cartsMerchants[k].Products, err = p.cartRepository.GetCartsProductsByUserIdMerchants(v.MerchantID, userId)
	}
	return cartsMerchants, err
}
