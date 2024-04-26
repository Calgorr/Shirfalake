package shirfalake

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type GCRA struct {
	rdb               *redis.Client
	prefix            string
	gcra              *redis.Script
	timeout           time.Duration
	emission_interval float64
	burst             int
}

func NewGCRA(rdb *redis.Client, prefix string, timeout time.Duration, emission_interval float64, burst int) *GCRA {
	gcra := &GCRA{
		rdb:               rdb,
		prefix:            prefix,
		timeout:           timeout,
		emission_interval: emission_interval,
		burst:             burst,
	}

	gcra.gcra = redis.NewScript(`
	redis.replicate_commands()
	local rate_limit_key = KEYS[1]
	local burst = ARGV[1]
	local emission_interval = tonumber(ARGV[2])
	
	local jan_1_2017 = 1483228800
	local now = redis.call("TIME")
	now = (now[1] - jan_1_2017) + (now[2] / 1000000)
	
	local tat = redis.call("GET", rate_limit_key)
	if not tat then
		tat = now
	else
		tat = tonumber(tat)
	end
	
	tat = math.max(tat, now)
	
	local burst_offset = emission_interval * burst
	local allow_at = tat - burst_offset
	local diff = now - allow_at
	if diff < 0 then
		return {
			0,
			0,
			-diff,
		}
	end
	
	local new_tat = tat + emission_interval
	local key_expiration = new_tat - now
	
	redis.call("SET", rate_limit_key, new_tat, "EX", math.ceil(key_expiration))
	
	return {
		1,
		tostring(math.ceil(diff / emission_interval)), -- remaining
		0,
	}
	`)

	return gcra
}

func (g *GCRA) Allow(key string) (bool, int) {
	values, err := g.gcra.Run(g.rdb, []string{key}, g.burst, g.emission_interval).Result()
	if err != nil {
		fmt.Println(err)
		return false, 0
	}

	return values.([]interface{})[0].(int64) == 1, 0
}
