package services_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/models"
	"github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/services"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_sealabsPayService_GetAll(t *testing.T) {
	wantModel := &models.SealabsPay{
		ID:            uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		AccountNumber: "2783993778282621",
		NameOnCard:    "Timothy",
		IsDefault:     true,
	}
	wantModels := make([]*models.SealabsPay, 0)
	wantModels = append(wantModels, wantModel)

	type fields struct {
		sealabsPayRepository *repositories.MockSealabsPayRepository
		redisRepository      *repositories.MockRedisRepository
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args, w []*models.SealabsPay)
		want    []*models.SealabsPay
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				userID: wantModel.UserID,
			},
			mock: func(r *fields, a args, w []*models.SealabsPay) {
				r.sealabsPayRepository.On("GetAll", mock.Anything).Return(wantModels, nil)
			},
			want: wantModels,
		},
		{
			args: args{
				userID: wantModel.UserID,
			},
			mock: func(r *fields, a args, w []*models.SealabsPay) {
				r.sealabsPayRepository.On("GetAll", mock.Anything).Return(nil, e.ErrDataNotFound)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				repositories.NewMockSealabsPayRepository(t),
				repositories.NewMockRedisRepository(t),
			}

			sealabsPayService := services.NewSealabsPayService(&services.SealabsPayServiceConfig{
				SealabsPayRepository: tt.fields.sealabsPayRepository,
				RedisRepository:      tt.fields.redisRepository,
			})

			tt.mock(&tt.fields, tt.args, tt.want)

			got, err := sealabsPayService.GetAll(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("sealabsPayService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sealabsPayService.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sealabsPayService_GetById(t *testing.T) {
	wantModel := &models.SealabsPay{
		ID:            uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		AccountNumber: "2783993778282621",
		NameOnCard:    "Timothy",
		IsDefault:     true,
	}

	type fields struct {
		sealabsPayRepository *repositories.MockSealabsPayRepository
		redisRepository      *repositories.MockRedisRepository
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args, w *models.SealabsPay)
		want    *models.SealabsPay
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				id: wantModel.ID,
			},
			mock: func(r *fields, a args, w *models.SealabsPay) {
				r.sealabsPayRepository.On("GetById", mock.Anything).Return(wantModel, nil)
			},
			want: wantModel,
		},
		{
			args: args{
				id: wantModel.UserID,
			},
			mock: func(r *fields, a args, w *models.SealabsPay) {
				r.sealabsPayRepository.On("GetById", mock.Anything).Return(nil, e.ErrDataNotFound)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				repositories.NewMockSealabsPayRepository(t),
				repositories.NewMockRedisRepository(t),
			}

			sealabsPayService := services.NewSealabsPayService(&services.SealabsPayServiceConfig{
				SealabsPayRepository: tt.fields.sealabsPayRepository,
				RedisRepository:      tt.fields.redisRepository,
			})

			tt.mock(&tt.fields, tt.args, tt.want)

			got, err := sealabsPayService.GetById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("sealabsPayService.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sealabsPayService.GetById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sealabsPayService_GetByUserIDAndCardNumber(t *testing.T) {
	wantModel := &models.SealabsPay{
		ID:            uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserID:        uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		AccountNumber: "2783993778282621",
		NameOnCard:    "Timothy",
		IsDefault:     true,
	}

	type fields struct {
		sealabsPayRepository *repositories.MockSealabsPayRepository
		redisRepository      *repositories.MockRedisRepository
	}
	type args struct {
		userID uuid.UUID
		cardNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args, w *models.SealabsPay)
		want    *models.SealabsPay
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				userID: wantModel.ID,
				cardNumber: wantModel.AccountNumber,
			},
			mock: func(r *fields, a args, w *models.SealabsPay) {
				r.sealabsPayRepository.On("GetByUserIDAndCardNumber", mock.Anything, mock.Anything).Return(wantModel, nil)
			},
			want: wantModel,
		},
		{
			args: args{
				userID: wantModel.UserID,
				cardNumber: wantModel.AccountNumber,
			},
			mock: func(r *fields, a args, w *models.SealabsPay) {
				r.sealabsPayRepository.On("GetByUserIDAndCardNumber", mock.Anything, mock.Anything).Return(nil, e.ErrDataNotFound)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				repositories.NewMockSealabsPayRepository(t),
				repositories.NewMockRedisRepository(t),
			}

			sealabsPayService := services.NewSealabsPayService(&services.SealabsPayServiceConfig{
				SealabsPayRepository: tt.fields.sealabsPayRepository,
				RedisRepository:      tt.fields.redisRepository,
			})

			tt.mock(&tt.fields, tt.args, tt.want)

			got, err := sealabsPayService.GetByUserIDAndCardNumber(tt.args.userID, tt.args.cardNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("sealabsPayService.GetByUserIDAndCardNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got.CardNumber, tt.want.AccountNumber) {
				t.Errorf("sealabsPayService.GetByUserIDAndCardNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
