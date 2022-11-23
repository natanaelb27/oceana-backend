package repositories

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	ProductDetailsGetByID(id uint, userID uuid.UUID) (*dtos.ProductDetails, error)
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

func (r *productRepository) ProductDetailsGetByID(id uint, userID uuid.UUID) (*dtos.ProductDetails, error) {
	var productDetails *dtos.ProductDetails

	subQueryUserMerchantPromotion := r.db.Model(&models.UserMerchantPromotion{}).
		Select(`
			usage
		`).
		Where("user_merchant_promotions.product_merchant_promotion_id = product_merchant_promotions.id").
		Where("user_merchant_promotions.user_id = ?", userID)

	subQueryProductDiscount := r.db.Model(&models.ProductMerchantPromotion{}).
		Select(`
			product_merchant_promotions.discount_type,
			product_merchant_promotions.discount,
			product_merchant_promotions.maximum_discount_amount
		`).
		Joins("JOIN merchant_promotions ON merchant_promotions.id = product_merchant_promotions.merchant_promotion_id").
		Where("product_merchant_promotions.status = 'active'").
		Where("product_merchant_promotions.current_quota < product_merchant_promotions.maximum_quota").
		Where("NOW() BETWEEN merchant_promotions.start_datetime AND merchant_promotions.end_datetime").
		Where("product_merchant_promotions.maximum_user_limit > COALESCE((?), 0)", subQueryUserMerchantPromotion).
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
					THEN GREATEST(product_item_agg.price_min - product_discount_agg.discount, 0)
					ELSE product_item_agg.price_min - LEAST((product_item_agg.price_min * product_discount_agg.discount) / 100.0, product_discount_agg.maximum_discount_amount)
				END
			) AS price_after_discount,
			product_item_agg.price_min AS price,
			(
				CASE 
					WHEN LOWER(product_discount_agg.discount_type) = 'nominal' 
					THEN GREATEST(product_item_agg.price_max - product_discount_agg.discount, 0)
					ELSE product_item_agg.price_max - LEAST((product_item_agg.price_max * product_discount_agg.discount) / 100.0, product_discount_agg.maximum_discount_amount)
				END
			) AS price_max_after_discount,
			product_item_agg.price_max AS price_max,
			(
				CASE 
					WHEN LOWER(product_discount_agg.discount_type) = 'nominal' 
					THEN GREATEST(product_item_agg.price_min - product_discount_agg.discount, 0)
					ELSE product_item_agg.price_min - LEAST((product_item_agg.price_min * product_discount_agg.discount) / 100.0, product_discount_agg.maximum_discount_amount)
				END
			) AS price_min_after_discount,
			product_item_agg.price_min AS price_min
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
