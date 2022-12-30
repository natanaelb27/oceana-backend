package dtos

import (
	"time"
)

type Chat struct {
	ID         uint   `json:"id"`
	FromUserID uint   `json:"from_user_id"`
	ToUserID   uint   `json:"to_user_id"`
	Message    string `json:"description"`
	MediaURL   string `json:"media_url"`
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

type ChatsResponse struct {
	ToUserID        string `json:"to_user_id"`
	ToUserPhoto     string `json:"to_user_photo"`
	ToUserFullName  string `json:"to_user_full_name"`
	ToMerchantID    string `json:"to_merchant_id"`
	ToMerchantName  string `json:"to_merchant_name"`
	ToMerchantPhoto string `json:"to_merchant_photo"`
	// LastMessage Chat    `json:"last_message"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatsRequest struct {
	ToUserID string `json:"to_user_id" binding:"required"`
	Message  string `json:"message" binding:"required"`
	MediaURL string `json:"media_url"`
}
