package services_test

import (
	"reflect"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_categoryService_GetAll(t *testing.T) {
	type fields struct {
		categoryRepository *r.MockCategoryRepository
	}
	type args struct {
		query *externalDtos.CategoryQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    []*dtos.Category
		wantErr bool
	}{
		{
			args: args{
				query: &externalDtos.CategoryQuery{},
			},
			mock: func(r *fields, a args) {
				r.categoryRepository.On("GetAllCategories").Return([]*models.Category{}, nil)
			},
			want:    []*dtos.Category{},
			wantErr: false,
		},
		{
			args: args{
				query: &externalDtos.CategoryQuery{},
			},
			mock: func(r *fields, a args) {
				r.categoryRepository.On("GetAllCategories").Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				query: &externalDtos.CategoryQuery{
					CategoryID: 1,
				},
			},
			mock: func(r *fields, a args) {
				r.categoryRepository.On("GetAllByID", mock.Anything).Return([]*models.Category{}, nil).Once()
				r.categoryRepository.On("GetByID", mock.Anything).Return(&models.Category{
					ParentID: 1,
				}, nil).Once()
				r.categoryRepository.On("GetByID", mock.Anything).Return(nil, e.ErrDataNotFound).Once()
			},
			want:    []*dtos.Category{},
			wantErr: false,
		},
		{
			args: args{
				query: &externalDtos.CategoryQuery{
					CategoryID: 1,
				},
			},
			mock: func(r *fields, a args) {
				r.categoryRepository.On("GetAllByID", mock.Anything).Return([]*models.Category{}, nil).Once()
				r.categoryRepository.On("GetByID", mock.Anything).Return(&models.Category{}, nil).Once()
			},
			want:    []*dtos.Category{},
			wantErr: false,
		},
		{
			args: args{
				query: &externalDtos.CategoryQuery{
					CategoryID: 1,
				},
			},
			mock: func(r *fields, a args) {
				r.categoryRepository.On("GetAllByID", mock.Anything).Return([]*models.Category{}, nil).Once()
				r.categoryRepository.On("GetByID", mock.Anything).Return(&models.Category{
					ParentID: 1,
				}, nil).Once()
				r.categoryRepository.On("GetByID", mock.Anything).Return(nil, e.ErrUnknown).Once()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				categoryRepository: r.NewMockCategoryRepository(t),
			}
			s := services.NewCategoryService(&services.CategoryServiceConfig{
				CategoryRepository: tt.fields.categoryRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := s.GetAll(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("categoryService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("categoryService.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_categoryService_GetAncestorsByID(t *testing.T) {
	type fields struct {
		categoryRepository *r.MockCategoryRepository
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    []*externalDtos.GetCategoriesResponse
		wantErr bool
	}{
		{
			args: args{
				id: 1,
			},
			mock: func(r *fields, a args) {
				r.categoryRepository.On("GetAncestorsByID", mock.Anything).Return([]*models.Category{}, nil)
			},
			want:    externalDtos.NewGetCategoriesResponse([]*models.Category{}),
			wantErr: false,
		},
		{
			args: args{
				id: 1,
			},
			mock: func(r *fields, a args) {
				r.categoryRepository.On("GetAncestorsByID", mock.Anything).Return(nil, e.ErrUnknown)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				categoryRepository: r.NewMockCategoryRepository(t),
			}
			s := services.NewCategoryService(&services.CategoryServiceConfig{
				CategoryRepository: tt.fields.categoryRepository,
			})
			tt.mock(&tt.fields, tt.args)
			got, err := s.GetAncestorsByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("categoryService.GetAncestorsByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("categoryService.GetAncestorsByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_categoryService_GetMerchantCategories(t *testing.T) {
	type fields struct {
		categoryRepository *r.MockCategoryRepository
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(r *fields, a args)
		want    []*externalDtos.GetCategoriesResponse
		wantErr bool
	}{
		{
			args: args{
				id: 1,
			},
			mock: func(r *fields, a args) {
				r.categoryRepository.On("GetMerchantCategoriesByMerchantID", mock.Anything).Return([]*models.Category{}, nil)
			},
			want:    externalDtos.NewGetCategoriesResponse([]*models.Category{}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields = fields{
				categoryRepository: r.NewMockCategoryRepository(t),
			}
			s := services.NewCategoryService(&services.CategoryServiceConfig{
				CategoryRepository: tt.fields.categoryRepository,
			})
			tt.mock(&tt.fields, tt.args)
			_, err := s.GetMerchantCategories(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("categoryService.GetMerchantCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
