package services_test

import (
	"reflect"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/internal/mocks"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func Test_walletService_Activate(t *testing.T) {
	modelBefore := &models.Wallet{
		ID:       uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserId:   uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		IsActive: false,
	}
	model := &models.Wallet{
		ID:       uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
		UserId:   uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
		IsActive: true,
	}

	type fields struct {
		walletRepository *mocks.WalletRepository
	}
	type args struct {
		args1 uuid.UUID
		args2 *models.Wallet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args, w *models.Wallet)
		want    *models.Wallet
		wantErr bool
	}{
		{
			args: args{
				args2: model,
			},
			mock: func(r *fields, a args, w *models.Wallet) {
				r.walletRepository.On("GetByUserId", mock.Anything).Return(modelBefore, nil)
				r.walletRepository.On("Activate", mock.Anything, mock.Anything).Return(model, nil)
			},
			want: model,
		},
		{
			args: args{
				args2: model,
			},
			mock: func(r *fields, a args, w *models.Wallet) {
				r.walletRepository.On("GetByUserId", mock.Anything).Return(nil, e.ErrDatabaseOperationFailed)
			},
			wantErr: true,
		},
		{
			args: args{
				args2: &models.Wallet{
					ID:       uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b851"),
					UserId:   uuid.MustParse("e56a013f-58c3-4dd3-93d1-190d6a46b852"),
					IsActive: true,
				},
			},
			mock: func(r *fields, a args, w *models.Wallet) {
				r.walletRepository.On("GetByUserId", mock.Anything).Return(a.args2, nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				mocks.NewWalletRepository(t),
			}

			walletService := services.NewWalletService(&services.WalletServiceConfig{
				WalletRepository: tt.fields.walletRepository,
			})

			tt.mock(&tt.fields, tt.args, tt.want)

			got, err := walletService.Activate(tt.args.args1, tt.args.args2)
			if (err != nil) != tt.wantErr {
				t.Errorf("walletService.Activate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("walletService.Activate() = %v, want %v", got, tt.want)
			}
		})
	}
}
