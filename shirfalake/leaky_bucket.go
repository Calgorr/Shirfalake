package shirfalake

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type LeakyBucket struct {
	rdb               *redis.Client
	prefix            string
	leakyBucket       *redis.Script
	timeout           time.Duration
	emission_interval int
}

func NewLeakyBucket(rdb *redis.Client, prefix string, timeout time.Duration, emission_interval int) *LeakyBucket {
	leakyBucket := &LeakyBucket{
		rdb:               rdb,
		prefix:            prefix,
		timeout:           timeout,
		emission_interval: emission_interval,
	}

	leakyBucket.leakyBucket = redis.NewScript(`
	local key = KEYS[1]
	local ts = tonumber(ARGV[1])
	local emission_interval = tonumber(ARGV[2])
	
	local last_ts = redis.call("GET",key)
	if not last_ts then
		last_ts = ts
	else
		last_ts = tonumber(last_ts)
	end
	
	local limited=1
	local remaining = last_ts+emission_interval-ts
	
	if ts==last_ts then
		limited=0
		remaining=0
		redis.call("SET", key, ts, "PX", math.ceil(emission_interval))
	
	end
	
	
	
	return {
		limited, -- allowed
		remaining , -- remaining
		last_ts, -- last_ts
		ts, -- ts
		emission_interval, -- emission_interval
	}
`)

	return leakyBucket
}

func (leakyBucket *LeakyBucket) Allow(key string) (bool, int) {
	ts := time.Now().UnixMilli()
	res, err := leakyBucket.leakyBucket.Run(leakyBucket.rdb, []string{key}, ts, leakyBucket.emission_interval).Result()
	if err != nil {
		fmt.Println(err)
		return false, 0
	}

	limited := res.([]interface{})[0].(int64)
	remaining := res.([]interface{})[1].(int64)

	return limited == 0, int(remaining)
}
