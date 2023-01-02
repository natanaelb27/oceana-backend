package services_test

import (
	"testing"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
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
