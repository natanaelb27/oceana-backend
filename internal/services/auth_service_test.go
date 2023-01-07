package services_test

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	internalDtos "github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/services"
	"github.com/natanaelb27/oceana-backend/internal/utils"
	externalDtos "github.com/natanaelb27/oceana-backend/pkg/dtos"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_authService_Register(t *testing.T) {
	type fields struct {
		userRepository   *r.MockUserRepository
		emailRepository  *r.MockEmailRepository
		walletRepository *r.MockWalletRepository
		redisRepository  *r.MockRedisRepository
	}
	type args struct {
		userCreateRequest *externalDtos.UserCreateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *externalDtos.UserCreateResponse
		wantErr bool
	}{
		{
			args: args{
				userCreateRequest: &externalDtos.UserCreateRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.emailRepository.On("Create", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("Create", mock.Anything).Return(&models.User{}, nil)
				r.walletRepository.On("Create", mock.Anything).Return(&models.Wallet{}, nil)
			},
			want:    &externalDtos.UserCreateResponse{},
			wantErr: false,
		},
		{
			args: args{
				userCreateRequest: &externalDtos.UserCreateRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				userCreateRequest: &externalDtos.UserCreateRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				userCreateRequest: &externalDtos.UserCreateRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				userCreateRequest: &externalDtos.UserCreateRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.emailRepository.On("Create", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				userCreateRequest: &externalDtos.UserCreateRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.emailRepository.On("Create", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("Create", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				userCreateRequest: &externalDtos.UserCreateRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.emailRepository.On("Create", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("Create", mock.Anything).Return(&models.User{}, nil)
				r.walletRepository.On("Create", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:   r.NewMockUserRepository(t),
				emailRepository:  r.NewMockEmailRepository(t),
				walletRepository: r.NewMockWalletRepository(t),
				redisRepository:  r.NewMockRedisRepository(t),
			}
			s := services.NewAuthService(&services.AuthServiceConfig{
				UserRepository:   tt.fields.userRepository,
				EmailRepository:  tt.fields.emailRepository,
				WalletRepository: tt.fields.walletRepository,
				RedisRepository:  tt.fields.redisRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := s.Register(tt.args.userCreateRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authService.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authService_Login(t *testing.T) {
	t.Setenv("REFRESH_TOKEN_DAYS_DURATION", "1")
	t.Setenv("JWT_EXPIRATION_HOURS", "1")
	t.Setenv("JWT_SECRET_KEY", "1")
	t.Setenv("JWT_ISSUER", "1")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	password, _ := utils.Hash("password")

	type fields struct {
		userRepository   *r.MockUserRepository
		emailRepository  *r.MockEmailRepository
		walletRepository *r.MockWalletRepository
		redisRepository  *r.MockRedisRepository
	}
	type args struct {
		c                *gin.Context
		userLoginRequest *externalDtos.UserLoginRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *internalDtos.UserLoginData
		wantErr bool
	}{
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(&models.User{
					Password: password,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			want:    &internalDtos.UserLoginData{},
			wantErr: false,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(&models.User{}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(&models.User{
					Password: password,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(&models.User{
					Password: password,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			want:    &internalDtos.UserLoginData{},
			wantErr: false,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(&models.User{}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(&models.User{
					Password: password,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},

		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrUnknown)
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(&models.User{
					Password: password,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			want:    &internalDtos.UserLoginData{},
			wantErr: false,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrUnknown)
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrUnknown)
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(&models.User{}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrUnknown)
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(&models.User{
					Password: password,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrUnknown)
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginRequest: &externalDtos.UserLoginRequest{
					LoginCredentials: "john.doe",
					Password:         "",
				},
			},
			mock:    func(r *fields, a args) {},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:   r.NewMockUserRepository(t),
				emailRepository:  r.NewMockEmailRepository(t),
				walletRepository: r.NewMockWalletRepository(t),
				redisRepository:  r.NewMockRedisRepository(t),
			}
			s := services.NewAuthService(&services.AuthServiceConfig{
				UserRepository:   tt.fields.userRepository,
				EmailRepository:  tt.fields.emailRepository,
				WalletRepository: tt.fields.walletRepository,
				RedisRepository:  tt.fields.redisRepository,
			})
			tt.mock(&tt.fields, tt.args)
			_, err := s.Login(tt.args.c, tt.args.userLoginRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_authService_RegisterWithGoogle(t *testing.T) {
	t.Setenv("REFRESH_TOKEN_DAYS_DURATION", "1")
	t.Setenv("JWT_EXPIRATION_HOURS", "1")
	t.Setenv("JWT_SECRET_KEY", "1")
	t.Setenv("JWT_ISSUER", "1")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	googleID := "google-id"

	type fields struct {
		userRepository   *r.MockUserRepository
		emailRepository  *r.MockEmailRepository
		walletRepository *r.MockWalletRepository
		redisRepository  *r.MockRedisRepository
	}
	type args struct {
		c                             *gin.Context
		userRegisterWithGoogleRequest *externalDtos.UserRegisterWithGoogleRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *internalDtos.UserLoginData
		wantErr bool
	}{
		{
			args: args{
				c: c,
				userRegisterWithGoogleRequest: &externalDtos.UserRegisterWithGoogleRequest{
					GoogleID: &googleID,
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.emailRepository.On("Create", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("Create", mock.Anything).Return(&models.User{}, nil)
				r.walletRepository.On("Create", mock.Anything).Return(&models.Wallet{}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			want:    &internalDtos.UserLoginData{},
			wantErr: false,
		},
		{
			args: args{
				c:                             c,
				userRegisterWithGoogleRequest: &externalDtos.UserRegisterWithGoogleRequest{},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:                             c,
				userRegisterWithGoogleRequest: &externalDtos.UserRegisterWithGoogleRequest{},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:                             c,
				userRegisterWithGoogleRequest: &externalDtos.UserRegisterWithGoogleRequest{},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.emailRepository.On("Create", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userRegisterWithGoogleRequest: &externalDtos.UserRegisterWithGoogleRequest{
					GoogleID: &googleID,
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.emailRepository.On("Create", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("Create", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userRegisterWithGoogleRequest: &externalDtos.UserRegisterWithGoogleRequest{
					GoogleID: &googleID,
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.emailRepository.On("Create", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("Create", mock.Anything).Return(&models.User{}, nil)
				r.walletRepository.On("Create", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userRegisterWithGoogleRequest: &externalDtos.UserRegisterWithGoogleRequest{
					GoogleID: &googleID,
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByUsername", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.userRepository.On("GetByPhoneNumber", mock.Anything).Return(nil, e.ErrDataNotFound)
				r.emailRepository.On("Create", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("Create", mock.Anything).Return(&models.User{}, nil)
				r.walletRepository.On("Create", mock.Anything).Return(&models.Wallet{}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:   r.NewMockUserRepository(t),
				emailRepository:  r.NewMockEmailRepository(t),
				walletRepository: r.NewMockWalletRepository(t),
				redisRepository:  r.NewMockRedisRepository(t),
			}
			s := services.NewAuthService(&services.AuthServiceConfig{
				UserRepository:   tt.fields.userRepository,
				EmailRepository:  tt.fields.emailRepository,
				WalletRepository: tt.fields.walletRepository,
				RedisRepository:  tt.fields.redisRepository,
			})
			tt.mock(&tt.fields, tt.args)
			_, err := s.RegisterWithGoogle(tt.args.c, tt.args.userRegisterWithGoogleRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.RegisterWithGoogle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_authService_LoginWithGoogle(t *testing.T) {
	t.Setenv("REFRESH_TOKEN_DAYS_DURATION", "1")
	t.Setenv("JWT_EXPIRATION_HOURS", "1")
	t.Setenv("JWT_SECRET_KEY", "1")
	t.Setenv("JWT_ISSUER", "1")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	googleID, _ := utils.Hash("google-id")

	type fields struct {
		userRepository   *r.MockUserRepository
		emailRepository  *r.MockEmailRepository
		walletRepository *r.MockWalletRepository
		redisRepository  *r.MockRedisRepository
	}
	type args struct {
		c                          *gin.Context
		userLoginWithGoogleRequest *externalDtos.UserLoginWithGoogleRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *internalDtos.UserLoginData
		wantErr bool
	}{
		{
			args: args{
				c: c,
				userLoginWithGoogleRequest: &externalDtos.UserLoginWithGoogleRequest{
					GoogleID: "google-id",
				},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(&models.User{
					GoogleID: &googleID,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			},
			want:    &internalDtos.UserLoginData{},
			wantErr: false,
		},
		{
			args: args{
				c: c,
				userLoginWithGoogleRequest: &externalDtos.UserLoginWithGoogleRequest{
					GoogleID: "google-id",
				},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrUnknown)

			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginWithGoogleRequest: &externalDtos.UserLoginWithGoogleRequest{
					GoogleID: "google-id",
				},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(nil, e.ErrUnknown)

			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginWithGoogleRequest: &externalDtos.UserLoginWithGoogleRequest{
					GoogleID: "wrong-google-id",
				},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(&models.User{
					GoogleID: &googleID,
				}, nil)

			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c: c,
				userLoginWithGoogleRequest: &externalDtos.UserLoginWithGoogleRequest{
					GoogleID: "google-id",
				},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(&models.User{
					GoogleID: &googleID,
				}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(e.ErrUnknown)

			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:   r.NewMockUserRepository(t),
				emailRepository:  r.NewMockEmailRepository(t),
				walletRepository: r.NewMockWalletRepository(t),
				redisRepository:  r.NewMockRedisRepository(t),
			}
			s := services.NewAuthService(&services.AuthServiceConfig{
				UserRepository:   tt.fields.userRepository,
				EmailRepository:  tt.fields.emailRepository,
				WalletRepository: tt.fields.walletRepository,
				RedisRepository:  tt.fields.redisRepository,
			})
			tt.mock(&tt.fields, tt.args)
			_, err := s.LoginWithGoogle(tt.args.c, tt.args.userLoginWithGoogleRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.LoginWithGoogle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_authService_StepUpPassword(t *testing.T) {
	t.Setenv("REFRESH_TOKEN_DAYS_DURATION", "1")
	t.Setenv("JWT_EXPIRATION_HOURS", "1")
	t.Setenv("JWT_SECRET_KEY", "1")
	t.Setenv("JWT_ISSUER", "1")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	password, _ := utils.Hash("password")

	type fields struct {
		userRepository   *r.MockUserRepository
		emailRepository  *r.MockEmailRepository
		walletRepository *r.MockWalletRepository
		redisRepository  *r.MockRedisRepository
	}
	type args struct {
		c       *gin.Context
		id      uuid.UUID
		request *externalDtos.AuthStepUpPasswordRequest
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mock      func(r *fields, a args)
		wantToken string
		wantErr   bool
	}{
		{
			args: args{
				c:  c,
				id: uuid.UUID{},
				request: &externalDtos.AuthStepUpPasswordRequest{
					Password: "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Password: password,
				}, nil)
			},
			wantToken: "",
			wantErr:   false,
		},
		{
			args: args{
				c:  c,
				id: uuid.UUID{},
				request: &externalDtos.AuthStepUpPasswordRequest{
					Password: "",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Password: password,
				}, nil)
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			args: args{
				c:  c,
				id: uuid.UUID{},
				request: &externalDtos.AuthStepUpPasswordRequest{
					Password: "password",
				},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByID", mock.Anything).Return(nil, e.ErrUnknown)
			},
			wantToken: "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:   r.NewMockUserRepository(t),
				emailRepository:  r.NewMockEmailRepository(t),
				walletRepository: r.NewMockWalletRepository(t),
				redisRepository:  r.NewMockRedisRepository(t),
			}
			s := services.NewAuthService(&services.AuthServiceConfig{
				UserRepository:   tt.fields.userRepository,
				EmailRepository:  tt.fields.emailRepository,
				WalletRepository: tt.fields.walletRepository,
				RedisRepository:  tt.fields.redisRepository,
			})
			tt.mock(&tt.fields, tt.args)
			_, err := s.StepUpPassword(tt.args.c, tt.args.id, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.StepUpPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_authService_SendResetPasswordEmail(t *testing.T) {
	t.Setenv("REFRESH_TOKEN_DAYS_DURATION", "1")
	t.Setenv("JWT_EXPIRATION_HOURS", "1")
	t.Setenv("JWT_SECRET_KEY", "1")
	t.Setenv("JWT_ISSUER", "1")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	type fields struct {
		userRepository   *r.MockUserRepository
		emailRepository  *r.MockEmailRepository
		walletRepository *r.MockWalletRepository
		redisRepository  *r.MockRedisRepository
	}
	type args struct {
		c       *gin.Context
		request *externalDtos.UserVerifyResetPasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    bool
		wantErr bool
	}{
		{
			args: args{
				c:       c,
				request: &externalDtos.UserVerifyResetPasswordRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(nil, e.ErrUnknown)

			},
			want:    false,
			wantErr: true,
		},
		{
			args: args{
				c:       c,
				request: &externalDtos.UserVerifyResetPasswordRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(nil, e.ErrUnknown)

			},
			want:    false,
			wantErr: true,
		},
		{
			args: args{
				c:       c,
				request: &externalDtos.UserVerifyResetPasswordRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(e.ErrUnknown)

			},
			want:    false,
			wantErr: true,
		},
		{
			args: args{
				c:       c,
				request: &externalDtos.UserVerifyResetPasswordRequest{},
			},
			mock: func(r *fields, a args) {
				r.emailRepository.On("GetByEmail", mock.Anything).Return(&models.Email{}, nil)
				r.userRepository.On("GetByEmailID", mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				r.redisRepository.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(e.ErrUnknown).Once()

			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:   r.NewMockUserRepository(t),
				emailRepository:  r.NewMockEmailRepository(t),
				walletRepository: r.NewMockWalletRepository(t),
				redisRepository:  r.NewMockRedisRepository(t),
			}
			s := services.NewAuthService(&services.AuthServiceConfig{
				UserRepository:   tt.fields.userRepository,
				EmailRepository:  tt.fields.emailRepository,
				WalletRepository: tt.fields.walletRepository,
				RedisRepository:  tt.fields.redisRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := s.SendResetPasswordEmail(tt.args.c, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.SendResetPasswordEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("authService.SendResetPasswordEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authService_ResetPassword(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	password, _ := utils.Hash("password")

	type fields struct {
		userRepository   *r.MockUserRepository
		emailRepository  *r.MockEmailRepository
		walletRepository *r.MockWalletRepository
		redisRepository  *r.MockRedisRepository
	}
	type args struct {
		c               *gin.Context
		resetPasswordID string
		request         *externalDtos.ResetPasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *externalDtos.ResetPasswordResponse
		wantErr bool
	}{
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "@Password123",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return(uuid.New().String(), nil).Once()
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil).Once()
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
					Password: password,
				}, nil)
				r.userRepository.On("Update", mock.Anything, mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
				r.redisRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
			},
			want:    externalDtos.NewResetPasswordResponse(true),
			wantErr: false,
		},
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "@Password123",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return(uuid.New().String(), nil).Once()
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil).Once()
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
					Password: password,
				}, nil)
				r.userRepository.On("Update", mock.Anything, mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
				r.redisRepository.On("Delete", mock.Anything, mock.Anything).Return(e.ErrUnknown).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "@Password123",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return(uuid.New().String(), nil).Once()
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil).Once()
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
					Password: password,
				}, nil)
				r.userRepository.On("Update", mock.Anything, mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Delete", mock.Anything, mock.Anything).Return(e.ErrUnknown).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "@Password123",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return(uuid.New().String(), nil).Once()
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil).Once()
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
					Password: password,
				}, nil)
				r.userRepository.On("Update", mock.Anything, mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "@Password123",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return(uuid.New().String(), nil).Once()
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil).Once()
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "password",
					Password: password,
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "password",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return(uuid.New().String(), nil).Once()
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil).Once()
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{
					Username: "username",
					Password: password,
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "password",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return(uuid.New().String(), nil).Once()
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil).Once()
				r.userRepository.On("GetByID", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "password",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return(uuid.New().String(), nil).Once()
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrUnknown).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "password",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil).Once()
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				c:               c,
				resetPasswordID: "",
				request: &externalDtos.ResetPasswordRequest{
					Password: "password",
				},
			},
			mock: func(r *fields, a args) {
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrUnknown).Once()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:   r.NewMockUserRepository(t),
				emailRepository:  r.NewMockEmailRepository(t),
				walletRepository: r.NewMockWalletRepository(t),
				redisRepository:  r.NewMockRedisRepository(t),
			}
			s := services.NewAuthService(&services.AuthServiceConfig{
				UserRepository:   tt.fields.userRepository,
				EmailRepository:  tt.fields.emailRepository,
				WalletRepository: tt.fields.walletRepository,
				RedisRepository:  tt.fields.redisRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := s.ResetPassword(tt.args.c, tt.args.resetPasswordID, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authService.ResetPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authService_RefreshAccessToken(t *testing.T) {
	t.Setenv("REFRESH_TOKEN_DAYS_DURATION", "1")
	t.Setenv("JWT_EXPIRATION_HOURS", "1")
	t.Setenv("JWT_SECRET_KEY", "1")
	t.Setenv("JWT_ISSUER", "1")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	type fields struct {
		userRepository   *r.MockUserRepository
		emailRepository  *r.MockEmailRepository
		walletRepository *r.MockWalletRepository
		redisRepository  *r.MockRedisRepository
	}
	type args struct {
		c            *gin.Context
		refreshToken string
		userID       uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    string
		wantErr bool
	}{
		{
			args: args{
				c:            c,
				refreshToken: "",
				userID:       uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", nil)
			},
			want:    "",
			wantErr: false,
		},
		{
			args: args{
				c:            c,
				refreshToken: "",
				userID:       uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByID", mock.Anything).Return(nil, e.ErrUnknown)

			},
			want:    "",
			wantErr: true,
		},
		{
			args: args{
				c:            c,
				refreshToken: "",
				userID:       uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("something", nil)

			},
			want:    "",
			wantErr: true,
		},
		{
			args: args{
				c:            c,
				refreshToken: "",
				userID:       uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.userRepository.On("GetByID", mock.Anything).Return(&models.User{}, nil)
				r.redisRepository.On("Get", mock.Anything, mock.Anything).Return("", e.ErrUnknown)
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				userRepository:   r.NewMockUserRepository(t),
				emailRepository:  r.NewMockEmailRepository(t),
				walletRepository: r.NewMockWalletRepository(t),
				redisRepository:  r.NewMockRedisRepository(t),
			}
			s := services.NewAuthService(&services.AuthServiceConfig{
				UserRepository:   tt.fields.userRepository,
				EmailRepository:  tt.fields.emailRepository,
				WalletRepository: tt.fields.walletRepository,
				RedisRepository:  tt.fields.redisRepository,
			})
			tt.mock(&tt.fields, tt.args)
			_, err := s.RefreshAccessToken(tt.args.c, tt.args.refreshToken, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.RefreshAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
