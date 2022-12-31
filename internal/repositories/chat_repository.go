package repositories

import (
	"database/sql"
	"log"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository interface {
	WithTx(tx *gorm.DB) *chatRepository
	GetAllChatList(userID uuid.UUID) (models []*internalDtos.ChatList, err error)
	GetProfileForToUserID(userID uuid.UUID) (model *internalDtos.ChatResponse, err error)
	GetChatWithToUserID(userID, toUserID uuid.UUID) (models []*internalDtos.Chat, err error)
	Create(model *models.Chat) (err error)
}

type chatRepository struct {
	db *gorm.DB
}

type ChatRepositoryConfig struct {
	DB *gorm.DB
}

func NewChatRepository(c *ChatRepositoryConfig) ChatRepository {
	return &chatRepository{
		db: c.DB,
	}
}

func (u chatRepository) WithTx(tx *gorm.DB) *chatRepository {
	if tx == nil {
		log.Print("Transaction Database not found")
		return &u
	}
	u.db = tx
	return &u
}

func (r *chatRepository) GetAllChatList(userID uuid.UUID) (models []*internalDtos.ChatList, err error) {
	queryResult := r.db.Raw(`
	-- Get chat lists
	SELECT 
		u.id AS to_user_id,
		u.photo AS to_user_photo,
		u.full_name AS to_user_full_name,
		COALESCE(m.id, 0) AS to_merchant_id,
		m.name AS to_merchant_name,
		m.photo AS to_merchant_photo,
		lc.message AS last_message,
		lc.created_at AS created_at
	-- get lists of ids this user chat with
	FROM (
		SELECT
			to_user_id AS to_user_id,
			created_at 
		FROM
			chats
		WHERE 
			from_user_id = (@userID)
		UNION ALL
		SELECT 
			from_user_id AS to_user_id,
			created_at
		FROM
			chats
		WHERE 
			to_user_id  = (@userID)
	) chat_lists
	JOIN users u ON u.id = chat_lists.to_user_id 
	LEFT JOIN merchants m ON m.user_id = chat_lists.to_user_id 
	-- get last message from with a user
	LEFT JOIN LATERAL (
		SELECT
			chat_lists.to_user_id AS to_user_id,
			message,
			created_at
		FROM
			chats
		WHERE 
			(from_user_id = (@userID)
				OR
			to_user_id = (@userID))
			AND 
			(from_user_id = chat_lists.to_user_id
				OR
			to_user_id = chat_lists.to_user_id)
		ORDER BY 
			created_at DESC
		LIMIT 1
	) lc ON lc.to_user_id = chat_lists.to_user_id
	GROUP BY chat_lists.to_user_id, m.id, u.id, lc.message, lc.created_at
	ORDER BY chat_lists.to_user_id DESC
	`,
		sql.Named("userID", userID),
	).Find(&models)
	if queryResult.Error != nil && queryResult.Error != gorm.ErrRecordNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *chatRepository) GetProfileForToUserID(userID uuid.UUID) (model *internalDtos.ChatResponse, err error) {
	queryResult := r.db.
		Select(`
			u.id AS to_user_id,
			u.photo AS to_user_photo,
			u.full_name AS to_user_full_name,
			m.id AS to_merchant_id,
			m.name AS to_merchant_name,
			m.photo AS to_merchant_photo
		`).
		Table("users u").
		Joins("LEFT JOIN merchants m ON m.user_id = u.id").
		Where("u.id = ?", userID).
		Scan(&model)
	if queryResult.Error != nil && queryResult.Error != gorm.ErrRecordNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *chatRepository) GetChatWithToUserID(userID, toUserID uuid.UUID) (models []*internalDtos.Chat, err error) {
	queryResult := r.db.
		Raw(`
			SELECT
				id AS id,
				from_user_id,
				to_user_id, 
				message, 
				media_url, 
				created_at
			FROM
				chats c
			WHERE 
				(from_user_id = (@userID)
					OR
				to_user_id = (@userID))
				AND 
				(from_user_id = (@toUserID)
					OR
				to_user_id = (@toUserID))
			ORDER BY 
				created_at DESC
		`,
			sql.Named("userID", userID),
			sql.Named("toUserID", toUserID),
		).
		Find(&models)
	if queryResult.Error != nil && queryResult.Error != gorm.ErrRecordNotFound {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.Error == gorm.ErrRecordNotFound {
		err = e.ErrDataNotFound
		return
	}

	return
}

func (r *chatRepository) Create(model *models.Chat) (err error) {
	queryResult := r.db.Debug().
		Create(&model)
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}
	if queryResult.RowsAffected == 0 {
		err = e.ErrMerchantVoucherCodeAlreadyExist
		return
	}

	return
}
