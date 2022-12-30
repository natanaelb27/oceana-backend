package services

import (
	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatService interface {
	WithTx(tx *gorm.DB) *chatService
	GetAllChatList(userID uuid.UUID) ([]*internalDtos.ChatList, error)
}

type chatService struct {
	chatRepository r.ChatRepository
}

type ChatServiceConfig struct {
	ChatRepository r.ChatRepository
}

func NewChatService(c *ChatServiceConfig) ChatService {
	return &chatService{
		chatRepository: c.ChatRepository,
	}
}

func (u chatService) WithTx(tx *gorm.DB) *chatService {
	u.chatRepository = u.chatRepository.WithTx(tx)
	return &u
}

func (p *chatService) GetAllChatList(userID uuid.UUID) ([]*internalDtos.ChatList, error) {
	return p.chatRepository.GetAllChatList(userID)
}
