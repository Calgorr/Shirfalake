package shirfalake

import (
	"time"

	"github.com/go-redis/redis"
)

type GCRA struct {
	rdb        *redis.Client
	prefix     string
	gcra       *redis.Script
	timeout    time.Duration
	rps, burst int
}

func NewGCRA(rdb *redis.Client, prefix string, timeout time.Duration, rps, burst int) *GCRA {
	gcra := &GCRA{
		rdb:     rdb,
		prefix:  prefix,
		timeout: timeout,
		rps:     rps,
		burst:   burst,
	}

	gcra.gcra = redis.NewScript(`
	local rate_limit_key = KEYS[1]
	local now            = ARGV[1]
	local burst          = ARGV[2]
	local rate           = ARGV[3]
	local period         = ARGV[4]
	
	local emission_interval = period / rate
	local burst_offset      = emission_interval * burst
	
	local tat = redis.call("GET", rate_limit_key)
	
	if not tat then
	  tat = now
	else
	  tat = tonumber(tat)
	end
	tat = math.max(tat, now)
	
	local new_tat = tat + emission_interval
	local allow_at = new_tat - burst_offset
	local diff = now - allow_at
	
	local limited
	local retry_in
	local reset_in
	
	local remaining = math.floor(diff / emission_interval) -- calculate how many tokens there actually are
	
	if remaining < 0 then
		limited = 1
		remaining = math.floor((now - (tat - burst_offset)) / emission_interval)
		reset_in = math.ceil(tat - now)
		retry_in = math.ceil(diff * -1)
	  elseif remaining == 0 and increment <= 0 then
		limited = 1
		remaining = 0
		reset_in = math.ceil(tat - now)
		retry_in = 0
	  else
		limited = 0
		reset_in = math.ceil(new_tat - now)
		retry_in = 0
		if increment > 0 then
		  redis.call("SET", rate_limit_key, new_tat, "PX", reset_in)
		end
	  end
	  
	  return {limited, remaining, retry_in, reset_in, tostring(diff), tostring(emission_interval)}
	`)

	return gcra
}

func (g *GCRA) Allow(key string) (bool, int) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	result, err := g.gcra.Run(g.rdb, []string{g.prefix + key}, now, g.burst, g.rps, g.timeout.Milliseconds()).Result()
	if err != nil {
		return false, 0
	}

	limited := result.([]interface{})[0].(int64) == 1
	retryIn := int(result.([]interface{})[2].(int64))

	return limited, retryIn
}
