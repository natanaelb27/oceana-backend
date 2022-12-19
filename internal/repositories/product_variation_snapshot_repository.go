package repositories

import (
	"database/sql"
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductVariationSnapshotRepository interface {
	WithTx(tx *gorm.DB) *productVariationSnapshotRepository
	CreateProductVariationSnapshot(productSnapshotID uint, productItemID uint) (model []*models.ProductVariationSnapshot, err error)
}

type productVariationSnapshotRepository struct {
	db *gorm.DB
}

type ProductVariationSnapshotRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductVariationSnapshotRepository(c *ProductVariationSnapshotRepositoryConfig) ProductVariationSnapshotRepository {
	return &productVariationSnapshotRepository{
		db: c.DB,
	}
}

func (u productVariationSnapshotRepository) WithTx(tx *gorm.DB) *productVariationSnapshotRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *productVariationSnapshotRepository) CreateProductVariationSnapshot(productSnapshotID uint, productItemID uint) (model []*models.ProductVariationSnapshot, err error) {
	var newModel []*models.ProductVariationSnapshot
	queryResult := r.db.Model(&models.ProductVariation{}).
		Select(`
			(@productSnapshotID)::int AS product_snapshot_id, 
			product_variant_types.name AS variant_type,
			product_variants.name AS variant
		`, sql.Named("productSnapshotID", productSnapshotID)).
		Joins("JOIN product_items ON product_items.id = product_variations.product_item_id").
		Joins("JOIN product_variants ON product_variations.product_variant_id = product_variants.id").
		Joins("JOIN product_variant_types ON product_variants.product_variant_type_id = product_variant_types.id").
		Where("product_variations.product_item_id = @productItemID", sql.Named("productItemID", productItemID)).
		Find(&newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	if len(newModel) <= 0 {
		return nil, nil
	}

	queryResult = r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	// queryResult = r.db.Exec(`
	// 		INSERT INTO product_variation_snapshots (product_snapshot_id, variant_type, variant)
	// 		SELECT @productSnapshotID AS product_snapshot_id,
	// 		product_variant_types.name AS variant_type,
	// 		product_variants.name AS variant
	// 		FROM product_variations
	// 		JOIN product_variations ON product_items.id = product_variations.product_item_id
	// 		JOIN product_variants ON product_variations.product_variant_id = product_variants.id
	// 		JOIN product_variant_types ON product_variants.product_variant_type_id = product_variant_types.id
	// 		WHERE product_variations.product_item_id = @productItemID
	// 	`, sql.Named("productSnapshotID", productSnapshotID),
	// 	sql.Named("productItemID", productItemID))

	model = newModel

	return
}
