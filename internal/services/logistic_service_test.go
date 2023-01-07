package services_test

import (
	"testing"

	"github.com/google/uuid"
	logistics_dtos "github.com/natanaelb27/oceana-backend/internal/dtos/logistic_dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/services"
	"github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_logisticService_GetByMerchantId(t *testing.T) {
	type fields struct {
		courierRepository                *r.MockCourierRepository
		courierServiceRepository         *r.MockCourierServiceRepository
		merchantCourierServiceRepository *r.MockMerchantCourierServiceRepository
		userRepository                   *r.MockUserRepository
		merchantRepository               *r.MockMerchantRepository
		merchantConfigurationRepository  *r.MockMerchantConfigurationRepository
		productRepository                *r.MockProductRepository
	}
	type args struct {
		userID uuid.UUID
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
				userID: uuid.UUID{},
			},
			mock: func(r *fields) {
				r.merchantRepository.On("GetByUserID", mock.Anything).Return(&models.Merchant{}, nil)
				r.courierRepository.On("GetAll").Return([]*models.Courier{
					{},
				}, nil)
				r.merchantCourierServiceRepository.On("GetAllByCourierID", mock.Anything, mock.Anything).Return([]*logistics_dtos.CourierService{
					{
						IsSelected: true,
					},
				}, nil)
				r.merchantConfigurationRepository.On("GetByMerchantID", mock.Anything).Return(&models.MerchantConfiguration{}, nil)
			},
			wantErr: false,
		},
		{
			args: args{
				userID: uuid.UUID{},
			},
			mock: func(r *fields) {
				r.merchantRepository.On("GetByUserID", mock.Anything).Return(&models.Merchant{}, nil)
				r.courierRepository.On("GetAll").Return([]*models.Courier{
					{},
				}, nil)
				r.merchantCourierServiceRepository.On("GetAllByCourierID", mock.Anything, mock.Anything).Return([]*logistics_dtos.CourierService{
					{
						IsSelected: true,
					},
				}, nil)
				r.merchantConfigurationRepository.On("GetByMerchantID", mock.Anything).Return(nil, errors.ErrUnknown)
			},
			wantErr: true,
		},
		{
			args: args{
				userID: uuid.UUID{},
			},
			mock: func(r *fields) {
				r.merchantRepository.On("GetByUserID", mock.Anything).Return(&models.Merchant{}, nil)
				r.courierRepository.On("GetAll").Return([]*models.Courier{
					{},
				}, nil)
				r.merchantCourierServiceRepository.On("GetAllByCourierID", mock.Anything, mock.Anything).Return([]*logistics_dtos.CourierService{
					{
						IsSelected: true,
					},
				}, nil)
				r.merchantConfigurationRepository.On("GetByMerchantID", mock.Anything).Return(nil, errors.ErrUnknown)
			},
			wantErr: true,
		},
		{
			args: args{
				userID: uuid.UUID{},
			},
			mock: func(r *fields) {
				r.merchantRepository.On("GetByUserID", mock.Anything).Return(&models.Merchant{}, nil)
				r.courierRepository.On("GetAll").Return([]*models.Courier{
					{},
				}, nil)
				r.merchantCourierServiceRepository.On("GetAllByCourierID", mock.Anything, mock.Anything).Return(nil, errors.ErrUnknown)
			},
			wantErr: true,
		},
		{
			args: args{
				userID: uuid.UUID{},
			},
			mock: func(r *fields) {
				r.merchantRepository.On("GetByUserID", mock.Anything).Return(&models.Merchant{}, nil)
				r.courierRepository.On("GetAll").Return(nil, errors.ErrUnknown)
			},
			wantErr: true,
		},
		{
			args: args{
				userID: uuid.UUID{},
			},
			mock: func(r *fields) {
				r.merchantRepository.On("GetByUserID", mock.Anything).Return(&models.Merchant{}, errors.ErrUnknown)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				courierRepository:                r.NewMockCourierRepository(t),
				courierServiceRepository:         r.NewMockCourierServiceRepository(t),
				merchantCourierServiceRepository: r.NewMockMerchantCourierServiceRepository(t),
				userRepository:                   r.NewMockUserRepository(t),
				merchantRepository:               r.NewMockMerchantRepository(t),
				merchantConfigurationRepository:  r.NewMockMerchantConfigurationRepository(t),
				productRepository:                r.NewMockProductRepository(t),
			}
			s := services.NewLogisticService(&services.LogisticServiceConfig{
				CourierRepository:                tt.fields.courierRepository,
				CourierServiceRepository:         tt.fields.courierServiceRepository,
				MerchantCourierServiceRepository: tt.fields.merchantCourierServiceRepository,
				UserRepository:                   tt.fields.userRepository,
				MerchantRepository:               tt.fields.merchantRepository,
				MerchantConfigurationRepository:  tt.fields.merchantConfigurationRepository,
				ProductRepository:                tt.fields.productRepository,
			})
			tt.mock(&tt.fields)
			_, err := s.GetByMerchantId(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("logisticService.GetByMerchantId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_logisticService_ProductGetByMerchantID(t *testing.T) {
	type fields struct {
		courierRepository                *r.MockCourierRepository
		courierServiceRepository         *r.MockCourierServiceRepository
		merchantCourierServiceRepository *r.MockMerchantCourierServiceRepository
		userRepository                   *r.MockUserRepository
		merchantRepository               *r.MockMerchantRepository
		merchantConfigurationRepository  *r.MockMerchantConfigurationRepository
		productRepository                *r.MockProductRepository
	}
	type args struct {
		merchantID uint
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
				merchantID: 1,
			},
			mock: func(r *fields) {
				r.merchantRepository.On("GetByMerchantIDJoinUser", mock.Anything).Return(&models.Merchant{}, nil)
				r.merchantRepository.On("GetByUserID", mock.Anything).Return(&models.Merchant{}, nil)
				r.courierRepository.On("GetAll").Return([]*models.Courier{
					{},
				}, nil)
				r.merchantCourierServiceRepository.On("GetAllByCourierID", mock.Anything, mock.Anything).Return([]*logistics_dtos.CourierService{
					{
						IsSelected: true,
					},
				}, nil)
				r.merchantConfigurationRepository.On("GetByMerchantID", mock.Anything).Return(&models.MerchantConfiguration{}, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				courierRepository:                r.NewMockCourierRepository(t),
				courierServiceRepository:         r.NewMockCourierServiceRepository(t),
				merchantCourierServiceRepository: r.NewMockMerchantCourierServiceRepository(t),
				userRepository:                   r.NewMockUserRepository(t),
				merchantRepository:               r.NewMockMerchantRepository(t),
				merchantConfigurationRepository:  r.NewMockMerchantConfigurationRepository(t),
				productRepository:                r.NewMockProductRepository(t),
			}
			s := services.NewLogisticService(&services.LogisticServiceConfig{
				CourierRepository:                tt.fields.courierRepository,
				CourierServiceRepository:         tt.fields.courierServiceRepository,
				MerchantCourierServiceRepository: tt.fields.merchantCourierServiceRepository,
				UserRepository:                   tt.fields.userRepository,
				MerchantRepository:               tt.fields.merchantRepository,
				MerchantConfigurationRepository:  tt.fields.merchantConfigurationRepository,
				ProductRepository:                tt.fields.productRepository,
			})
			tt.mock(&tt.fields)
			_, err := s.ProductGetByMerchantID(tt.args.merchantID)
			if (err != nil) != tt.wantErr {
				t.Errorf("logisticService.ProductGetByMerchantID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
