package repositories

import (
	"database/sql"
	"log"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository interface {
	WithTx(tx *gorm.DB) *chatRepository
	GetAllChatList(userID uuid.UUID) (models []*internalDtos.ChatList, err error)
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
	if queryResult.Error != nil {
		err = e.ErrDatabaseOperationFailed
		return
	}

	return
}
