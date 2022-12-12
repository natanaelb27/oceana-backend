package repositories

import (
	"fmt"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"gorm.io/gorm"
)

type ProductReviewRepository interface {
	Create(newModel *models.ProductReview) (model *models.ProductReview, err error)
	GetByProductID(id uint, pagination *utils.Pagination, query *externalDtos.ProductReviewQuery) ([]*dtos.ProductReviewList, *utils.Pagination, error)
}

type productReviewRepository struct {
	db *gorm.DB
}

type ProductReviewRepositoryConfig struct {
	DB *gorm.DB
}

func NewProductReviewRepository(c *ProductReviewRepositoryConfig) ProductReviewRepository {
	return &productReviewRepository{
		db: c.DB,
	}
}

func (pr *productReviewRepository) Create(newModel *models.ProductReview) (model *models.ProductReview, err error) {
	queryResult := pr.db.
		Where("order_detail_id = ?", newModel.OrderDetailID).
		FirstOrCreate(&newModel)
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

func (pr *productReviewRepository) GetByProductID(id uint, pagination *utils.Pagination, query *externalDtos.ProductReviewQuery) ([]*dtos.ProductReviewList, *utils.Pagination, error) {
	var productReviews []*dtos.ProductReviewList
	subQuery := pr.db.Table("product_snapshots AS ps").
		Select(`
			STRING_AGG(product_variation_snapshots.variant, ', ')
		`).
		Joins("JOIN product_variation_snapshots ON product_snapshots.id = product_variation_snapshots.product_snapshot_id").
		Where("ps.id = product_snapshots.id").
		Group("ps.id")
	queryResult := pr.db.Model(&models.ProductReview{}).
		Select(`
			product_reviews.id, users.id AS user_id, users.full_name, users.username, users.photo, product_reviews.product_id, product_reviews.rating, product_reviews.review, (?) AS product_variations, product_reviews.created_at
		`, subQuery).
		Joins("JOIN users ON product_reviews.user_id = users.id").
		Joins("JOIN order_details ON product_reviews.order_detail_id = order_details.id").
		Joins("JOIN product_snapshots ON order_details.product_snapshot_id = product_snapshots.id").
		Order(fmt.Sprintf("id %s", query.Sort)).
		Where("product_reviews.product_id = ?", id)
	if query.Comment {
		queryResult = queryResult.Where("product_reviews.review <> ''")
	}
	if len(query.Rating) > 0 {
		queryResult = queryResult.Where("product_reviews.rating IN (?)", query.Rating)
	}
	queryResult.Scopes(utils.Paginate(pagination, queryResult)).
		Find(&productReviews)
	if queryResult.Error != nil {
		return nil, nil, e.ErrDatabaseOperationFailed
	}
	return productReviews, pagination, nil
}
