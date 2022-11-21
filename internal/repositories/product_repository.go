package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	ProductDetailsGetByID(id uint, user_id uuid.UUID) (*dtos.ProductDetails, error)
}

type productRepository struct {
	db *gorm.DB
}

type ProductRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductRepository(c *ProductRepositoryConfig) ProductRepository {
	return &productRepository{
		db: c.DB,
	}
}

func (r *productRepository) ProductDetailsGetByID(id uint, user_id uuid.UUID) (*dtos.ProductDetails, error) {
	var productDetails *dtos.ProductDetails

	subQueryUserMerchantPromotion := r.db.Model(&models.OrderDetail{}).
		Select(`
			SUM()
		`).
		Joins("order_detail_merchant_promotions ON order_detail.id = order_detail_merchant_promotions.order_detail_id").
		Joins("product_merchant_promotions pmp ON order_detail_merchant_promotions.merchant_promotion_id = pmp.merchant_promotion_id").
		Joins("orders ON order_detail_id.order_id = orders.id").
		Joins("transactions ON orders.transaction_id = transactions.id").
		Where("transactions.user_id = ?", user_id).
		Where("pmp.id = product_merchant_promotions.id").
		Group("product_merchant_promotions.merchant_promotion_id")

		// Joins("refunds ON order_details.id = refunds.order_detail_id").
		// Where("LOWER(refunds.Status) = 'agreed'").

	// CEK ORDER DETAIL - REFUND
	subQueryProductDiscount := r.db.Model(&models.ProductMerchantPromotion{}).
		Select(`
			merchant_promotions.discount_type AS discount_type,
			merchant_promotions.discount AS discount,
			merchant_promotions.maximum_discount_amount AS maximum_discount_amount 
		`).
		Joins("JOIN merchant_promotions ON product_merchant_promotions.merchant_promotion_id = merchant_promotions.id").
		Where("product_merchant_promotions.status = 'active'").
		Where("product_merchant_promotions.current_quota < product_merchant_promotions.maximum_quota").
		Where("NOW() BETWEEN merchant_promotions.start_datetime AND merchant_promotions.end_datetime").
		Where("product_merchant_promotions.maximum_user_limit > ?", subQueryUserMerchantPromotion).
		Where("product_merchant_promotions.product_id = ?", id)

	subQueryProductItem := r.db.Model(&models.ProductItem{}).
		Select(`
			SUM(stock_quantity) AS total_stock,
			MIN(price) AS price_min,
			MAX(price) AS price_max
		`).
		Where("product_items.product_id = products.id").
		Group("product_items.product_id")

	subQueryProductReview := r.db.Model(&models.ProductItem{}).
		Select(`
		product_items.product_id AS product_id,
			COUNT(product_reviews.id) AS total_rating,
			AVG(product_reviews.rating) AS rating_star
		`).
		Joins("JOIN product_snapshots ON product_item.product_id = product_snapshots.product_id").
		Joins("JOIN order_details ON product_snapshots.id = order_details.product_snapshot_id").
		Joins("JOIN product_reviews ON order_details.id = product_reviews.order_detail_id").
		Where("product_items.product_id = ?", id).
		Group("product_items.product_id")

	queryResult := r.db.Model(&models.Product{}).
		Select(`
			products.name AS name,
			products.description AS description,
			product_item_agg.total_stock AS total_stock,
			(?) AS total_like,
			product_review_agg.total_rating AS total_rating,
			product_review_agg.rating_star AS rating_star,
			product_discount_agg.discount_type AS discount_type,
			product_discount_agg.discount AS discount,
			(
				CASE 
					WHEN LOWER(product_discount_agg.discount_type) = 'nominal' 
					THEN GREATEST(product_item_agg.price_min - product_discount_agg.maximum_discount_amount, 0)
					ELSE product_item_agg.price_min - LEAST((product_item_agg.price_min * product_discount_agg.discount)::float / 100, product_discount_agg.maximum_discount_amount)
				END
			) AS price,
			product_item_agg.price_min AS price_before_discount,
			(
				CASE 
					WHEN LOWER(product_discount_agg.discount_type) = 'nominal' 
					THEN GREATEST(product_item_agg.price_max - product_discount_agg.maximum_discount_amount, 0)
					ELSE product_item_agg.price_max - LEAST((product_item_agg.price_max * product_discount_agg.discount)::float / 100, product_discount_agg.maximum_discount_amount)
				END
			) AS price_max,
			product_item_agg.price_max AS price_max_before_discount,
			(
				CASE 
					WHEN LOWER(product_discount_agg.discount_type) = 'nominal' 
					THEN GREATEST(product_item_agg.price_min - product_discount_agg.maximum_discount_amount, 0)
					ELSE product_item_agg.price_min - LEAST((product_item_agg.price_min * product_discount_agg.discount)::float / 100, product_discount_agg.maximum_discount_amount)
				END
			) AS price_min,
			product_item_agg.price_min AS price_min_before_discount
		`).
		Joins("JOIN (?) product_item_agg ON products.id = product_item_agg.product_id", subQueryProductItem).
		Joins("JOIN (?) product_discount_agg ON products.id = product_discount_agg.product_id", subQueryProductDiscount).
		Joins("JOIN (?) product_review_agg ON products.id = product_review_agg.product_id", subQueryProductReview).
		Where("products.id = ?", id).
		First(&productDetails)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrDataNotFound
			// ini ganti jd product not found
		}
		return nil, e.ErrDatabaseOperationFailed
	}

	return productDetails, nil
}
