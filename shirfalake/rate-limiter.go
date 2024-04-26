package shirfalake

type RateLimit interface {
	// Check if the request is allowed
	Allow(string, int) (bool, int)
}

type RateLimiter struct {
	rl         *RateLimit
	burst, rps int
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter(rl RateLimit, rps, burst int) *RateLimiter {
	return &RateLimiter{
		rl:    &rl,
		burst: burst,
		rps:   rps,
	}
}
