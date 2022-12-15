package repositories

import (
	"database/sql"
	"log"

	in_cart_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/cart_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	WithTx(tx *gorm.DB) *cartRepository
	Create(cart *models.Cart) (model *models.Cart, err error)
	GetByUserIdProductItemId(userId uuid.UUID, productItemId uint) (model *models.Cart, err error)
	GetCartsMerchants(userId uuid.UUID) (result []*in_cart_dtos.Merchant, err error)
	GetCartsProductsByUserIdMerchants(merchantId uint, userId uuid.UUID) (result []*in_cart_dtos.Product, err error)
	Update(id uint, newModel *models.Cart) (model *models.Cart, err error)
	Delete(id uint) (model *models.Cart, err error)
}

type cartRepository struct {
	db *gorm.DB
}

type CartRepositoryConfig struct {
	DB *gorm.DB
}

func NewCartRepository(c *CartRepositoryConfig) CartRepository {
	return &cartRepository{
		db: c.DB,
	}
}

func (u cartRepository) WithTx(tx *gorm.DB) *cartRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *cartRepository) Create(newModel *models.Cart) (model *models.Cart, err error) {
	queryResult := r.db.Create(&newModel)
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

func (r *cartRepository) GetCartsProductsByUserIdMerchants(merchantId uint, userId uuid.UUID) (result []*in_cart_dtos.Product, err error) {
	sqTotalViewCount := sql.Named("sqTotalViewCount",
		r.db.Model(&models.Cart{}).
			Select("SUM(carts.quantity)").
			Joins(`JOIN product_items ON product_items.id = carts.product_item_id`).
			Joins(`JOIN products p ON p.id = product_items.product_id`).
			Where("carts.user_id = ?", userId).
			Where("merchants.id = ?", merchantId).
			Where("p.id = products.id"),
	)
	subQueryProductVariations := sql.Named("subQueryProductVariations",
		r.db.
			Table("product_variations pvariations").
			Select(`string_agg(CONCAT(pvarianttype."name", ':', pvariants."name"),', ' ORDER BY pvarianttype.id) AS value`).
			Joins(`JOIN product_variants pvariants ON pvariants.id = pvariations.product_variant_id `).
			Joins(`JOIN product_variant_types pvarianttype ON pvarianttype.id = pvariants.product_variant_type_id  `).
			Where("pvariations.product_item_id = carts.product_item_id").
			Group("pvariations.product_item_id"),
	)
	subQueryPromoValidation := sql.Named("subQueryPromoValidation",
		r.db.Raw(`
			CASE
				-- validate maximum user limit
				-- pmp.maximum_user_limit <> -1 means that if there is no user limit
				-- validate merchant promo user quota - ump.usage + carts.quantity > pmp.maximum_user_limit  means if user merchant promo usage + carts quantity then show original price
				WHEN pmp.maximum_user_limit <> -1 AND (@sqTotalViewCount) + COALESCE(ump.usage,0) > pmp.maximum_user_limit  
					THEN false
				-- validate merchant promo global quota
				WHEN pmp.maximum_quota <> -1 AND (@sqTotalViewCount) + pmp.current_quota > pmp.maximum_quota
					THEN false
				-- validate mp start and end date
				WHEN NOW() NOT BETWEEN mp.start_datetime AND mp.end_datetime
					THEN false
			ELSE true
			END		
		`,
			sqTotalViewCount,
		),
	)
	queryResult := r.db.Model(&models.Cart{}).Select(`
			carts.id AS carts_id,
			carts.user_id AS user_id,
			products.id AS product_id,
			products.name AS product_name,
			-- hide promo id when promo condition dont met
			(
				CASE 
					WHEN (@subQueryPromoValidation) = false
						THEN null
				ELSE pmp.merchant_promotion_id
				END
			) AS promo_id,
			carts.product_item_id AS product_item_id,
			carts.quantity AS quantity,
			product_items.stock_quantity AS stock_quantity,
			product_items.price AS price,
			-- select user limit
			(
				CASE 
					WHEN (@subQueryPromoValidation) = false
						THEN product_items.stock_quantity
				ELSE GREATEST(COALESCE(pmp.maximum_user_limit, 0) - COALESCE(ump.usage, 0), 0)
				END
			) AS promo_user_quota,
			-- maximum global quota
			(
				CASE 
					WHEN (@subQueryPromoValidation) = false
						THEN product_items.stock_quantity
				ELSE (
					GREATEST(
						COALESCE(
							pmp.maximum_quota
							, 0
						) - pmp.current_quota
						, 0
					)
				)
				END
			) AS promo_global_quota,
			-- count price after discount
			(
				CASE 
					WHEN (@subQueryPromoValidation) = false
						THEN product_items.price
					WHEN pmp.discount_type = 'percentage' 
						THEN product_items.price - LEAST
						(
							(pmp.discount / 100.0) * product_items.price, 
							pmp.maximum_discount_amount
						)
					-- price after nominal promo
					WHEN pmp.discount_type = 'nominal' 
						THEN GREATEST
						(
							product_items.price - pmp.discount, 
							0
						)
				ELSE product_items.price
				END
			) AS price_after_discount,
			(@subQueryProductVariations) AS product_variation
		`,
		sqTotalViewCount,
		subQueryProductVariations,
		subQueryPromoValidation,
	).
		Joins(`JOIN product_items ON product_items.id = carts.product_item_id`).
		Joins(`JOIN products ON products.id = product_items.product_id`).
		Joins("LEFT JOIN product_merchant_promotions pmp ON products.id = pmp.product_id AND pmp.status = ?", "active").
		Joins("LEFT JOIN user_merchant_promotions ump ON ump.product_merchant_promotion_id = pmp.id AND ump.user_id = carts.user_id").
		Joins("LEFT JOIN merchant_promotions mp ON mp.id = pmp.merchant_promotion_id").
		Joins(`JOIN merchants ON merchants.id = products.merchant_id`).
		Where("carts.user_id = ?", userId).
		Where("merchants.id = ?", merchantId).
		Order("carts.id DESC").
		Find(&result)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *cartRepository) GetCartsMerchants(userId uuid.UUID) (result []*in_cart_dtos.Merchant, err error) {
	// join products, product_items
	queryResult := r.db.Debug().Model(&models.Cart{}).Select(`
		merchants.id AS merchant_id,
		merchants.name AS name,
		provinces.name AS merchant_province,
		NULL as products
	`).
		Where("carts.user_id = ?", userId).
		Joins(`JOIN product_items ON product_items.id = carts.product_item_id`).
		Joins(`JOIN products ON products.id = product_items.product_id`).
		Joins(`JOIN merchants ON merchants.id = products.merchant_id`).
		Joins("JOIN addresses ON merchants.address_id = addresses.id").
		Joins("JOIN sub_districts ON addresses.sub_district_id = sub_districts.id").
		Joins("JOIN districts ON sub_districts.district_id = districts.id").
		Joins("JOIN cities ON districts.city_id = cities.id").
		Joins("JOIN provinces ON cities.province_id = provinces.id").
		Group("merchants.id, provinces.name").
		Find(&result)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}

func (r *cartRepository) GetByUserIdProductItemId(userId uuid.UUID, productItemId uint) (model *models.Cart, err error) {
	result := r.db.Where("user_id = ?", userId).Where("product_item_id = ?", productItemId).First(&model)
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

func (r *cartRepository) Update(id uint, newModel *models.Cart) (model *models.Cart, err error) {
	result := r.db.Model(&model).Where("id = ?", id).Updates(&newModel)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model = newModel
	model.ID = id

	return
}

func (r *cartRepository) Delete(id uint) (model *models.Cart, err error) {
	result := r.db.Delete(&model, id)
	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}
