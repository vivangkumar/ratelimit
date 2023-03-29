package ratelimiter

type Opt func(r *RateLimiter)

// WithNowFunc sets the NowFunc to return the current
// time.
//
// By default, it is set to time.Now.
func WithNowFunc(n NowFunc) Opt {
	return func(r *RateLimiter) {
		r.now = n
	}
}
