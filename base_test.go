package interval

import (
	"math"
	"reflect"
	"testing"
)

func TestNewBaseInterval(t *testing.T) {
	i1 := NewBaseInterval[int](1, 2)
	t.Log(i1)
	t.Log(i1.Contains(2))
	t.Log(i1)
}

func TestParseIntInterval(t *testing.T) {
	type args struct {
		intervalStr string
	}
	tests := []struct {
		name    string
		args    args
		wantI   *BaseInterval[int64]
		wantErr bool
	}{
		{
			name:  "1",
			args:  args{intervalStr: "[1,2]"},
			wantI: NewBaseInterval[int64](1, 2, Closed),
		},
		{
			name:  "2",
			args:  args{intervalStr: "[1,2)"},
			wantI: NewBaseInterval[int64](1, 2, ClosedOpen),
		},
		{
			name:  "3",
			args:  args{intervalStr: "[1,  20000)"},
			wantI: NewBaseInterval[int64](1, 20000, ClosedOpen),
		},
		{
			name:  "4",
			args:  args{intervalStr: "(0,  20000)"},
			wantI: NewBaseInterval[int64](0, 20000, Open),
		},
		{
			name:  "5",
			args:  args{intervalStr: "(-1,  20000)"},
			wantI: NewBaseInterval[int64](-1, 20000, Open),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseIntInterval(tt.args.intervalStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIntInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantI) {
				t.Errorf("ParseIntInterval() gotTi = %v, want %v", gotTi, tt.wantI)
			}
		})
	}
}

func TestParseFloatInterval(t *testing.T) {
	type args struct {
		intervalStr string
	}
	tests := []struct {
		name    string
		args    args
		wantI   *BaseInterval[float64]
		wantErr bool
	}{
		{
			name:  "1",
			args:  args{intervalStr: "[1.1,2.0]"},
			wantI: NewBaseInterval[float64](1.1, 2.0, Closed),
		},
		{
			name:  "2",
			args:  args{intervalStr: "[1.1,2)"},
			wantI: NewBaseInterval[float64](1.1, 2, ClosedOpen),
		},
		{
			name:  "3",
			args:  args{intervalStr: "[1.001,  20000)"},
			wantI: NewBaseInterval[float64](1.001, 20000, ClosedOpen),
		},
		{
			name:  "4",
			args:  args{intervalStr: "(0.001,  20000)"},
			wantI: NewBaseInterval[float64](0.001, 20000, Open),
		},
		{
			name:  "5",
			args:  args{intervalStr: "(-0.1,  Inf)"},
			wantI: NewBaseInterval[float64](-0.1, math.Inf(1), Open),
		},
		{
			name:  "6",
			args:  args{intervalStr: "(-Inf,  +Inf)"},
			wantI: NewBaseInterval[float64](math.Inf(-1), math.Inf(1), Open),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotI, err := ParseFloatInterval(tt.args.intervalStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFloatInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotI, tt.wantI) {
				t.Errorf("ParseFloatInterval() gotI = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}
