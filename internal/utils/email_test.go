package utils_test

import (
	"strings"
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
)

func TestWriteChangePasswordEmailMessage(t *testing.T) {
	type args struct {
		verification_code string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "WriteChangePasswordEmailMessage() success",
			args: args{
				verification_code: "ABCDE-12345",
			},
			want: "ABCDE-12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.WriteChangePasswordEmailMessage(tt.args.verification_code); !strings.Contains(got, tt.want) {
				t.Errorf("WriteChangePasswordEmailMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteResetPasswordEmailMessage(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "WriteResetPasswordEmailMessage() success",
			args: args{
				url: "http://url.com/abcde-12345",
			},
			want: "http://url.com/abcde-12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.WriteResetPasswordEmailMessage(tt.args.url); !strings.Contains(got, tt.want) {
				t.Errorf("WriteResetPasswordEmailMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
