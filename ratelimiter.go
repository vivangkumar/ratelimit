package ratelimiter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vivangkumar/ratelimit/internal/bucket"
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
	// bucket is the underlying storage for the rate limiter.
	bucket *bucket.Bucket

	// refillDuration is the duration after which tokens are refilled.
	//
	// The duration is calculated based on the limit specified
	// at creation time.
	refillDuration time.Duration
	now            NowFunc

	m sync.Mutex
	// lastRefillUnixMs keeps track of the time when the last
	// refresh of tokens took place.
	//
	// It is kept track in milliseconds.
	lastRefillUnixMs int64
}

// New constructs a rate limiter that accepts the max tokens (size) that
// the limiter holds, along with the limit per duration.
//
// For example: if the bucket is configured with a max of 100 tokens
// and limit is set to 10 over a duration of 1m, this implies that
// the bucket will be refilled with one token every ((1 * 60s) / limit)s.
// This refills the bucket with 1 token every 6s, while giving us a
// max "burst" of 100 tokens.
func New(
	max uint64,
	limit uint64,
	per time.Duration,
	opts ...Opt,
) (*RateLimiter, error) {
	if limit == 0 {
		return nil, fmt.Errorf("limit must be positive")
	}

	r := &RateLimiter{
		bucket:         bucket.New(max),
		refillDuration: time.Duration(float64(per) / float64(limit)),
		now:            time.Now,
	}

	for _, opt := range opts {
		opt(r)
	}
	r.lastRefillUnixMs = r.now().Unix()

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
	lastRefill := r.lastRefillUnixMs
	r.m.Unlock()

	now := r.now()
	tokens := (now.UnixMilli() - lastRefill) / r.refillDuration.Milliseconds()
	if tokens > 0 {
		r.m.Lock()
		r.lastRefillUnixMs = now.UnixMilli()
		r.m.Unlock()

		r.bucket.Refill(uint64(tokens))
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
	return r.bucket.TakeN(n)
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
	if n > r.bucket.Size() {
		return fmt.Errorf("tokens requested exceeds max tokens")
	}

	// Check refillEvery duration to see if a new token is available.
	t := time.NewTicker(r.refillDuration)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			if !r.AddN(n) {
				continue
			}

			return nil
		}
	}
}
