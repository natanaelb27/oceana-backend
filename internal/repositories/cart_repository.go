package repositories

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	in_cart_dtos "github.com/natanaelb27/oceana-backend/internal/dtos/cart_dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	sqCartQuantityTotal := sql.Named("sqCartQuantityTotal",
		r.db.Model(&models.Cart{}).
			Select("SUM(carts.quantity)").
			Joins(`JOIN product_items ON product_items.id = carts.product_item_id`).
			Joins(`JOIN products p ON p.id = product_items.product_id`).
			Where("carts.user_id = ?", userId).
			Where("merchants.id = ?", merchantId).
			Where("p.id = products.id"),
	)
	sqProductVariantMedia := r.db.
		Table("product_variant_media pvm").
		Select(`
			pvm.id as product_variant_id,
			pvm.media_url as media_url
		`).
		Where("pvm.product_variant_id = pvariations.product_variant_id").
		Order("pvm.id ASC").
		Group("pvm.id, pvm.media_url").
		Limit(1)

	subQueryProductVariations := r.db.
		Table("product_variations pvariations").
		Select(`
			string_agg(pvariants."name",', ' ORDER BY pvarianttype.id) AS variations,
			pvariations.product_item_id as product_item_id,
			MIN(pvariants.id) AS product_variant_id
		`).
		Joins(`JOIN product_variants pvariants ON pvariants.id = pvariations.product_variant_id`).
		Joins(`JOIN product_variant_types pvarianttype ON pvarianttype.id = pvariants.product_variant_type_id `).
		Where("pvariations.product_item_id = carts.product_item_id").
		Group("pvariations.product_item_id")

	subQueryPromoValidation := sql.Named("subQueryPromoValidation",
		r.db.Raw(`
			CASE
				WHEN pmp.maximum_user_limit <> -1 AND (@sqCartQuantityTotal) + COALESCE(ump.usage,0) > pmp.maximum_user_limit  
					THEN false
				WHEN pmp.maximum_quota <> -1 AND (@sqCartQuantityTotal) + pmp.current_quota > pmp.maximum_quota
					THEN false
				WHEN NOW() NOT BETWEEN mp.start_datetime AND mp.end_datetime
					THEN false
			ELSE true
			END		
			`,
			sqCartQuantityTotal,
		),
	)
	sqPromo := r.db.Table("product_merchant_promotions pmp").
		Select(`ump.id as ump_id, pmp.id as pmp_id, mp.id as mp_id, *`).
		Joins("LEFT JOIN user_merchant_promotions ump ON ump.product_merchant_promotion_id = pmp.id AND ump.user_id = carts.user_id").
		Joins("LEFT JOIN merchant_promotions mp ON mp.id = pmp.merchant_promotion_id").
		Where("pmp.product_id = products.id").
		Where("(@subQueryPromoValidation) = true", subQueryPromoValidation).
		Order("pmp.id ASC").
		Limit(1)
	sqProductMedia := r.db.
		Table("product_media pm").
		Select(`
			pm.id,
			pm.product_id,
			pm.media_url
		`).
		Where("pm.type = 'image'").
		Where("pm.product_id = products.id").
		Order("pm.id ASC").
		Limit(1)
	queryResult := r.db.Model(&models.Cart{}).Select(`
			carts.id AS carts_id,
			carts.user_id AS user_id,
			products.id AS product_id,
			products.name AS product_name,
			products.weight AS product_weight,
			(
				CASE 
					WHEN pvm.media_url <> ''
						THEN pvm.media_url
					ELSE pm.media_url
				END
			) AS product_media_url,
			-- hide promo id when promo condition dont met
			promo.pmp_id AS product_merchant_promotion_id,
			carts.product_item_id AS product_item_id,
			carts.quantity AS quantity,
			product_items.stock_quantity AS stock_quantity,
			product_items.price AS price,
			-- select user limit
			(
				CASE 
					WHEN promo.pmp_id IS NULL
						THEN product_items.stock_quantity
					WHEN promo.maximum_user_limit = -1
						THEN product_items.stock_quantity
				ELSE GREATEST(
					COALESCE(
						promo.maximum_user_limit, 
						0
					) - COALESCE(
						promo.usage, 
						0
					), 
					0
				)
				END
			) AS promo_user_quota,
			-- maximum global quota
			(
				CASE 	
					WHEN promo.pmp_id IS NULL
						THEN product_items.stock_quantity
					WHEN promo.maximum_quota = -1
						THEN product_items.stock_quantity
				ELSE (
					GREATEST(
						COALESCE(
							promo.maximum_quota, 
							0
						) - promo.current_quota, 
						0
					)
				)
				END
			) AS promo_global_quota,
			-- count price after discount
			(
				CASE 
					WHEN promo.discount_type = 'percentage' 
						THEN product_items.price - LEAST
						(
							(promo.discount / 100.0) * product_items.price, 
							promo.maximum_discount_amount
						)
					-- price after nominal promo
					WHEN promo.discount_type = 'nominal' 
						THEN GREATEST
						(
							product_items.price - promo.discount, 
							0
						)
				ELSE product_items.price
				END
			) AS price_after_discount,
			pvariations.variations AS product_variation
		`,
	).
		Joins(`JOIN product_items ON product_items.id = carts.product_item_id`).
		Joins(`JOIN products ON products.id = product_items.product_id`).
		Joins(`JOIN merchants ON merchants.id = products.merchant_id`).
		Joins("LEFT JOIN LATERAL (?) promo ON products.id = promo.product_id", sqPromo).
		Joins("LEFT JOIN LATERAL (?) pm ON products.id = pm.product_id", sqProductMedia).
		Joins("LEFT JOIN LATERAL (?) pvariations ON carts.product_item_id = pvariations.product_item_id", subQueryProductVariations).
		Joins(`LEFT JOIN LATERAL (?) pvm ON pvm.product_variant_id = pvariations.product_variant_id`, sqProductVariantMedia).
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
	queryResult := r.db.Model(&models.Cart{}).Select(`
		merchants.id AS merchant_id,
		merchants.name AS name,
		provinces.name AS merchant_province,
		cities.id AS city_id,
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
		Group("merchants.id, provinces.name, cities.id").
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
	result := r.db.
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Delete(&model)

	if result.RowsAffected == 0 {
		err = e.ErrDataNotFound
		return
	}

	model.ID = id

	return
}
