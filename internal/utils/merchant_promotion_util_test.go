package utils_test

import (
	"testing"
	"time"

	"github.com/natanaelb27/oceana-backend/internal/constants"
	"github.com/natanaelb27/oceana-backend/internal/utils"
)

func TestValidateMpCreateDatetime(t *testing.T) {
	type args struct {
		startDatetime time.Time
		endDatetime   time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ValidateMpCreateDatetime() success",
			args: args{
				startDatetime: time.Now().Add(36 * time.Hour),
				endDatetime:   time.Now().Add(37 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "ValidateMpCreateDatetime() failed due to start time before now",
			args: args{
				startDatetime: time.Now().Add(-1 * time.Hour),
				endDatetime:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "ValidateMpCreateDatetime() failed due to end time before start time",
			args: args{
				startDatetime: time.Now().Add(1 * time.Hour),
				endDatetime:   time.Now().Add(-1 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "ValidateMpCreateDatetime() failed due to start time less than one day",
			args: args{
				startDatetime: time.Now().Add(1 * time.Hour),
				endDatetime:   time.Now().Add(1 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "ValidateMpCreateDatetime() failed due to duration less than hour",
			args: args{
				startDatetime: time.Now().Add(36 * time.Hour),
				endDatetime:   time.Now().Add(36 * time.Hour).Add(30 * time.Minute),
			},
			wantErr: true,
		},
		{
			name: "ValidateMpCreateDatetime() failed due to end date more than 90 days",
			args: args{
				startDatetime: time.Now().Add(36 * time.Hour),
				endDatetime:   time.Now().Add(24 * 100 * time.Hour),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := utils.ValidateMpCreateDatetime(tt.args.startDatetime, tt.args.endDatetime); (err != nil) != tt.wantErr {
				t.Errorf("ValidateMpCreateDatetime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMpUpdateDatetime(t *testing.T) {
	type args struct {
		startDatetime        time.Time
		endDatetime          time.Time
		currentStartDatetime time.Time
		currentEndDatetime   time.Time
		promotionStatus      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ValidateMpUpdateDatetime() success",
			args: args{
				startDatetime:        time.Now().Add(11 * time.Hour),
				endDatetime:          time.Now().Add(12 * time.Hour),
				currentStartDatetime: time.Now(),
				currentEndDatetime:   time.Now().Add(24 * time.Hour),
				promotionStatus:      constants.UpcomingPromotionOrVoucher,
			},
			wantErr: false,
		},
		{
			name: "ValidateMpUpdateDatetime() failed due to ended status",
			args: args{
				startDatetime:        time.Now().Add(11 * time.Hour),
				endDatetime:          time.Now().Add(12 * time.Hour),
				currentStartDatetime: time.Now(),
				currentEndDatetime:   time.Now().Add(24 * time.Hour),
				promotionStatus:      constants.EndedPromotionOrVoucher,
			},
			wantErr: true,
		},
		{
			name: "ValidateMpUpdateDatetime() failed due to ongoing status",
			args: args{
				startDatetime:        time.Now().Add(11 * time.Hour),
				endDatetime:          time.Now().Add(12 * time.Hour),
				currentStartDatetime: time.Now(),
				currentEndDatetime:   time.Now().Add(24 * time.Hour),
				promotionStatus:      constants.OngoingPromotionOrVoucher,
			},
			wantErr: true,
		},
		{
			name: "ValidateMpUpdateDatetime() failed due new start/end date is not between previous dates",
			args: args{
				startDatetime:        time.Now().Add(11 * time.Hour),
				endDatetime:          time.Now().Add(25 * time.Hour),
				currentStartDatetime: time.Now(),
				currentEndDatetime:   time.Now().Add(24 * time.Hour),
				promotionStatus:      constants.UpcomingPromotionOrVoucher,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := utils.ValidateMpUpdateDatetime(tt.args.startDatetime, tt.args.endDatetime, tt.args.currentStartDatetime, tt.args.currentEndDatetime, tt.args.promotionStatus); (err != nil) != tt.wantErr {
				t.Errorf("ValidateMpUpdateDatetime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMpProductIDList(t *testing.T) {
	type args struct {
		mpRequestProductIDList  []int
		pmpRequestProductIDList []int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ValidateMpProductIDList() success",
			args: args{
				mpRequestProductIDList:  []int{1, 2, 3},
				pmpRequestProductIDList: []int{1, 2, 3},
			},
			wantErr: false,
		},
		{
			name: "ValidateMpProductIDList() failed due to different length",
			args: args{
				mpRequestProductIDList:  []int{1, 2, 3, 4},
				pmpRequestProductIDList: []int{1, 2, 3},
			},
			wantErr: true,
		},
		{
			name: "ValidateMpProductIDList() failed due to different data",
			args: args{
				mpRequestProductIDList:  []int{1, 2, 4},
				pmpRequestProductIDList: []int{1, 2, 3},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := utils.ValidateMpProductIDList(tt.args.mpRequestProductIDList, tt.args.pmpRequestProductIDList); (err != nil) != tt.wantErr {
				t.Errorf("ValidateMpProductIDList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
