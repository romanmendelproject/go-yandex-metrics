package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbInit(t *testing.T) {
	ctx := context.Background()

	storage := dbInit(ctx)

	assert.NotNil(t, storage)

}
