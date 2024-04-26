package shirfalake

type RateLimit interface {
	// Check if the request is allowed
	Allow(string) (bool, int)
}

type RateLimiter struct {
	rl *RateLimit
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter(rl RateLimit, rps, burst int) *RateLimiter {
	return &RateLimiter{
		rl: &rl,
	}
}
