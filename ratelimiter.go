package ratelimiter

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// NowFunc helps with mocking time.
type NowFunc = func() time.Time

// RateLimiter represents a token bucket based rate limiter.
//
// It is built on top of a bucket that accepts a max size.
// The bucket is refilled at an interval determined by the limiter.
//
// Most callers should use either Wait or WaitN to wait for tokens
// to be available.
type RateLimiter struct {
	m sync.Mutex
	// bucket is the underlying storage structure
	// for the rate limiter.
	bucket *bucket

	// lastRefillAt keeps track of the time when the last
	// refresh of tokens took place.
	//
	// This is tracked as a unix nanoseconds timestamp for
	// maximum granularity.
	lastRefillAt int64

	// refillEvery is the duration after which tokens are refilled.
	//
	// The duration is calculated based on the limit specified
	// at creation time.
	refillEvery time.Duration
	now         NowFunc
}

// New constructs a rate limiter that accepts the max tokens that
// the limiter holds, along with the limit per duration.
//
// For example: if the bucket is configured with a max of 100 tokens
// and limit is set to 10 over a duration of 1m, this implies that
// the bucket will be refilled with one token every ((1 * 60s) / limit)s.
// This refills the bucket with 1 token every 6s, while giving us a
// max "burst" of 100 tokens.
func New(
	maxTokens uint64,
	limit uint64,
	per time.Duration,
	opts ...Opt,
) (*RateLimiter, error) {
	if limit == 0 {
		return nil, fmt.Errorf("limit must be positive")
	}

	r := &RateLimiter{
		bucket:      newBucket(maxTokens),
		refillEvery: time.Duration(float64(per) / float64(limit)),
		now:         time.Now,
	}

	for _, opt := range opts {
		opt(r)
	}
	r.lastRefillAt = r.now().UnixNano()

	return r, nil
}

// refill is responsible for refilling the bucket with
// one token every refill period.
//
// This method is called when attempting to Add as we
// might have to refresh our token count before allowing
// the token to be taken.
func (r *RateLimiter) refill() {
	r.m.Lock()
	defer r.m.Unlock()

	now := r.now().UnixNano()
	diff := now - r.lastRefillAt
	if diff == 0 {
		return
	}

	// Round up the value so that we can be as accurate as possible.
	refills := math.Round(float64(diff) / float64(r.refillEvery.Nanoseconds()))
	if refills > 0 {
		r.bucket.refill(uint64(math.Round(refills)))
		r.lastRefillAt = now
	}
}

// Add attempts to take a single token from the bucket.
//
// If there are tokens available, it returns true.
// Otherwise, the method returns false, indicating that we
// have reached the rate limit.
//
// Callers should retry the request to take a token from the
// rate limiter the next time.
func (r *RateLimiter) Add() bool {
	return r.AddN(1)
}

// AddN attempts to acquire n tokens from the bucket.
// Its behaviour is details in bucket.takeN.
func (r *RateLimiter) AddN(n uint64) bool {
	r.refill()

	r.m.Lock()
	ok := r.bucket.takeN(n)
	r.m.Unlock()

	return ok
}

// Wait blocks until a token is available.
//
// It returns an error if the context is cancelled,
// or if the wait time for the context is exceeded.
//
// This method consumes a token if successful.
func (r *RateLimiter) Wait(ctx context.Context) error {
	return r.WaitN(ctx, 1)
}

// WaitN blocks until n tokens are available.
//
// It returns an error if the context is cancelled,
// if the wait time for the context is exceeded, or
// if the number of tokens request exceeds the maximum
// available tokens.
//
// This method also consumes n tokens, if successful.
func (r *RateLimiter) WaitN(ctx context.Context, n uint64) error {
	if n > r.bucket.size {
		return fmt.Errorf("tokens requested exceeds max tokens")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if r.AddN(n) {
				return nil
			}
		}
	}
}
