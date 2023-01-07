package services_test

import (
	"reflect"
	"testing"

	"github.com/natanaelb27/oceana-backend/internal/models"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/services"
	"github.com/natanaelb27/oceana-backend/pkg/dtos"
	"github.com/natanaelb27/oceana-backend/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_cityService_GetAll(t *testing.T) {
	type fields struct {
		cityRepository *r.MockCityRepository
	}
	type args struct {
		query *dtos.CitiesQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    []*dtos.CityResponse
		wantErr bool
	}{
		{
			args: args{
				query: &dtos.CitiesQuery{},
			},
			mock: func(r *fields, a args) {
				r.cityRepository.On("CitiesGetAll", mock.Anything).Return([]*models.City{}, nil)
			},
			want:    dtos.NewCitiesResponse([]*models.City{}),
			wantErr: false,
		},
		{
			args: args{
				query: &dtos.CitiesQuery{},
			},
			mock: func(r *fields, a args) {
				r.cityRepository.On("CitiesGetAll", mock.Anything).Return(nil, errors.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				cityRepository: r.NewMockCityRepository(t),
			}
			s := services.NewCityService(&services.CityServiceConfig{
				CityRepository: tt.fields.cityRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := s.GetAll(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("cityService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cityService.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
