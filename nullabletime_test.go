package interval

import (
	"reflect"
	"testing"
	"time"
)

func TestParseNullableTimeInterval(t *testing.T) {
	type args struct {
		intervalStr string
		layout      []string
	}
	tm := time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		args    args
		wantTi  *NullableTimeInterval
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				intervalStr: "[null,2022-10-01T00:00:00Z]",
			},
			wantTi: NewNullableTimeInterval(nil, &tm, Closed),
		},
		{
			name: "2",
			args: args{
				intervalStr: "[null,null)",
			},
			wantTi: NewNullableTimeInterval(nil, nil, ClosedOpen),
		},
		{
			name: "2",
			args: args{
				intervalStr: "(2022-10-01 00:00:00,null)",
				layout:      []string{"2006-01-02 15:04:05"},
			},
			wantTi: NewNullableTimeInterval(&tm, nil, Open),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseNullableTimeInterval(tt.args.intervalStr, tt.args.layout...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNullableTimeInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseNullableTimeInterval() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}

func TestNullableTimeInterval_Contains(t *testing.T) {
	ti0, err := ParseNullableTimeInterval("[2022-10-01T00:00:00Z, null]")
	if err != nil {
		panic(err)
	}
	ti1, err := ParseNullableTimeInterval("(2022-10-01T00:00:00Z, null]")
	if err != nil {
		panic(err)
	}
	ti2, err := ParseNullableTimeInterval("(null, 2022-10-01T00:00:00Z]")
	if err != nil {
		panic(err)
	}
	ti3, err := ParseNullableTimeInterval("(null, 2022-10-01T00:00:00Z)")
	if err != nil {
		panic(err)
	}
	ti4, err := ParseNullableTimeInterval("(2022-09-01T00:00:00Z, 2022-10-01T00:00:00Z)")
	if err != nil {
		panic(err)
	}

	tm0 := time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC)
	tm1 := time.Date(2022, 9, 30, 0, 0, 0, 0, time.UTC)

	type args struct {
		e *time.Time
	}
	tests := []struct {
		name string
		ti   *NullableTimeInterval
		args args
		want bool
	}{
		{
			name: "0",
			ti:   ti0,
			args: args{&tm0},
			want: true,
		},
		{
			name: "1",
			ti:   ti0,
			args: args{&tm1},
			want: false,
		},
		{
			name: "2",
			ti:   ti1,
			args: args{&tm1},
			want: false,
		},
		{
			name: "3",
			ti:   ti1,
			args: args{&tm0},
			want: false,
		},
		{
			name: "4",
			ti:   ti1,
			args: args{&tm1},
			want: false,
		},
		{
			name: "5",
			ti:   ti2,
			args: args{&tm0},
			want: true,
		},
		{
			name: "6",
			ti:   ti2,
			args: args{&tm1},
			want: true,
		},
		{
			name: "7",
			ti:   ti3,
			args: args{&tm0},
			want: false,
		},
		{
			name: "8",
			ti:   ti3,
			args: args{&tm1},
			want: true,
		},
		{
			name: "9",
			ti:   ti4,
			args: args{&tm0},
			want: false,
		},
		{
			name: "10",
			ti:   ti4,
			args: args{&tm1},
			want: true,
		},
		{
			name: "11",
			ti:   ti4,
			args: args{nil},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ti.Contains(tt.args.e); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			} else {
				t.Logf("Contains() = %v, interval = %s", got, tt.ti.String())
			}
		})
	}
}
