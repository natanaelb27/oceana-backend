package utils_test

import (
	"testing"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
)

func TestIsDateMoreThanOneDay(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "IsDateMoreThanOneDay() is more than one day",
			args: args{
				date: time.Now(),
			},
			want: true,
		},
		{
			name: "IsDateMoreThanOneDay() is less than one day",
			args: args{
				date: time.Now().Add(-24 * time.Hour),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.IsDateMoreThanOneDay(tt.args.date); got != tt.want {
				t.Errorf("IsDateMoreThanOneDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPromotionOrVoucherStatus(t *testing.T) {
	type args struct {
		startDatetime time.Time
		endDatetime   time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "PromotionOrVoucherStatus() return upcoming",
			args: args{
				startDatetime: time.Now().Add(1 * time.Hour),
				endDatetime:   time.Now().Add(1 * time.Hour),
			},
			want: constants.UpcomingPromotionOrVoucher,
		},
		{
			name: "PromotionOrVoucherStatus() return ended",
			args: args{
				startDatetime: time.Now().Add(-1 * time.Hour),
				endDatetime:   time.Now().Add(-1 * time.Hour),
			},
			want: constants.EndedPromotionOrVoucher,
		},
		{
			name: "PromotionOrVoucherStatus() return ongoing",
			args: args{
				startDatetime: time.Now().Add(-1 * time.Hour),
				endDatetime:   time.Now().Add(1 * time.Hour),
			},
			want: constants.OngoingPromotionOrVoucher,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.PromotionOrVoucherStatus(tt.args.startDatetime, tt.args.endDatetime); got != tt.want {
				t.Errorf("PromotionOrVoucherStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
