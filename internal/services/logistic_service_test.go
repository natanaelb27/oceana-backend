package services_test

import (
	"testing"

	logistics_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos/logistic_dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
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
