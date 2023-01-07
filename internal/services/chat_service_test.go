package services_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	internalDtos "github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/services"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_chatService_GetAllChatList(t *testing.T) {
	type fields struct {
		chatRepository *r.MockChatRepository
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    []*internalDtos.ChatList
		wantErr bool
	}{
		{
			args: args{
				userID: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.chatRepository.On("GetAllChatList", mock.Anything).Return([]*internalDtos.ChatList{}, nil)
			},
			want:    []*internalDtos.ChatList{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				chatRepository: r.NewMockChatRepository(t),
			}
			s := services.NewChatService(&services.ChatServiceConfig{
				ChatRepository: tt.fields.chatRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := s.GetAllChatList(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatService.GetAllChatList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("chatService.GetAllChatList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_chatService_GetChat(t *testing.T) {
	type fields struct {
		chatRepository *r.MockChatRepository
	}
	type args struct {
		userID   uuid.UUID
		toUserID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *internalDtos.ChatResponse
		wantErr bool
	}{
		{
			args: args{
				userID:   uuid.UUID{},
				toUserID: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.chatRepository.On("GetProfileForToUserID", mock.Anything).Return(&internalDtos.ChatResponse{}, nil)
				r.chatRepository.On("GetChatWithToUserID", mock.Anything, mock.Anything).Return([]*internalDtos.Chat{}, nil)
			},
			want:    &internalDtos.ChatResponse{},
			wantErr: false,
		},
		{
			args: args{
				userID:   uuid.UUID{},
				toUserID: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.chatRepository.On("GetProfileForToUserID", mock.Anything).Return(&internalDtos.ChatResponse{}, nil)
				r.chatRepository.On("GetChatWithToUserID", mock.Anything, mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				userID:   uuid.UUID{},
				toUserID: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.chatRepository.On("GetProfileForToUserID", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				chatRepository: r.NewMockChatRepository(t),
			}
			s := services.NewChatService(&services.ChatServiceConfig{
				ChatRepository: tt.fields.chatRepository,
			})
			tt.mock(&tt.fields, tt.args)
			_, err := s.GetChat(tt.args.userID, tt.args.toUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatService.GetChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_chatService_Create(t *testing.T) {
	type fields struct {
		chatRepository *r.MockChatRepository
	}
	type args struct {
		request *internalDtos.ChatRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *models.Chat
		wantErr bool
	}{
		{
			args: args{
				request: &internalDtos.ChatRequest{},
			},
			mock: func(r *fields, a args) {
				r.chatRepository.On("Create", mock.Anything).Return(nil)
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				chatRepository: r.NewMockChatRepository(t),
			}
			s := services.NewChatService(&services.ChatServiceConfig{
				ChatRepository: tt.fields.chatRepository,
			})
			tt.mock(&tt.fields, tt.args)
			err := s.Create(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
