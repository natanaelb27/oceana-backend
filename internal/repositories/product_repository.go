package repositories

import (
	"database/sql"

	"fmt"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/product_media_constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"github.com/google/uuid"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos/product_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductRepository interface {
	ProductDetailsGetByID(productID uint) (*dtos.ProductDetails, error)
	GetPurchaseQuantitiesByID(userID uuid.UUID, productID uint, productItemID uint, quantity uint) (*dtos.SelectedProductPurchaseQuantities, error)
	GetAll(query *product_dtos.ProductQuery, pagination *utils.Pagination) ([]*product_dtos.Product, *utils.Pagination, error)
	CreateProduct(product *models.Product) (model *models.Product, err error)
	GetByProductItemID(productItemID uint) (*models.Product, error)
	UpdateUsingMapByMapCondition(conditionMap map[string]any, modelMap map[string]any) (err error)
	GetUpdateByProductID(id uint) (model *product_dtos.ProductModel, err error)
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

func (r *productRepository) ProductDetailsGetByID(productID uint) (*dtos.ProductDetails, error) {
	var productDetails *dtos.ProductDetails

	subQueryProductDiscount := r.db.Model(&models.ProductMerchantPromotion{}).
		Select(`
			product_merchant_promotions.product_id,
			product_merchant_promotions.discount_type,
			product_merchant_promotions.discount,
			product_merchant_promotions.maximum_discount_amount
		`).
		Joins("JOIN merchant_promotions ON merchant_promotions.id = product_merchant_promotions.merchant_promotion_id").
		Where(`
			product_merchant_promotions.maximum_quota > 
			CASE
				WHEN product_merchant_promotions.maximum_quota > -1 THEN product_merchant_promotions.current_quota
				ELSE -2
			END
		`).
		Where("NOW() BETWEEN merchant_promotions.start_datetime AND merchant_promotions.end_datetime").
		Where("product_merchant_promotions.product_id = ?", productID)

	subQueryProductItem := r.db.Model(&models.ProductItem{}).
		Select(`
			product_id AS product_id,
			SUM(stock_quantity) AS total_stock,
			MIN(price) AS price_min,
			MAX(price) AS price_max
		`).
		Where("product_items.product_id = ?", productID).
		Group("product_items.product_id")

	// subQueryProductReview := r.db.Model(&models.ProductItem{}).
	// 	Select(`
	// 		product_items.product_id AS product_id,
	// 		COUNT(product_reviews.id) AS rating_count,
	// 		AVG(product_reviews.rating) AS rating_star
	// 	`).
	// 	Joins("JOIN product_snapshots ON product_items.product_id = product_snapshots.product_id").
	// 	Joins("JOIN order_details ON product_snapshots.id = order_details.product_snapshot_id").
	// 	Joins("JOIN product_reviews ON order_details.id = product_reviews.order_detail_id").
	// 	Where("product_items.product_id = ?", productID).
	// 	Group("product_items.product_id")

	subQueryProductFavorite := r.db.Model(&models.ProductFavorite{}).
		Select(`
			product_favorites.product_id AS product_id,
			COUNT(*) AS like_count
		`).
		Where("product_favorites.product_id = ?", productID).
		Where("is_favorite = TRUE").
		Group("product_favorites.product_id")

	subQueryProductImage := r.db.Model(&models.ProductMedia{}).
		Select(`
			product_id,
			ARRAY_AGG(media_url) AS images
		`).
		Where("product_id = ?", productID).
		Where("type = ?", product_media_constants.TypeImage).
		Group("product_id")

	subQueryProductVideo := r.db.Model(&models.ProductMedia{}).
		Select(`
			product_id,
			ARRAY_AGG(media_url) AS videos
		`).
		Where("product_id = ?", productID).
		Where("type = ?", product_media_constants.TypeVideo).
		Group("product_id")

	queryResult := r.db.Model(&models.Product{}).
		Select(`
			products.name AS name,
			products.merchant_id AS merchant_id,
			products.is_hazardous AS is_hazardous,
			products.weight AS weight,
			products.package_size AS package_size,
			products.condition AS condition,
			products.category_id AS category_id,
			categories.name AS category,
			COALESCE(product_image_agg.images, '{}') AS images,
			COALESCE(product_video_agg.videos, '{}') AS videos,
			products.description AS description,
			product_item_agg.total_stock AS total_stock,
			product_favorite_agg.like_count AS like_count,
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
		Joins("LEFT JOIN (?) product_item_agg ON products.id = product_item_agg.product_id", subQueryProductItem).
		Joins("LEFT JOIN (?) product_discount_agg ON products.id = product_discount_agg.product_id", subQueryProductDiscount).
		Joins("LEFT JOIN (?) product_favorite_agg ON products.id = product_favorite_agg.product_id", subQueryProductFavorite).
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Joins("LEFT JOIN (?) product_image_agg ON products.id = product_image_agg.product_id", subQueryProductImage).
		Joins("LEFT JOIN (?) product_video_agg ON products.id = product_video_agg.product_id", subQueryProductVideo).
		Where("products.id = ?", productID).
		Where("products.status = 'active'").
		First(&productDetails)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrProductNotFound
		}
		return nil, e.ErrDatabaseOperationFailed
	}

	return productDetails, nil
}

func (r *productRepository) GetPurchaseQuantitiesByID(userID uuid.UUID, productID uint, productItemID uint, quantity uint) (*dtos.SelectedProductPurchaseQuantities, error) {
	var selectedProductPurchaseQuantities *dtos.SelectedProductPurchaseQuantities

	subQueryUserCart := r.db.Model(&models.Cart{}).
		Select(`
			carts.user_id AS user_id,
			product_items.product_id AS product_id, 
			SUM(carts.quantity) AS total_quantity
		`).
		Joins("JOIN product_items ON carts.product_item_id = product_items.id").
		Where("carts.user_id = ?", userID).
		Where("product_items.product_id = ?", productID).
		Group("carts.user_id, product_items.product_id")

	subQueryUserMerchantPromotion := r.db.Model(&models.UserMerchantPromotion{}).
		Select(`
			user_merchant_promotions.product_merchant_promotion_id AS product_merchant_promotion_id,
			product_merchant_promotions.product_id AS product_id,
			user_merchant_promotions.usage AS user_merchant_promotion_usage
		`).
		Joins("JOIN product_merchant_promotions ON product_merchant_promotions.id = user_merchant_promotions.product_merchant_promotion_id").
		Where("user_merchant_promotions.user_id = ?", userID).
		Where("product_merchant_promotions.product_id = ?", productID)

	subQueryUserUsage := r.db.Model(&models.ProductMerchantPromotion{}).
		Select(`
			product_merchant_promotions.id AS product_merchant_promotion_id, 
			COALESCE(user_cart.total_quantity, 0) + COALESCE(user_merchant_promo.user_merchant_promotion_usage, 0) AS total_usage
		`).
		Joins("LEFT JOIN (?) user_cart ON product_merchant_promotions.product_id = user_cart.product_id", subQueryUserCart).
		Joins("LEFT JOIN (?) user_merchant_promo ON product_merchant_promotions.id = user_merchant_promo.product_merchant_promotion_id", subQueryUserMerchantPromotion).
		Where("product_merchant_promotions.product_id = ?", productID)

	subQueryProductDiscount := r.db.Model(&models.ProductMerchantPromotion{}).
		Select(`
			product_merchant_promotions.product_id,
			product_merchant_promotions.id AS product_merchant_promotion_id,
			product_merchant_promotions.discount_type,
			product_merchant_promotions.discount,
			product_merchant_promotions.maximum_discount_amount,
			product_merchant_promotions.maximum_quota,
			product_merchant_promotions.maximum_user_limit,
			product_merchant_promotions.current_quota,
			COALESCE(user_usage.total_usage, 0) AS total_usage
		`).
		Joins("JOIN merchant_promotions ON merchant_promotions.id = product_merchant_promotions.merchant_promotion_id").
		Joins("LEFT JOIN (?) user_usage ON user_usage.product_merchant_promotion_id = product_merchant_promotions.id", subQueryUserUsage).
		Where(`
			product_merchant_promotions.maximum_quota > 
			CASE
				WHEN product_merchant_promotions.maximum_quota > -1 THEN product_merchant_promotions.current_quota
				ELSE -2
			END
		`).
		Where("NOW() BETWEEN merchant_promotions.start_datetime AND merchant_promotions.end_datetime").
		Where(`
			CASE
				WHEN product_merchant_promotions.maximum_user_limit > -1 
				THEN COALESCE(user_usage.total_usage, 0) + @quantity <= product_merchant_promotions.maximum_user_limit
				ELSE TRUE
			END
		`,
			sql.Named("quantity", quantity)).
		Where(`
			CASE
				WHEN product_merchant_promotions.maximum_quota > -1 
				THEN @quantity <= product_merchant_promotions.maximum_quota - product_merchant_promotions.current_quota
				ELSE TRUE
			END
		`,
			sql.Named("quantity", quantity)).
		Where("product_merchant_promotions.product_id = ?", productID)

	subQueryProduct := r.db.Model(&models.ProductItem{}).
		Select(`
			product_items.product_id AS product_id,
			CASE 
				WHEN product_discount_agg.product_merchant_promotion_id IS NOT NULL
				THEN
					(CASE
						WHEN LOWER(product_discount_agg.discount_type) = 'nominal' 
						THEN GREATEST(product_items.price - product_discount_agg.discount, 0)
						ELSE product_items.price - LEAST((product_items.price * product_discount_agg.discount) / 100.0, product_discount_agg.maximum_discount_amount)
					END)
				ELSE product_items.price
			END AS price,
			product_items.stock_quantity AS stock
		`).
		Joins("LEFT JOIN (?) product_discount_agg ON product_items.product_id = product_discount_agg.product_id", subQueryProductDiscount).
		Where("product_items.id = ?", productItemID)

	queryResult := r.db.Model(&models.Product{}).
		Select(`
			selected_product.price AS price,
			selected_product.stock AS stock,
			(
				CASE
					WHEN product_discount_agg.product_merchant_promotion_id IS NULL
					THEN selected_product.stock
					ELSE 
						(CASE
							WHEN product_discount_agg.maximum_quota <= -1 THEN selected_product.stock
							ELSE product_discount_agg.maximum_quota - product_discount_agg.current_quota
						END)
				END
			) AS max_purchase_quantity,
			(
				CASE
					WHEN product_discount_agg.product_merchant_promotion_id IS NULL
					THEN selected_product.stock
					ELSE 
						(CASE
							WHEN product_discount_agg.maximum_user_limit <= -1 THEN selected_product.stock
							ELSE product_discount_agg.maximum_user_limit - product_discount_agg.total_usage
						END)
				END
			) AS user_max_purchase_quantity,
			product_discount_agg.product_merchant_promotion_id AS product_merchant_promotion_id
		`).
		Joins("JOIN (?) selected_product ON products.id = selected_product.product_id", subQueryProduct).
		Joins("LEFT JOIN (?) product_discount_agg ON products.id = product_discount_agg.product_id", subQueryProductDiscount).
		Where("products.id = ?", productID).
		Where("products.status = 'active'").
		First(&selectedProductPurchaseQuantities)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrProductNotFound
		}
		return nil, e.ErrDatabaseOperationFailed
	}
	return selectedProductPurchaseQuantities, nil
}

func (r *productRepository) subQueryProductItemsMinimumMaximumPrice(priceType string) *gorm.DB {
	minMax := "MIN"
	if priceType == "max" {
		minMax = "MAX"
	}
	query := r.db.Model(&models.ProductItem{}).
		Select(minMax + "(product_items.price)").
		Where("product_items.product_id = products.id")
	return query
}

func (r *productRepository) subQueryProductViewsSum() *gorm.DB {
	query := r.db.Model(&models.ProductView{}).
		Select("SUM(product_views.count)").
		Where("product_views.product_id = products.id")
	return query
}

func (r *productRepository) subQueryProductPurchaseSum() *gorm.DB {
	query := r.db.Model(&models.ProductSnapshot{}).
		Select(`
			COALESCE(SUM(order_details.quantity),0)
		`).
		Joins("JOIN order_details ON order_details.product_snapshot_id = product_snapshots.id").
		Where("product_snapshots.product_id = products.id")
		// subtract by refund quantity?
		// condition where order details is completed only?
	return query
}

func (r *productRepository) subQueryProductCourierService(courierServices []string) *gorm.DB {
	query := r.db.Model(&models.ProductCourierService{}).
		Select("product_courier_services.courier_service_id").
		Where("product_courier_services.product_id = products.id").
		Where("product_courier_services.courier_service_id in (?)", courierServices)
	return query
}

func (r *productRepository) subQueryProductMediaByProductID() *gorm.DB {
	query := r.db.Model(&models.ProductMedia{}).
		Select("product_media.media_url").
		Where("product_media.product_id = products.id").
		Order("product_media.id ASC").
		Limit(1)
	return query
}

func (r *productRepository) GetAll(queryParams *product_dtos.ProductQuery, pagination *utils.Pagination) ([]*product_dtos.Product, *utils.Pagination, error) {
	sqPriceMinimum := sql.Named("price_minimum", r.subQueryProductItemsMinimumMaximumPrice("min"))
	sqPriceMaximum := sql.Named("price_maximum", r.subQueryProductItemsMinimumMaximumPrice("max"))
	sqTotalViewCount := sql.Named("total_view_count", r.subQueryProductViewsSum())
	sqTotalPurchaseSum := sql.Named("total_purchase_count", r.subQueryProductPurchaseSum())
	sqMediaURL := sql.Named("media_url", r.subQueryProductMediaByProductID())

	var ProductDetails []*product_dtos.Product

	subQueryProductDiscount := r.db.Model(&models.ProductMerchantPromotion{}).
		Select(`
			product_merchant_promotions.id,
			merchant_promotions.name,
			product_merchant_promotions.product_id,
			product_merchant_promotions.discount_type,
			product_merchant_promotions.discount,
			product_merchant_promotions.maximum_discount_amount
		`).
		Joins("JOIN merchant_promotions ON merchant_promotions.id = product_merchant_promotions.merchant_promotion_id").
		Where(`
			product_merchant_promotions.maximum_quota > 
			CASE
				WHEN product_merchant_promotions.maximum_quota > -1 THEN product_merchant_promotions.current_quota
				ELSE -2
			END
		`).
		Where("NOW() BETWEEN merchant_promotions.start_datetime AND merchant_promotions.end_datetime")

	subQueryProductReview := r.db.Model(&models.ProductItem{}).
		Select(`
			product_items.product_id AS product_id,
			AVG(product_reviews.rating) AS rating_star
		`).
		Joins("JOIN product_snapshots ON product_items.product_id = product_snapshots.product_id").
		Joins("JOIN order_details ON product_snapshots.id = order_details.product_snapshot_id").
		Joins("JOIN product_reviews ON order_details.id = product_reviews.order_detail_id").
		Group("product_items.product_id")

	query := r.db.Model(&models.Product{}).
		Select(`
			products.id,
            products.merchant_id,
            products.category_id,
            products.name,
            products.description,
            (@media_url) AS media_url,
            products.is_hazardous,
            products.weight,
            products.package_size,
            products.condition,
            products.status AS product_status,
			(@price_minimum) AS price_minimum,
			(@price_maximum) AS price_maximum,
			-- count price minimum after discount
			(
				CASE
					WHEN pmp.discount_type = 'percentage' 
						THEN (@price_minimum) - LEAST
							(
								(pmp.discount / 100.0) * (@price_minimum), 
								pmp.maximum_discount_amount
							)
					WHEN pmp.discount_type = 'nominal' 
						THEN GREATEST
							(
								(@price_minimum) - pmp.discount, 
								0
							)
				ELSE (@price_minimum)
				END
			) AS price_minimum_after_discount,
			-- count price maximum after discount
			(
				CASE
					WHEN pmp.discount_type = 'percentage' 
						THEN (@price_maximum) - LEAST
							(
								(pmp.discount / 100.0) * (@price_maximum), 
								pmp.maximum_discount_amount
							)
					WHEN pmp.discount_type = 'nominal' 
						THEN GREATEST
							(
								(@price_maximum) - pmp.discount, 
								0
							)
				ELSE (@price_maximum)
				END
			) AS price_maximum_after_discount,
			(@total_view_count) AS total_view_count,
			(@total_purchase_count) AS total_purchase_count,
			COALESCE(product_review.rating_star, 0) AS rating_star,
			pmp.id as product_merchant_promotion_id,
			pmp.name as merchant_promotion_name,
			pmp.discount as discount,
			pmp.discount_type as discount_type,
			pmp.maximum_discount_amount as discount_maximum_discount_amount,
            products.created_at,
            products.updated_at,
			merchants.id as "Merchant__id",
			merchants.name as "Merchant__name",
			cities.name as "Merchant__city"
			`,
			sqPriceMinimum,
			sqPriceMaximum,
			sqTotalViewCount,
			sqTotalPurchaseSum,
			sqMediaURL,
		).
		Joins("JOIN merchants ON merchants.id = products.merchant_id").
		Joins("JOIN addresses ON addresses.id = merchants.address_id").
		Joins("JOIN sub_districts ON sub_districts.id = addresses.sub_district_id").
		Joins("JOIN districts ON districts.id = sub_districts.district_id").
		Joins("JOIN cities ON cities.id = districts.city_id").
		Joins("LEFT JOIN (?) product_review ON products.id = product_review.product_id", subQueryProductReview).
		Joins("LEFT JOIN (?) pmp ON products.id = pmp.product_id", subQueryProductDiscount).
		Where("products.status = 'active'")

	// filter by merchant
	if queryParams.MerchantID != 0 {
		query.Where("products.merchant_id = ?", queryParams.MerchantID)
	}

	// filter by price category
	if queryParams.CategoryID != 0 {
		query.Where("products.category_id IN ?", queryParams.Categories)
	}

	// filter by price minimum
	if queryParams.PriceMinimum != 0 {
		query.Where(`(@price_minimum) >= @query_price_minimum`,
			sqPriceMinimum,
			sql.Named("query_price_minimum", queryParams.PriceMinimum),
		)
	}
	// filter by price maximum
	if queryParams.PriceMaximum != 0 {
		query.Where(`(@price_maximum) <= @query_price_maximum`,
			sqPriceMaximum,
			sql.Named("query_price_maximum", queryParams.PriceMaximum),
		)
	}

	if queryParams.RatingStarMinimum > 0 {
		query.Where("COALESCE(product_review.rating_star, 0) >= ?", queryParams.RatingStarMinimum)
	}

	// filter by merchant cities
	if len(queryParams.MerchantCities) > 0 {
		query = query.Where("cities.id in (?)", queryParams.MerchantCities)
	}

	// filter by cities
	if len(queryParams.CourierServices) > 0 {
		query.Where(`EXISTS (?)`,
			r.subQueryProductCourierService(queryParams.CourierServices),
		)
	}

	// search by products name
	query.Where("products.name ILIKE ?", "%"+queryParams.Search+"%")

	query.
		Scopes(utils.Paginate(pagination, query)).
		Order(fmt.Sprintf("%s %s", queryParams.SortBy, queryParams.Sort)).
		Find(&ProductDetails)

	if query.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}

	return ProductDetails, pagination, nil
}

func (p *productRepository) CreateProduct(newModel *models.Product) (model *models.Product, err error) {
	queryResult := p.db.Create(&newModel)
	if queryResult.RowsAffected == 0 {
		err = e.ErrInsertFailed
		return
	}
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}

func (r *productRepository) GetByProductItemID(productItemID uint) (*models.Product, error) {
	var model *models.Product

	result := r.db.Joins("LEFT JOIN product_items ON id = product_items.product_id").
		First(&model, "product_items.id = ?", productItemID)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, e.ErrProductNotFound
	}

	return model, result.Error
}

func (r *productRepository) UpdateUsingMapByMapCondition(conditionMap map[string]any, modelMap map[string]any) (err error) {
	result := r.db.Model(&models.Product{}).Where(conditionMap).Updates(&modelMap)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *productRepository) GetUpdateByProductID(id uint) (productModel *product_dtos.ProductModel, err error) {
	sqProductMediaImage := r.db.Model(&models.ProductMedia{}).
		Select("ARRAY_AGG(media_url)").
		Where("product_media.type = ?", product_media_constants.TypeImage).
		Where("product_media.product_id = products.id")
	sqProductMediaVideo := r.db.Model(&models.ProductMedia{}).
		Select("ARRAY_AGG(media_url)").
		Where("product_media.type = ?", product_media_constants.TypeVideo).
		Where("product_media.product_id = products.id")
	result := r.db.
		Table("products").
		Select(`
			products.id AS id,
			products.name AS name,
			products.description AS description,
			products.category_id AS category_id,
			products.merchant_id AS merchant_id,
			products.weight AS weight,
			products.is_hazardous AS is_hazardous,
			products.package_size AS package_size,
			products.condition AS condition,
			products.status AS status,
			products.process_day AS process_day,
			products.preorder_day AS preorder_day,
			(@sqProductMediaImage) AS images,
			(@sqProductMediaVideo) AS videos,
			NULL AS product_items
		`,
			sql.Named("sqProductMediaImage", sqProductMediaImage),
			sql.Named("sqProductMediaVideo", sqProductMediaVideo),
		).
		Where("id = ?", id).
		First(&productModel)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if result.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}

	return
}
