package bucket

import (
	"sync"
)

// Bucket represents a token bucket.
//
// It is safe for concurrent use.
type Bucket struct {
	m sync.Mutex
	// size is the max tokens the bucket can hold.
	size uint64

	// available is the current number of available
	// tokens in the bucket.
	available uint64
}

// New returns a token bucket configured with
// size number of max tokens.
func New(size uint64) *Bucket {
	return &Bucket{
		m:         sync.Mutex{},
		size:      size,
		available: size,
	}
}

// Take uses up a single token from the bucket.
//
// It returns true if a token could be acquired.
// Otherwise, it returns false.
func (b *Bucket) Take() bool {
	return b.TakeN(1)
}

// TakeN acquires n tokens from the bucket, if available.
//
// If n tokens are not available, no tokens are removed
// from the bucket.
func (b *Bucket) TakeN(n uint64) bool {
	b.m.Lock()
	defer b.m.Unlock()

	if b.available >= n {
		b.available -= n
		return true
	}

	return false
}

// Refill refills the bucket with n tokens.
//
// If the quantity to refill exceeds the size of the bucket,
// the bucket is refilled upto its size.
func (b *Bucket) Refill(n uint64) {
	b.m.Lock()
	defer b.m.Unlock()

	t := b.available + n
	if t > b.size {
		t = b.size
	}
	b.available = t
}

// Available returns the tokens currently available in the bucket.
func (b *Bucket) Available() uint64 {
	b.m.Lock()
	defer b.m.Unlock()

	return b.available
}

// Size returns the max tokens a bucket can have.
func (b *Bucket) Size() uint64 {
	b.m.Lock()
	defer b.m.Unlock()

	return b.size
}
