package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURLUpdate(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    URLParams
		wantErr bool
	}{
		{
			name: "Good url gauge",
			args: args{
				url: "/update/gauge/test/0.1",
			},
			want: URLParams{
				MetricType:  "gauge",
				MetricName:  "test",
				MetricValue: "0.1",
			},
			wantErr: false,
		},
		{
			name: "Good url count",
			args: args{
				url: "/update/gauge/test/1",
			},
			want: URLParams{
				MetricType:  "gauge",
				MetricName:  "test",
				MetricValue: "1",
			},
			wantErr: false,
		},
		{
			name: "Bad (too many arguments)",
			args: args{
				url: "/update/gauge/test/test/1",
			},
			want:    URLParams{},
			wantErr: true,
		},
		{
			name: "Bad (too few arguments)",
			args: args{
				url: "/update/gauge/1",
			},
			want:    URLParams{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURLUpdate(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseURLValue(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    URLParams
		wantErr bool
	}{
		{
			name: "Good url gauge",
			args: args{
				url: "/update/gauge/test",
			},
			want: URLParams{
				MetricType: "gauge",
				MetricName: "test",
			},
			wantErr: false,
		},
		{
			name: "Good url count",
			args: args{
				url: "/update/gauge/test",
			},
			want: URLParams{
				MetricType: "gauge",
				MetricName: "test",
			},
			wantErr: false,
		},
		{
			name: "Bad (too many arguments)",
			args: args{
				url: "/update/gauge/test/test/1",
			},
			want:    URLParams{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURLValue(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURLValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseURLValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFloatPtr(t *testing.T) {
	value := 5.1
	address := &value
	valuePtr := GetFloatPtr(value)

	assert.Equal(t, address, valuePtr)
}
