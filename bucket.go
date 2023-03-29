package ratelimiter

// bucket represents a token bucket.
//
// It is not safe for concurrent use.
type bucket struct {
	// size is the max tokens the bucket can hold.
	size uint64

	// available is the current number of available
	// tokens in the bucket.
	available uint64
}

// newBucket returns a token bucket configured with
// size number of max tokens.
func newBucket(size uint64) *bucket {
	return &bucket{
		size:      size,
		available: size,
	}
}

// take uses up a single token from the bucket.
//
// It returns true if a token could be acquired.
// Otherwise, it returns false.
func (b *bucket) take() bool {
	return b.takeN(1)
}

// takeN acquires n tokens from the bucket, if available.
//
// If n tokens are not available, no tokens are removed
// from the bucket.
func (b *bucket) takeN(n uint64) bool {
	if b.available >= n {
		b.available -= n
		return true
	}

	return false
}

// refill refills the bucket with n tokens.
//
// If the quantity to refill exceeds the size of the bucket,
// the bucket is refilled upto its size.
func (b *bucket) refill(n uint64) {
	t := b.available + n
	if t > b.size {
		t = b.size
	}
	b.available = t
}
