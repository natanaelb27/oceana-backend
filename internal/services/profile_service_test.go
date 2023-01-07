package services_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/models"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/services"
	externalDtos "github.com/natanaelb27/oceana-backend/pkg/dtos"
	"github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_profileService_UpdateProfilePassword(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	type fields struct {
		userRepository  *r.MockUserRepository
		redisRepository *r.MockRedisRepository
	}
	type args struct {
		c       *gin.Context
		id      uuid.UUID
		request *externalDtos.UpdateProfilePasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields)
		wantErr bool
	}{
		{
			args: args{
				c:  c,
				id: uuid.UUID{},
				request: &externalDtos.UpdateProfilePasswordRequest{
					Password:         "@Password123",
					VerificationCode: "code",
				},
			},
			mock: func(r *fields) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
				}, nil)
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("code", nil)
				r.userRepository.On("Update", mock.Anything, mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Delete", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			args: args{
				c:  c,
				id: uuid.UUID{},
				request: &externalDtos.UpdateProfilePasswordRequest{
					Password:         "@Password123",
					VerificationCode: "code",
				},
			},
			mock: func(r *fields) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
				}, nil)
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("code", nil)
				r.userRepository.On("Update", mock.Anything, mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Delete", mock.Anything, mock.Anything).Return(errors.ErrUnknown)
			},
			wantErr: true,
		},
		{
			args: args{
				c:  c,
				id: uuid.UUID{},
				request: &externalDtos.UpdateProfilePasswordRequest{
					Password:         "@Password123",
					VerificationCode: "code",
				},
			},
			mock: func(r *fields) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
				}, nil)
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("code", nil)
				r.userRepository.On("Update", mock.Anything, mock.Anything).Return(&models.User{}, errors.ErrUnknown)
			},
			wantErr: true,
		},
		{
			args: args{
				c:  c,
				id: uuid.UUID{},
				request: &externalDtos.UpdateProfilePasswordRequest{
					Password:         "@Password123",
					VerificationCode: "code",
				},
			},
			mock: func(r *fields) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
				}, nil)
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("code", errors.ErrUnknown)
			},
			wantErr: true,
		},
		{
			args: args{
				c:  c,
				id: uuid.UUID{},
				request: &externalDtos.UpdateProfilePasswordRequest{
					Password:         "@Password123",
					VerificationCode: "code",
				},
			},
			mock: func(r *fields) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
				}, errors.ErrUnknown)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:  r.NewMockUserRepository(t),
				redisRepository: r.NewMockRedisRepository(t),
			}
			s := services.NewProfileService(&services.ProfileServiceConfig{
				UserRepository:  tt.fields.userRepository,
				RedisRepository: tt.fields.redisRepository,
			})
			tt.mock(&tt.fields)
			_, err := s.UpdateProfilePassword(tt.args.c, tt.args.id, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("profileService.UpdateProfilePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_profileService_CheckPasswordVerificationCode(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	type fields struct {
		userRepository  *r.MockUserRepository
		redisRepository *r.MockRedisRepository
	}
	type args struct {
		c       *gin.Context
		id      uuid.UUID
		request *externalDtos.ProfileCheckPasswordVerificationCodeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields)
		wantErr bool
	}{
		{
			args: args{
				c:       c,
				id:      uuid.UUID{},
				request: &externalDtos.ProfileCheckPasswordVerificationCodeRequest{},
			},
			mock: func(r *fields) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil)
			},
			wantErr: false,
		},
		{
			args: args{
				c:       c,
				id:      uuid.UUID{},
				request: &externalDtos.ProfileCheckPasswordVerificationCodeRequest{},
			},
			mock: func(r *fields) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("code", nil)
			},
			wantErr: true,
		},
		{
			args: args{
				c:       c,
				id:      uuid.UUID{},
				request: &externalDtos.ProfileCheckPasswordVerificationCodeRequest{},
			},
			mock: func(r *fields) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", errors.ErrUnknown)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:  r.NewMockUserRepository(t),
				redisRepository: r.NewMockRedisRepository(t),
			}
			s := services.NewProfileService(&services.ProfileServiceConfig{
				UserRepository:  tt.fields.userRepository,
				RedisRepository: tt.fields.redisRepository,
			})
			tt.mock(&tt.fields)
			_, err := s.CheckPasswordVerificationCode(tt.args.c, tt.args.id, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("profileService.CheckPasswordVerificationCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
