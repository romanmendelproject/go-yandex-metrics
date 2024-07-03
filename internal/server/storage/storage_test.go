package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemStorage_SetGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "God",
			args: args{
				name:  "Test",
				value: 0.5,
			},
			want: 0.5,
		},
	}

	storage := NewMemStorage("test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.SetGauge(tt.args.name, tt.args.value)
			require.Equal(t, storage.gauge["Test"], tt.want)
		})
	}
}

func TestMemStorage_SetCount(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "God",
			args: args{
				name:  "Test",
				value: 1,
			},
			want: 1,
		},
	}

	storage := NewMemStorage("test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.SetCounter(tt.args.name, tt.args.value)
			require.Equal(t, storage.counter["Test"], tt.want)
		})
	}
}
