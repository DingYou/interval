package interval

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTimeInterval(t *testing.T) {
	type args struct {
		intervalStr string
		layout      []string
	}
	tests := []struct {
		name    string
		args    args
		wantTi  *TimeInterval
		wantErr bool
	}{
		{
			name: "standardTimeClosed",
			args: args{
				intervalStr: "[2022-01-01T13:12:11Z, 2022-05-01T13:12:11Z]",
			},
			wantTi: NewTimeInterval(
				time.Date(2022, 1, 1, 13, 12, 11, 0, time.UTC),
				time.Date(2022, 5, 1, 13, 12, 11, 0, time.UTC),
				Closed,
			),
		},
		{
			name: "standardTimeOpen",
			args: args{
				intervalStr: "(2022-01-01T13:12:11Z, 2022-05-01T13:12:11Z)",
			},
			wantTi: NewTimeInterval(
				time.Date(2022, 1, 1, 13, 12, 11, 0, time.UTC),
				time.Date(2022, 5, 1, 13, 12, 11, 0, time.UTC),
				Open,
			),
		},
		{
			name: "standardTimeOpenClosed",
			args: args{
				intervalStr: "(2022-01-01T13:12:11Z, 2022-05-01T13:12:11Z]",
			},
			wantTi: NewTimeInterval(
				time.Date(2022, 1, 1, 13, 12, 11, 0, time.UTC),
				time.Date(2022, 5, 1, 13, 12, 11, 0, time.UTC),
				OpenClosed,
			),
		},
		{
			name: "standardTimeClosedOpen",
			args: args{
				intervalStr: "[2022-01-01T13:12:11Z, 2022-05-01T13:12:11Z)",
			},
			wantTi: NewTimeInterval(
				time.Date(2022, 1, 1, 13, 12, 11, 0, time.UTC),
				time.Date(2022, 5, 1, 13, 12, 11, 0, time.UTC),
				ClosedOpen,
			),
		},
		{
			name: "standardDateClosed",
			args: args{
				intervalStr: "[2022-01-01, 2022-05-01)",
				layout:      []string{"2006-01-02"},
			},
			wantTi: NewTimeInterval(
				time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC),
				ClosedOpen,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseTimeInterval(tt.args.intervalStr, tt.args.layout...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseTimeInterval() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}
