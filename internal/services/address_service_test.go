package services_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
	r "github.com/natanaelb27/oceana-backend/internal/repositories"
	"github.com/natanaelb27/oceana-backend/internal/services"
	"github.com/natanaelb27/oceana-backend/pkg/errors"
)

func Test_addressService_Create(t *testing.T) {
	type fields struct {
		addressRepository  *r.MockAddressRepository
		merchantRepository *r.MockMerchantRepository
	}
	type args struct {
		model *models.Address
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *models.Address
		wantErr bool
	}{
		{
			args: args{
				model: &models.Address{},
			},
			mock: func(r *fields, a args) {
				r.addressRepository.On("GetProfileAddressess", a.model.UserID).Return([]*dtos.AddressDetails{
					{},
					{},
					{},
				}, nil)
				r.addressRepository.On("Create", a.model).Return(&models.Address{}, nil)
			},
			want:    &models.Address{},
			wantErr: false,
		},
		{
			args: args{model: &models.Address{}},
			mock: func(r *fields, a args) {
				r.addressRepository.On("GetProfileAddressess", a.model.UserID).Return([]*dtos.AddressDetails{}, nil)
				r.addressRepository.On("Create", &models.Address{
					IsDefault: true,
				}).Return(&models.Address{}, nil)
			},
			want:    &models.Address{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				addressRepository:  r.NewMockAddressRepository(t),
				merchantRepository: r.NewMockMerchantRepository(t),
			}
			a := services.NewAddressService(&services.AddressServiceConfig{
				AddressRepository:  tt.fields.addressRepository,
				MerchantRepository: tt.fields.merchantRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := a.Create(tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("addressService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addressService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressService_GetProfileAddressess(t *testing.T) {
	type fields struct {
		addressRepository  *r.MockAddressRepository
		merchantRepository *r.MockMerchantRepository
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    []*dtos.AddressDetails
		wantErr bool
	}{
		{
			args: args{
				userID: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.addressRepository.On("GetProfileAddressess", uuid.UUID{}).Return([]*dtos.AddressDetails{}, nil)
			},
			want:    []*dtos.AddressDetails{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				addressRepository:  r.NewMockAddressRepository(t),
				merchantRepository: r.NewMockMerchantRepository(t),
			}
			a := services.NewAddressService(&services.AddressServiceConfig{
				AddressRepository:  tt.fields.addressRepository,
				MerchantRepository: tt.fields.merchantRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := a.GetProfileAddressess(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("addressService.GetProfileAddressess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addressService.GetProfileAddressess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressService_GetProfileAddressById(t *testing.T) {
	type fields struct {
		addressRepository  *r.MockAddressRepository
		merchantRepository *r.MockMerchantRepository
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *dtos.AddressDetails
		wantErr bool
	}{
		{
			args: args{
				userID: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.addressRepository.On("GetProfileAddressById", uuid.UUID{}).Return(&dtos.AddressDetails{}, nil)
			},
			want:    &dtos.AddressDetails{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				addressRepository:  r.NewMockAddressRepository(t),
				merchantRepository: r.NewMockMerchantRepository(t),
			}
			a := services.NewAddressService(&services.AddressServiceConfig{
				AddressRepository:  tt.fields.addressRepository,
				MerchantRepository: tt.fields.merchantRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := a.GetProfileAddressById(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("addressService.GetProfileAddressById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addressService.GetProfileAddressById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressService_UpdateProfileAddress(t *testing.T) {
	type fields struct {
		addressRepository  *r.MockAddressRepository
		merchantRepository *r.MockMerchantRepository
	}
	type args struct {
		userID   uuid.UUID
		newModel *models.Address
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *models.Address
		wantErr bool
	}{
		{
			args: args{
				userID:   uuid.UUID{},
				newModel: &models.Address{},
			},
			mock: func(r *fields, a args) {
				r.addressRepository.On("UpdateProfileAddress", uuid.UUID{}, &models.Address{}).Return(&models.Address{}, nil)
			},
			want:    &models.Address{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				addressRepository:  r.NewMockAddressRepository(t),
				merchantRepository: r.NewMockMerchantRepository(t),
			}
			a := services.NewAddressService(&services.AddressServiceConfig{
				AddressRepository:  tt.fields.addressRepository,
				MerchantRepository: tt.fields.merchantRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := a.UpdateProfileAddress(tt.args.userID, tt.args.newModel)
			if (err != nil) != tt.wantErr {
				t.Errorf("addressService.UpdateProfileAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addressService.UpdateProfileAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressService_DeleteProfileAddress(t *testing.T) {
	type fields struct {
		addressRepository  *r.MockAddressRepository
		merchantRepository *r.MockMerchantRepository
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *models.Address
		wantErr bool
	}{
		{
			args: args{
				id: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.addressRepository.On("GetProfileAddressById", uuid.UUID{}).Return(&dtos.AddressDetails{}, nil)
				r.addressRepository.On("DeleteProfileAddress", uuid.UUID{}).Return(&models.Address{}, nil)
			},
			want:    &models.Address{},
			wantErr: false,
		},
		{
			args: args{
				id: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.addressRepository.On("GetProfileAddressById", uuid.UUID{}).Return(nil, errors.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				id: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.addressRepository.On("GetProfileAddressById", uuid.UUID{}).Return(&dtos.AddressDetails{
					IsDefault: true,
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				id: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.addressRepository.On("GetProfileAddressById", uuid.UUID{}).Return(&dtos.AddressDetails{
					IsShopAddress: true,
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				addressRepository:  r.NewMockAddressRepository(t),
				merchantRepository: r.NewMockMerchantRepository(t),
			}
			a := services.NewAddressService(&services.AddressServiceConfig{
				AddressRepository:  tt.fields.addressRepository,
				MerchantRepository: tt.fields.merchantRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := a.DeleteProfileAddress(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("addressService.DeleteProfileAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addressService.DeleteProfileAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressService_UpdateProfileAddressDefault(t *testing.T) {
	type fields struct {
		addressRepository  *r.MockAddressRepository
		merchantRepository *r.MockMerchantRepository
	}
	type args struct {
		id     uuid.UUID
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *models.Address
		wantErr bool
	}{
		{
			args: args{
				id:     uuid.UUID{},
				userID: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.addressRepository.On("UpdateProfileAddressDefault", uuid.UUID{}, uuid.UUID{}).Return(&models.Address{}, nil)
			},
			want:    &models.Address{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				addressRepository:  r.NewMockAddressRepository(t),
				merchantRepository: r.NewMockMerchantRepository(t),
			}
			a := services.NewAddressService(&services.AddressServiceConfig{
				AddressRepository:  tt.fields.addressRepository,
				MerchantRepository: tt.fields.merchantRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := a.UpdateProfileAddressDefault(tt.args.id, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("addressService.UpdateProfileAddressDefault() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addressService.UpdateProfileAddressDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressService_UpdateProfileAddressShopDefault(t *testing.T) {
	type fields struct {
		addressRepository  *r.MockAddressRepository
		merchantRepository *r.MockMerchantRepository
	}
	type args struct {
		id     uuid.UUID
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    *models.Merchant
		wantErr bool
	}{
		{
			args: args{
				id:     uuid.UUID{},
				userID: uuid.UUID{},
			},
			mock: func(r *fields, a args) {
				r.merchantRepository.On("UpdateProfileAddressShopDefault", uuid.UUID{}, uuid.UUID{}).Return(&models.Merchant{}, nil)
			},
			want:    &models.Merchant{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				addressRepository:  r.NewMockAddressRepository(t),
				merchantRepository: r.NewMockMerchantRepository(t),
			}
			s := services.NewAddressService(&services.AddressServiceConfig{
				AddressRepository:  tt.fields.addressRepository,
				MerchantRepository: tt.fields.merchantRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := s.UpdateProfileAddressShopDefault(tt.args.id, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("addressService.UpdateProfileAddressShopDefault() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addressService.UpdateProfileAddressShopDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
