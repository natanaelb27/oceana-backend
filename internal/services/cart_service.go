package services

import (
	"fmt"
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	in_cart_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/cart_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartService interface {
	Update(request *in_cart_dtos.CartRequest) (*models.Cart, error)
	GetAllByUserId(userId uuid.UUID) ([]*in_cart_dtos.Merchant, error)
	PayCartWithWallet(c *gin.Context, request *dtos.PayCartWithWalletRequest) (interface{}, error)
}

type cartService struct {
	cartRepository                     r.CartRepository
	productRepository                  r.ProductRepository
	productItemRepository              r.ProductItemRepository
	walletRepository                   r.WalletRepository
	productMerchantPromotionRepository r.ProductMerchantPromotionRepository
	merchantVoucherRepository          r.MerchantVoucherRepository
	marketplaceVoucherRepository       r.MarketplaceVoucherRepository
	redisRepository                    r.RedisRepository
}

type CartServiceConfig struct {
	CartRepository                     r.CartRepository
	ProductRepository                  r.ProductRepository
	ProductItemRepository              r.ProductItemRepository
	WalletRepository                   r.WalletRepository
	ProductMerchantPromotionRepository r.ProductMerchantPromotionRepository
	MerchantVoucherRepository          r.MerchantVoucherRepository
	MarketplaceVoucherRepository       r.MarketplaceVoucherRepository
	RedisRepository                    r.RedisRepository
}

func NewCartService(c *CartServiceConfig) CartService {
	return &cartService{
		cartRepository:                     c.CartRepository,
		productRepository:                  c.ProductRepository,
		productItemRepository:              c.ProductItemRepository,
		walletRepository:                   c.WalletRepository,
		productMerchantPromotionRepository: c.ProductMerchantPromotionRepository,
		merchantVoucherRepository:          c.MerchantVoucherRepository,
		marketplaceVoucherRepository:       c.MarketplaceVoucherRepository,
		redisRepository:                    c.RedisRepository,
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

func (p *cartService) PayCartWithWallet(c *gin.Context, request *dtos.PayCartWithWalletRequest) (interface{}, error) {
	redisKey := "wallet-pin-failed-count:" + request.UserID.String()
	redisValue, err := p.redisRepository.Get(c, redisKey)

	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}

	var failedAttemptCount int

	if err == e.ErrDataNotFound {
		failedAttemptCount = 0
	} else {
		value, err := strconv.Atoi(redisValue)

		if err != nil {
			failedAttemptCount = 0
		} else {
			failedAttemptCount = value
		}
	}

	if failedAttemptCount >= constants.MaxWalletFailAttempt {
		return nil, e.ErrBlockedWallet
	}

	wallet, err := p.walletRepository.GetByUserId(request.UserID)

	if err != nil {
		return nil, err
	}

	if err := utils.CheckHash(wallet.Pin, request.Pin); err != nil {
		if err := p.redisRepository.Set(c, redisKey, fmt.Sprintf("%d", failedAttemptCount+1), constants.BlockDuration); err != nil {
			return nil, err
		}

		return nil, e.ErrWrongWalletPin
	}

	// for _, cartOrderItem := range request.CartOrderItems {
	// 	merchantVoucher, err := p.merchantVoucherRepository.GetDetailsByID(*cartOrderItem.MerchantVoucherID)

	// 	if err != nil && err != e.ErrDataNotFound {
	// 		return nil, err
	// 	}

	// 	for _, orderCartItem := range cartOrderItem.OrderCartItems {
	// 		orderCartItemPrice := 0
	// 		if cartOrderItem.MerchantPromotionID != nil {
	// 			productItem, err := p.productItemRepository.GetById(orderCartItem.ProductItemID)

	// 			if err != nil {
	// 				return nil, err
	// 			}

	// 			productMerchantPromotion, err := p.productMerchantPromotionRepository.GetOnGoingByProductID(productItem.ProductID)

	// 			if err != nil && err != e.ErrDataNotFound {
	// 				return nil, err
	// 			}

	// 			if err == nil {
	// 				totalPrice := productItem.Price * *orderCartItem.Quantity
	// 				discount := uint(0)

	// 				if productMerchantPromotion.DiscountType == constants.NominalDiscountType {
	// 					discount = productMerchantPromotion.Discount
	// 				}

	// 				if productMerchantPromotion.DiscountType == constants.PercentageDiscountType {
	// 					discount = uint(float32(float32(totalPrice) * float32(discount/100)))

	// 					if discount > productMerchantPromotion.MaximumDiscountAmount {
	// 						discount = productMerchantPromotion.MaximumDiscountAmount
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// TO-DO: validate user wallet balance

	// TO-DO: validate user wallet balance

	// TO-DO: validate product quantity

	// TO-DO: create snapshot for bunch of stuffs

	// TO-DO: do payment stuffs

	return request, nil
}
