package repositories

import (
	"log"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/constants/product_media_constants"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductSnapshotRepository interface {
	WithTx(tx *gorm.DB) *productSnapshotRepository
	CreateProductSnapshot(model *models.ProductSnapshot) (newModel *models.ProductSnapshot, err error)
	GetByOrderDetailID(orderDetailID uuid.UUID) (*dtos.ProductSnapshotDetails, error)
}

type productSnapshotRepository struct {
	db *gorm.DB
}

type ProductSnapshotRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductSnapshotRepository(c *ProductSnapshotRepositoryConfig) ProductSnapshotRepository {
	return &productSnapshotRepository{
		db: c.DB,
	}
}

func (u productSnapshotRepository) WithTx(tx *gorm.DB) *productSnapshotRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *productSnapshotRepository) CreateProductSnapshot(newModel *models.ProductSnapshot) (model *models.ProductSnapshot, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}

func (r *productSnapshotRepository) GetByOrderDetailID(orderDetailID uuid.UUID) (*dtos.ProductSnapshotDetails, error) {
	var productSnapshotDetails *dtos.ProductSnapshotDetails

	subQueryProductSnapshotImage := r.db.Model(&models.ProductSnapshot{}).
		Select(`
			product_snapshots.id AS product_snapshot_id,
			ARRAY_AGG(product_media_snapshots.media_url) AS images
		`).
		Joins("JOIN product_media_snapshots ON product_snapshots.id = product_media_snapshots.product_snapshot_id").
		Where("product_media_snapshots.type = ?", product_media_constants.TypeImage).
		Group("product_snapshots.id")

	subQueryProductSnapshotVideo := r.db.Model(&models.ProductSnapshot{}).
		Select(`
			product_snapshots.id AS product_snapshot_id,
			ARRAY_AGG(product_media_snapshots.media_url) AS videos
		`).
		Joins("JOIN product_media_snapshots ON product_snapshots.id = product_media_snapshots.product_snapshot_id").
		Where("product_media_snapshots.type = ?", product_media_constants.TypeVideo).
		Group("product_snapshots.id")

	subQueryProductSnapshotVariation := r.db.Model(&models.ProductSnapshot{}).
		Select(`
			product_snapshots.id AS product_snapshot_id,
			STRING_AGG(product_variation_snapshots.variant, ', ') AS product_variations
		`).
		Joins("JOIN product_variation_snapshots ON product_snapshots.id = product_variation_snapshots.product_snapshot_id").
		Group("product_snapshots.id")

	queryResult := r.db.Model(&models.ProductSnapshot{}).
		Select(`
			product_snapshots.id AS product_snapshot_id,
			transactions.user_id AS user_id,
			transactions.transaction_datetime AS transaction_datetime,
			product_snapshots.product_id AS product_id,
			merchants.photo AS merchant_photo,
			merchants.id AS merchant_id,
			merchants.name AS shop_name,
			product_snapshots.name AS product_name,
			products.name AS current_product_name,
			product_snapshots.is_hazardous AS is_hazardous,
			product_snapshots.weight AS weight,
			product_snapshots.package_size AS package_size,
			product_snapshots.condition AS condition,
			product_snapshots.description AS description,
			products.status AS current_status,
			COALESCE(tbl_product_snapshot_image.images, '{}') AS images,
			COALESCE(tbl_product_snapshot_video.videos, '{}') AS videos,
			tbl_product_snapshot_variation.product_variations AS product_variations,
			product_snapshots.price AS price
		`).
		Joins("LEFT JOIN products ON products.id = product_snapshots.product_id").
		Joins("JOIN order_details ON order_details.product_snapshot_id = product_snapshots.id").
		Joins("JOIN orders ON order_details.order_id = orders.id").
		Joins("JOIN merchants ON merchants.id = orders.merchant_id").
		Joins("JOIN transactions ON orders.transaction_id = transactions.id").
		Joins("LEFT JOIN (?) tbl_product_snapshot_variation ON tbl_product_snapshot_variation.product_snapshot_id = product_snapshots.id", subQueryProductSnapshotVariation).
		Joins("LEFT JOIN (?) tbl_product_snapshot_image ON tbl_product_snapshot_image.product_snapshot_id = product_snapshots.id", subQueryProductSnapshotImage).
		Joins("LEFT JOIN (?) tbl_product_snapshot_video ON tbl_product_snapshot_video.product_snapshot_id = product_snapshots.id", subQueryProductSnapshotVideo).
		Where("order_details.id = ?", orderDetailID).
		First(&productSnapshotDetails)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrDataNotFound
		}

		return nil, e.ErrDatabaseOperationFailed
	}

	return productSnapshotDetails, nil
}
