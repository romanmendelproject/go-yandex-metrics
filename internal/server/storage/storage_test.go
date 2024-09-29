package storage

import (
	"context"
	"testing"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/metrics"
	"github.com/stretchr/testify/assert"
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

func TestMemStorage_GetCount(t *testing.T) {
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

	ctx := context.Background()
	storage := NewMemStorage("")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.SetCounter(ctx, tt.args.name, int64(tt.args.value))
			require.NoError(t, err)
			getVal, _ := storage.GetCounter(context.Background(), tt.args.name)
			require.Equal(t, getVal, tt.want)
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
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
				value: 0.1,
			},
			want: 0.1,
		},
	}

	ctx := context.Background()
	storage := NewMemStorage("")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.SetGauge(ctx, tt.args.name, float64(tt.args.value))
			require.NoError(t, err)
			getVal, _ := storage.GetGauge(context.Background(), tt.args.name)
			require.Equal(t, getVal, tt.want)
		})
	}
}

func TestMemStorage_GetAll(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name string
		args args
		want []Value
	}{
		{
			name: "God",
			args: args{
				name:  "test",
				value: int64(1),
			},
			want: []Value{
				{Name: "test", Type: "counter", Value: int64(1)},
			},
		},
	}

	ctx := context.Background()
	storage := NewMemStorage("")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.SetCounter(ctx, tt.args.name, int64(tt.args.value))
			require.NoError(t, err)
			getVal, _ := storage.GetAll(context.Background())
			require.Equal(t, getVal, tt.want)
		})
	}
}

func TestNewStorage(t *testing.T) {
	stor := NewMemStorage("test")

	assert.NotEmpty(t, stor)
}

func TestSaveToFile(t *testing.T) {
	stor := NewMemStorage("test")
	err := stor.SaveToFile()

	require.NoError(t, err)
}

func TestRestoreFromFile(t *testing.T) {
	stor := NewMemStorage("test")
	err := stor.RestoreFromFile()

	require.NoError(t, err)
}

func TestPing(t *testing.T) {
	ctx := context.Background()
	stor := NewMemStorage("test")
	err := stor.Ping(ctx)

	require.NoError(t, err)
}

func TestSetBatch(t *testing.T) {
	var metrics []metrics.Metric
	ctx := context.Background()
	stor := NewMemStorage("test")
	err := stor.SetBatch(ctx, metrics)

	require.NoError(t, err)
}
