package storage

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	storage := NewMemStorage("test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.SetGauge(ctx, tt.args.name, tt.args.value)
			v, ok := storage.gauge.Load("Test")
			if ok {
				require.Equal(t, v, tt.want)
			}

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	storage := NewMemStorage("test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.SetCounter(ctx, tt.args.name, tt.args.value)
			v, ok := storage.counter.Load("Test")
			if ok {
				require.Equal(t, v, tt.want)
			}
		})
	}
}
