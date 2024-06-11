package handlers

import (
	"reflect"
	"testing"
)

func TestParseURLUpdate(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    URLParamsUpdate
		wantErr bool
	}{
		{
			name: "Good url gauge",
			args: args{
				url: "/update/gauge/test/0.1",
			},
			want: URLParamsUpdate{
				metricType:  "gauge",
				metricName:  "test",
				metricValue: "0.1",
			},
			wantErr: false,
		},
		{
			name: "Good url count",
			args: args{
				url: "/update/gauge/test/1",
			},
			want: URLParamsUpdate{
				metricType:  "gauge",
				metricName:  "test",
				metricValue: "1",
			},
			wantErr: false,
		},
		{
			name: "Bad (too many arguments)",
			args: args{
				url: "/update/gauge/test/test/1",
			},
			want:    URLParamsUpdate{},
			wantErr: true,
		},
		{
			name: "Bad (too few arguments)",
			args: args{
				url: "/update/gauge/1",
			},
			want:    URLParamsUpdate{},
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
