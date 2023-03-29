package ratelimiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucket_take(t *testing.T) {
	b := newBucket(1)
	assert.True(t, b.take())

	assert.EqualValues(t, b.size, 1)
	assert.EqualValues(t, b.available, 0)

	assert.False(t, b.take())
}

func TestBucket_takeN(t *testing.T) {
	b := newBucket(10)
	assert.True(t, b.takeN(10))

	assert.EqualValues(t, b.size, 10)
	assert.EqualValues(t, b.available, 0)

	assert.False(t, b.takeN(10))
}

func TestBucket_refill(t *testing.T) {
	b := newBucket(10)
	assert.True(t, b.takeN(10))

	b.refill(1)
	assert.EqualValues(t, b.available, 1)

	// Fill to max.
	b.refill(20)
	assert.EqualValues(t, b.available, 10)
}
