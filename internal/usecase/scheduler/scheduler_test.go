package scheduler

import (
	"testing"
	"time"
)

func TestNextDate(t *testing.T) {
	now := time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC) // Понедельник

	tests := []struct {
		name    string
		dateStr string
		repeat  string
		want    string
		wantErr bool
	}{
		{
			name:    "Daily - next day",
			dateStr: "20260415",
			repeat:  "daily 1",
			want:    "20260416",
		},
		{
			name:    "Daily - after 5 days",
			dateStr: "20260415",
			repeat:  "daily 4",
			want:    "20260420",
		},
		{
			name:    "Even - from odd day",
			dateStr: "20260415",
			repeat:  "even",
			want:    "20260416",
		},
		{
			name:    "Even - from even day (should be next even)",
			dateStr: "20260415",
			repeat:  "even",
			want:    "20260417",
		},
		{
			name:    "Specific dates - pick nearest",
			dateStr: "20260413",
			repeat:  "specific 20260509,20260418,20260412",
			want:    "20260416",
		},
		{
			name:    "Incorrect repeat type",
			dateStr: "20260413",
			repeat:  "unknown 1",
			wantErr: true,
		},
		{
			name:    "Monthly - day 31 in short months",
			dateStr: "20260131",
			repeat:  "monthly 31",
			want: "20260531",
		},
		{
			name:    "Daily - jump from past to future",
			dateStr: "20200101", 
			repeat:  "daily 1",
			want:    "20260414",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextDate(now, tt.dateStr, tt.repeat)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NextDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
