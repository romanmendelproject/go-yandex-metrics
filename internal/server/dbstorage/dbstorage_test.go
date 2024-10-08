package dbstorage

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/dbstorage/mocks"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"
	"github.com/romanmendelproject/go-yandex-metrics/utils"
	"github.com/stretchr/testify/require"
)

func TestNewPostgresStorage(t *testing.T) {
	// Test case 1: Successful connection
	connString := "postgres://user:password@localhost:5432/database"
	ctx := context.Background()
	storage := NewPostgresStorage(ctx, connString)
	if storage == nil {
		t.Errorf("Expected non-nil PostgresStorage instance, got nil")
	}
	if storage.db == nil {
		t.Errorf("Expected non-nil db instance, got nil")
	}
}

func TestGetCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	db.EXPECT().GetCounter(ctx, "test").Return(int64(1), nil)

	got, err := db.GetCounter(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, int64(1), got)
}

func TestGetCounterErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	errExp := errors.New("error")

	db.EXPECT().GetCounter(ctx, "test").Return(int64(0), errExp)

	got, err := db.GetCounter(ctx, "test")
	require.Error(t, err)
	require.EqualError(t, err, errExp.Error())
	require.Equal(t, int64(0), got)
}

func TestGetGauge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	db.EXPECT().GetGauge(ctx, "test").Return(float64(1), nil)

	got, err := db.GetGauge(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, float64(1), got)
}

func TestGetGaugeErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	errExp := errors.New("error")

	db.EXPECT().GetGauge(ctx, "test").Return(float64(0), errExp)

	got, err := db.GetGauge(ctx, "test")
	require.Error(t, err)
	require.EqualError(t, err, errExp.Error())
	require.Equal(t, float64(0), got)
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()
	values := []storage.Value{{Name: "Alloc", Type: "gauge", Value: float64(0)}}
	db.EXPECT().GetAll(ctx).Return(values, nil)

	got, err := db.GetAll(ctx)
	require.NoError(t, err)
	require.Equal(t, values, got)
}

func TestGetAllErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	errExp := errors.New("error")
	values := []storage.Value{{Name: "Alloc", Type: "gauge", Value: float64(0)}}

	db.EXPECT().GetAll(ctx).Return(values, errExp)

	got, err := db.GetAll(ctx)
	require.Error(t, err)
	require.EqualError(t, err, errExp.Error())
	require.Equal(t, values, got)
}

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	db.EXPECT().Ping(ctx).Return(nil)

	err := db.Ping(ctx)
	require.NoError(t, err)
}

func TestPingErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	errExp := errors.New("error")

	db.EXPECT().Ping(ctx).Return(errExp)

	err := db.Ping(ctx)
	require.Error(t, err)
	require.EqualError(t, err, errExp.Error())
}

func TestSetCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	db.EXPECT().SetCounter(ctx, "test", int64(1)).Return(nil)

	err := db.SetCounter(ctx, "test", int64(1))
	require.NoError(t, err)
}

func TestSetCounterErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	errExp := errors.New("error")

	db.EXPECT().SetCounter(ctx, "test", int64(1)).Return(errExp)

	err := db.SetCounter(ctx, "test", int64(1))
	require.Error(t, err)
	require.EqualError(t, err, errExp.Error())
}

func TestSetGauge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	db.EXPECT().SetGauge(ctx, "test", float64(0)).Return(nil)

	err := db.SetGauge(ctx, "test", float64(0))
	require.NoError(t, err)
}

func TestSetGaugeErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	errExp := errors.New("error")

	db.EXPECT().SetGauge(ctx, "test", float64(0)).Return(errExp)

	err := db.SetGauge(ctx, "test", float64(0))
	require.Error(t, err)
	require.EqualError(t, err, errExp.Error())
}

func TestSetBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	values := []metrics.Metric{
		{ID: "Alloc", MType: "gauge", Value: utils.GetFloatPtr(float64(0))}}

	db.EXPECT().SetBatch(ctx, values).Return(nil)

	err := db.SetBatch(ctx, values)
	require.NoError(t, err)
}

func TestSetBatchErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorage(ctrl)
	ctx := context.Background()

	values := []metrics.Metric{
		{ID: "Alloc", MType: "gauge", Value: utils.GetFloatPtr(float64(0))}}
	errExp := errors.New("error")

	db.EXPECT().SetBatch(ctx, values).Return(errExp)

	err := db.SetBatch(ctx, values)
	require.Error(t, err)
	require.EqualError(t, err, errExp.Error())
}
