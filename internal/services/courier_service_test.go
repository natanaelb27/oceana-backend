package services_test

import (
	"testing"

	internalDtos "github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/services"
	"github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_courierService_GetAll(t *testing.T) {
	type fields struct {
		courierRepository        *r.MockCourierRepository
		courierServiceRepository *r.MockCourierServiceRepository
	}
	tests := []struct {
		name    string
		fields  fields
		mock    func(r *fields)
		want    []*internalDtos.Courier
		wantErr bool
	}{
		{
			mock: func(r *fields) {
				r.courierRepository.On("GetAll").Return([]*models.Courier{{}}, nil)
				r.courierServiceRepository.On("GetAllByCourierID", mock.Anything).Return([]*internalDtos.CourierService{}, nil)
			},
			want:    []*internalDtos.Courier{{}},
			wantErr: false,
		},
		{
			mock: func(r *fields) {
				r.courierRepository.On("GetAll").Return([]*models.Courier{{}}, nil)
				r.courierServiceRepository.On("GetAllByCourierID", mock.Anything).Return(nil, errors.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			mock: func(r *fields) {
				r.courierRepository.On("GetAll").Return(nil, errors.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				courierRepository:        r.NewMockCourierRepository(t),
				courierServiceRepository: r.NewMockCourierServiceRepository(t),
			}
			s := services.NewCourierService(&services.CourierServiceConfig{
				CourierRepository:        tt.fields.courierRepository,
				CourierServiceRepository: tt.fields.courierServiceRepository,
			})
			tt.mock(&tt.fields)
			_, err := s.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("courierService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
