package bucket_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vivangkumar/ratelimit/internal/bucket"
)

func TestBucket_Take(t *testing.T) {
	b := bucket.New(1)

	assert.True(t, b.Take())
	assert.EqualValues(t, b.Available(), 0)
	assert.False(t, b.Take())
}

func TestBucket_TakeN(t *testing.T) {
	b := bucket.New(10)

	assert.True(t, b.TakeN(10))
	assert.EqualValues(t, b.Available(), 0)
	assert.False(t, b.TakeN(10))
}

func TestBucket_Refill(t *testing.T) {
	b := bucket.New(10)
	assert.True(t, b.TakeN(10))

	b.Refill(1)
	assert.EqualValues(t, b.Available(), 1)

	// Fill to max.
	b.Refill(20)
	assert.EqualValues(t, b.Available(), 10)
}
