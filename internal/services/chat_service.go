package services

import (
	"fmt"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatService interface {
	WithTx(tx *gorm.DB) *chatService
	GetAllChatList(userID uuid.UUID) ([]*internalDtos.ChatList, error)
	GetChat(userID uuid.UUID, toUserID uuid.UUID) (model *internalDtos.ChatResponse, err error)
	Create(request *internalDtos.ChatRequest) (err error)
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

func (p *chatService) GetChat(userID uuid.UUID, toUserID uuid.UUID) (model *internalDtos.ChatResponse, err error) {
	user, err := p.chatRepository.GetProfileForToUserID(userID)
	if err != nil {
		return nil, err
	}

	user.Chats, err = p.chatRepository.GetChatWithToUserID(userID, toUserID)
	fmt.Println(len(user.Chats), err)
	if err != nil && err != e.ErrDataNotFound {
		return nil, err
	}

	return user, nil
}

func (p *chatService) Create(request *internalDtos.ChatRequest) (err error) {
	return p.chatRepository.Create(&models.Chat{
		FromUserID: request.FromUserID,
		ToUserID:   request.ToUserID,
		Message:    request.Message,
		MediaURL:   request.MediaURL,
	})
}
