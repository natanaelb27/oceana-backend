package repositories

import (
	"fmt"
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants/product_media_constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductReviewRepository interface {
	WithTx(tx *gorm.DB) *productReviewRepository
	Create(newModel *models.ProductReview) (model *models.ProductReview, err error)
	GetByProductID(id uint, pagination *utils.Pagination, query *externalDtos.ProductReviewQuery) ([]*dtos.ProductReviewList, *utils.Pagination, error)
	GetDetailsByProductID(id uint) (*dtos.ProductReviewDetails, error)
	GetByOrderDetailID(id uuid.UUID) (model *models.ProductReview, err error)
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

func (u productReviewRepository) WithTx(tx *gorm.DB) *productReviewRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
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
	productVariantSubQuery := pr.db.Table("product_snapshots AS ps").
		Select(`
			STRING_AGG(product_variation_snapshots.variant, ', ')
		`).
		Joins("JOIN product_variation_snapshots ON product_snapshots.id = product_variation_snapshots.product_snapshot_id").
		Where("ps.id = product_snapshots.id").
		Group("ps.id")

	subQueryReviewImage := pr.db.Model(&models.ProductReviewMedia{}).
		Select(`
			product_reviews.product_id,
			ARRAY_AGG(product_review_media.media_url) AS review_images
		`).
		Joins("JOIN product_reviews ON product_reviews.id = product_review_media.product_review_id").
		Where("type = ?", product_media_constants.TypeImage).
		Group("product_reviews.product_id")

	subQueryReviewVideo := pr.db.Model(&models.ProductReviewMedia{}).
		Select(`
			product_reviews.product_id,
			ARRAY_AGG(product_review_media.media_url) AS review_videos
		`).
		Joins("JOIN product_reviews ON product_reviews.id = product_review_media.product_review_id").
		Where("type = ?", product_media_constants.TypeVideo).
		Group("product_reviews.product_id")

	// mediaSubQuery := pr.db.Model(&models.ProductReview{}).
	// 	Select(`
	// 		product_reviews.id AS product_review_id,
	// 		ARRAY_REMOVE(ARRAY_AGG(product_review_media.media_url), NULL) AS medias
	// 	`).
	// 	Joins("LEFT JOIN product_review_media ON product_review_media.product_review_id = product_reviews.id").
	// 	Group("product_reviews.id")

	queryResult := pr.db.Model(&models.ProductReview{}).
		Select(`
			product_reviews.id, 
			users.id AS user_id, 
			users.full_name, 
			users.username, 
			users.photo, 
			product_reviews.product_id, 
			product_reviews.rating, 
			product_reviews.review, 
			(?) AS product_variations, 
			tbl_review_image.review_images AS review_images,
			tbl_review_video.review_videos AS review_videos,
			product_reviews.created_at
		`, productVariantSubQuery).
		Joins("JOIN users ON product_reviews.user_id = users.id").
		Joins("JOIN order_details ON product_reviews.order_detail_id = order_details.id").
		Joins("JOIN product_snapshots ON order_details.product_snapshot_id = product_snapshots.id").
		Joins("LEFT JOIN (?) tbl_review_image ON product_reviews.product_id = tbl_review_image.product_id", subQueryReviewImage).
		Joins("LEFT JOIN (?) tbl_review_video ON product_reviews.product_id = tbl_review_video.product_id", subQueryReviewVideo).
		Order(fmt.Sprintf("%s %s", query.SortBy, query.Sort)).
		Where("product_reviews.product_id = ?", id)

	if query.Media {
		queryResult = queryResult.Where("ARRAY_LENGTH(COALESCE(medias_agg.medias, '{}'), 1) > 0")
	}

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

func (pr *productReviewRepository) GetDetailsByProductID(id uint) (*dtos.ProductReviewDetails, error) {
	var productReviewDetails *dtos.ProductReviewDetails

	subQueryReview := pr.db.Model(&models.ProductReview{}).
		Select(`
			product_reviews.product_id AS product_id,
			COUNT(product_reviews.id) AS review_count
		`).
		Where("product_reviews.product_id = ?", id).
		Where("product_reviews.review IS NOT NULL").
		Where("TRIM(product_reviews.review) <> ''").
		Group("product_reviews.product_id")

	subQueryRating := pr.db.Model(&models.ProductReview{}).
		Select(`
			product_reviews.product_id AS product_id,
			COUNT(product_reviews.id) AS rating_count,
			AVG(product_reviews.rating) AS rating_star
		`).
		Where("product_reviews.product_id = ?", id).
		Group("product_reviews.product_id")

	subQuerySatisfied := pr.db.Model(&models.ProductReview{}).
		Select(`
			product_reviews.product_id AS product_id,
			COUNT(product_reviews.id) AS satisfied_rating_count
		`).
		Where("product_reviews.product_id = ?", id).
		Where("product_reviews.rating >= 4").
		Group("product_reviews.product_id")

	subQueryPerRating := pr.db.Model(&models.ProductReview{}).
		Select(`
			product_reviews.product_id AS product_id,
			product_reviews.rating AS rating,
			COUNT(product_reviews.id) AS rating_count
		`).
		Where("product_reviews.product_id = ?", id).
		Group("product_reviews.product_id, product_reviews.rating")

	queryResult := pr.db.Model(&models.Product{}).
		Select(`
			tbl_rating.rating_count AS rating_count,
			tbl_review_count.review_count AS review_count,
			tbl_rating.rating_star AS rating_star,
			COALESCE(tbl_satisfied_rating.satisfied_rating_count::FLOAT / NULLIF(tbl_rating.rating_count, 0)::FLOAT, 0) * 100 AS satisfied_percentage,
			COALESCE(tbl_rating_5.rating_count, 0) AS star_5,
			COALESCE(tbl_rating_4.rating_count, 0) AS star_4,
			COALESCE(tbl_rating_3.rating_count, 0) AS star_3,
			COALESCE(tbl_rating_2.rating_count, 0) AS star_2,
			COALESCE(tbl_rating_1.rating_count, 0) AS star_1
		`).
		Joins("LEFT JOIN (?) tbl_review_count ON products.id = tbl_review_count.product_id", subQueryReview).
		Joins("LEFT JOIN (?) tbl_rating ON products.id = tbl_rating.product_id", subQueryRating).
		Joins("LEFT JOIN (?) tbl_satisfied_rating ON products.id = tbl_satisfied_rating.product_id", subQuerySatisfied).
		Joins("LEFT JOIN (?) tbl_rating_5 ON products.id = tbl_rating_5.product_id AND tbl_rating_5.rating = 5", subQueryPerRating).
		Joins("LEFT JOIN (?) tbl_rating_4 ON products.id = tbl_rating_4.product_id AND tbl_rating_4.rating = 4", subQueryPerRating).
		Joins("LEFT JOIN (?) tbl_rating_3 ON products.id = tbl_rating_3.product_id AND tbl_rating_3.rating = 3", subQueryPerRating).
		Joins("LEFT JOIN (?) tbl_rating_2 ON products.id = tbl_rating_2.product_id AND tbl_rating_2.rating = 2", subQueryPerRating).
		Joins("LEFT JOIN (?) tbl_rating_1 ON products.id = tbl_rating_1.product_id AND tbl_rating_1.rating = 1", subQueryPerRating).
		Where("products.id = ?", id).
		First(&productReviewDetails)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrProductNotFound
		}
		return nil, e.ErrDatabaseOperationFailed
	}

	return productReviewDetails, nil
}
func (pr *productReviewRepository) GetByOrderDetailID(id uuid.UUID) (model *models.ProductReview, err error) {
	queryResult := pr.db.Model(&models.ProductReview{}).
		Where("order_detail_id = ?", id).
		First(&model)

	if queryResult.Error != nil {
		if queryResult.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrDataNotFound
		}

		return nil, e.ErrDatabaseOperationFailed
	}

	return

}
