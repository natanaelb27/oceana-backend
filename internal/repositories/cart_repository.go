package repositories

import (
	"database/sql"

	in_cart_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/cart_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
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

func (r *cartRepository) subQueryCartsQuantitySum(merchantId uint, userId uuid.UUID) *gorm.DB {
	query := r.db.Model(&models.Cart{}).
		Select("SUM(carts.quantity)").
		Joins(`JOIN product_items ON product_items.id = carts.product_item_id`).
		Joins(`JOIN products p ON p.id = product_items.product_id`).
		Where("carts.user_id = ?", userId).
		Where("merchants.id = ?", merchantId).
		Where("p.id = products.id")
	return query
}

func (r *cartRepository) subQueryProductVariations() *gorm.DB {
	query := r.db.
		Table("product_variations pvariations").
		Select(`string_agg(CONCAT(pvarianttype."name", ':', pvariants."name"),', ' ORDER BY pvarianttype.id) AS value`).
		Joins(`JOIN product_variants pvariants ON pvariants.id = pvariations.product_variant_id `).
		Joins(`JOIN product_variant_types pvarianttype ON pvarianttype.id = pvariants.product_variant_type_id  `).
		Where("pvariations.product_item_id = carts.product_item_id").
		Group("pvariations.product_item_id")
	return query
}

func (r *cartRepository) GetCartsProductsByUserIdMerchants(merchantId uint, userId uuid.UUID) (result []*in_cart_dtos.Product, err error) {
	sqTotalViewCount := sql.Named("all_product_items_quantity", r.subQueryCartsQuantitySum(merchantId, userId))
	subQueryProductVariations := sql.Named("subQueryProductVariations", r.subQueryProductVariations())
	queryResult := r.db.Model(&models.Cart{}).Select(`
			carts.id AS carts_id,
			carts.user_id AS user_id,
			products.id AS product_id,
			products.name AS product_name,
			carts.product_item_id AS product_item_id,
			carts.quantity AS quantity,
			product_items.stock_quantity AS stock_quantity,
			product_items.price AS price,
			-- count price after discount
			(
				CASE
					WHEN (@all_product_items_quantity) > pmp.current_quota
						THEN product_items.price
					WHEN pmp.discount_type = 'percentage' 
						THEN product_items.price - LEAST
						(
							(pmp.discount / 100.0) * product_items.price, 
							pmp.maximum_discount_amount
						)
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
	).
		Joins(`JOIN product_items ON product_items.id = carts.product_item_id`).
		Joins(`JOIN products ON products.id = product_items.product_id`).
		Joins("LEFT JOIN product_merchant_promotions pmp ON products.id = pmp.product_id").
		Joins(`JOIN merchants ON merchants.id = products.merchant_id`).
		Where("carts.user_id = ?", userId).
		Where("merchants.id = ?", merchantId).
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
		Group("products.id, merchants.id, provinces.name").
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
