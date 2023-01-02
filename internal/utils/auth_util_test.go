package utils_test

import (
	"reflect"
	"testing"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/google/uuid"
)

func TestAuthenticateUser(t *testing.T) {
	t.Setenv("JWT_EXPIRATION_HOURS", "1")

	password, _ := utils.Hash("password")

	type args struct {
		user             *models.User
		providedPassword string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "AuthenticateUser() success",
			args: args{
				user: &models.User{
					ID:          uuid.New(),
					FullName:    "",
					Username:    "",
					EmailID:     uuid.New(),
					Email:       models.Email{},
					Password:    password,
					PhoneNumber: "",
					Gender:      nil,
					Birthdate:   nil,
					Photo:       nil,
					GoogleID:    nil,
					Role:        "",
				},
				providedPassword: "password",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := utils.AuthenticateUser(tt.args.user, tt.args.providedPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthenticateUserWithoutEnv(t *testing.T) {
	password, _ := utils.Hash("password")

	type args struct {
		user             *models.User
		providedPassword string
	}
	tests := []struct {
		name    string
		args    args
		want    *internalDtos.UserLoginData
		wantErr bool
	}{
		{
			name: "AuthenticateUser() failed due to wrong password",
			args: args{
				user: &models.User{
					ID:          uuid.New(),
					FullName:    "",
					Username:    "",
					EmailID:     uuid.New(),
					Email:       models.Email{},
					Password:    password,
					PhoneNumber: "",
					Gender:      nil,
					Birthdate:   nil,
					Photo:       nil,
					GoogleID:    nil,
					Role:        "",
				},
				providedPassword: "wrong password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "AuthenticateUser() failed due to error when generating token",
			args: args{
				user: &models.User{
					ID:          uuid.New(),
					FullName:    "",
					Username:    "",
					EmailID:     uuid.New(),
					Email:       models.Email{},
					Password:    password,
					PhoneNumber: "",
					Gender:      nil,
					Birthdate:   nil,
					Photo:       nil,
					GoogleID:    nil,
					Role:        "",
				},
				providedPassword: "password",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.AuthenticateUser(tt.args.user, tt.args.providedPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthenticateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateUserTokenString(t *testing.T) {
	t.Setenv("JWT_EXPIRATION_HOURS", "1")

	password, _ := utils.Hash("password")

	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GenerateUserTokenString() success",
			args: args{
				user: &models.User{
					ID:          uuid.New(),
					FullName:    "",
					Username:    "",
					EmailID:     uuid.New(),
					Email:       models.Email{},
					Password:    password,
					PhoneNumber: "",
					Gender:      nil,
					Birthdate:   nil,
					Photo:       nil,
					GoogleID:    nil,
					Role:        "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := utils.GenerateUserTokenString(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateUserTokenString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGenerateUserTokenStringWithoutEnv(t *testing.T) {
	password, _ := utils.Hash("password")

	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "GenerateUserTokenString() failed due to error when generating token",
			args: args{
				user: &models.User{
					ID:          uuid.New(),
					FullName:    "",
					Username:    "",
					EmailID:     uuid.New(),
					Email:       models.Email{},
					Password:    password,
					PhoneNumber: "",
					Gender:      nil,
					Birthdate:   nil,
					Photo:       nil,
					GoogleID:    nil,
					Role:        "",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.GenerateUserTokenString(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateUserTokenString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateUserTokenString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRefreshTokenString(t *testing.T) {
	t.Setenv("JWT_EXPIRATION_HOURS", "1")
	t.Setenv("REFRESH_TOKEN_DAYS_DURATION", "1")
	t.Setenv("JWT_SECRET_KEY", "secret")
	t.Setenv("JWT_ISSUER", "secret")

	password, _ := utils.Hash("password")

	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GenerateRefreshTokenString() success",
			args: args{
				user: &models.User{
					ID:          uuid.New(),
					FullName:    "",
					Username:    "",
					EmailID:     uuid.New(),
					Email:       models.Email{},
					Password:    password,
					PhoneNumber: "",
					Gender:      nil,
					Birthdate:   nil,
					Photo:       nil,
					GoogleID:    nil,
					Role:        "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := utils.GenerateRefreshTokenString(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRefreshTokenString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGenerateRefreshTokenStringWithoutEnv(t *testing.T) {
	password, _ := utils.Hash("password")

	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "GenerateRefreshTokenString() failed due to error when generating token",
			args: args{
				user: &models.User{
					ID:          uuid.New(),
					FullName:    "",
					Username:    "",
					EmailID:     uuid.New(),
					Email:       models.Email{},
					Password:    password,
					PhoneNumber: "",
					Gender:      nil,
					Birthdate:   nil,
					Photo:       nil,
					GoogleID:    nil,
					Role:        "",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.GenerateRefreshTokenString(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRefreshTokenString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateRefreshTokenString() = %v, want %v", got, tt.want)
			}
		})
	}
}
