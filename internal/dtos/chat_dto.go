package dtos

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID         uint      `json:"id"`
	FromUserID uuid.UUID `json:"from_user_id"`
	ToUserID   uuid.UUID `json:"to_user_id"`
	Message    string    `json:"message"`
	MediaURL   string    `json:"media_url"`
	// Read       bool `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatList struct {
	ToUserID        string    `json:"to_user_id"`
	ToUserPhoto     string    `json:"to_user_photo"`
	ToUserFullName  string    `json:"to_user_full_name"`
	ToMerchantID    string    `json:"to_merchant_id"`
	ToMerchantName  string    `json:"to_merchant_name"`
	ToMerchantPhoto string    `json:"to_merchant_photo"`
	LastMessage     string    `json:"last_message"`
	CreatedAt       time.Time `json:"created_at"`
}

type ChatResponse struct {
	ToUserID        string  `json:"to_user_id"`
	ToUserPhoto     string  `json:"to_user_photo"`
	ToUserFullName  string  `json:"to_user_full_name"`
	ToMerchantID    string  `json:"to_merchant_id"`
	ToMerchantName  string  `json:"to_merchant_name"`
	ToMerchantPhoto string  `json:"to_merchant_photo"`
	Chats           []*Chat `json:"chats" gorm:"-"`
}

type ChatsRequest struct {
	ToUserID string `json:"to_user_id" binding:"required"`
	Message  string `json:"message" binding:"required"`
	MediaURL string `json:"media_url"`
}
