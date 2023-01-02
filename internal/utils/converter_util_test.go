package utils_test

import (
	"testing"

	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
)

func TestStringToUint(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name       string
		args       args
		wantTarget uint
	}{
		{
			name: "StringToUint() success",
			args: args{
				source: "1",
			},
			wantTarget: 1,
		},
		{
			name: "StringToUint() failed due to less than 0",
			args: args{
				source: "-1",
			},
			wantTarget: 0,
		},
		{
			name: "StringToUint() failed due to not a number",
			args: args{
				source: "string",
			},
			wantTarget: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTarget := utils.StringToUint(tt.args.source); gotTarget != tt.wantTarget {
				t.Errorf("StringToUint() = %v, want %v", gotTarget, tt.wantTarget)
			}
		})
	}
}
