package repositories

import (
	"database/sql"
	"log"

	"github.com/natanaelb27/oceana-backend/internal/constants/product_media_constants"
	"github.com/natanaelb27/oceana-backend/internal/models"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductMediaSnapshotRepository interface {
	WithTx(tx *gorm.DB) *productMediaSnapshotRepository
	CreateProductMediaSnapshot(productSnapshotID uint, productItemID uint, productID uint) (model []*models.ProductMediaSnapshot, err error)
}

type productMediaSnapshotRepository struct {
	db *gorm.DB
}

type ProductMediaSnapshotRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductMediaSnapshotRepository(c *ProductMediaSnapshotRepositoryConfig) ProductMediaSnapshotRepository {
	return &productMediaSnapshotRepository{
		db: c.DB,
	}
}

func (u productMediaSnapshotRepository) WithTx(tx *gorm.DB) *productMediaSnapshotRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *productMediaSnapshotRepository) CreateProductSnapshot(newModel *models.ProductSnapshot) (model *models.ProductSnapshot, err error) {
	queryResult := r.db.Create(newModel)

	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	model = newModel

	return
}
func (r *productMediaSnapshotRepository) CreateProductMediaSnapshot(productSnapshotID uint, productItemID uint, productID uint) (model []*models.ProductMediaSnapshot, err error) {
	var newModel []*models.ProductMediaSnapshot

	subQueryProductMedia := r.db.Model(&models.ProductMedia{}).
		Select(`
			(@productSnapshotID)::int AS product_snapshot_id, 
			media_url,
			type
		`, sql.Named("productSnapshotID", productSnapshotID)).
		Where("product_media.product_id = @productID", sql.Named("productID", productID))

	subQueryProductVariantMedia := r.db.Model(&models.ProductVariantMedia{}).
		Select(`
		(@productSnapshotID)::int AS product_snapshot_id, 
		media_url,
		@type_image AS type
	`,
			sql.Named("productSnapshotID", productSnapshotID),
			sql.Named("type_image", product_media_constants.TypeImage)).
		Joins("JOIN product_variations ON product_variant_media.product_variant_id = product_variations.product_variant_id").
		Where("product_variations.product_item_id = @productItemID", sql.Named("productItemID", productItemID))

	queryResult := r.db.Raw("(?) UNION (?)",
		subQueryProductVariantMedia,
		subQueryProductMedia).
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

	model = newModel

	return
}
