package shirfalake

type RateLimit interface {
	// Check if the request is allowed
	Allow(string) (bool, int)
}

type RateLimiter struct {
	Rl RateLimit
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter(Rl RateLimit) *RateLimiter {
	return &RateLimiter{
		Rl: Rl,
	}
}
