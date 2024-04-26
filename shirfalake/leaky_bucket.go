package shirfalake

import (
	"time"

	"github.com/go-redis/redis"
)

type LeakyBucket struct {
	rdb     *redis.Client
	prefix  string
	gcra    *redis.Script
	timeout time.Duration
}

func NewLeakyBucket(rdb *redis.Client, prefix string, timeout time.Duration) *LeakyBucket {
	gcra := &LeakyBucket{
		rdb:     rdb,
		prefix:  prefix,
		timeout: timeout,
	}

	gcra.gcra = redis.NewScript(`local time_stamp  = tonumber(ARGV[1])
	local rps = tonumber(ARGV[2])
	local key = KEYS[1]
	
	local min = time_stamp -1
	redis.call('ZREMRANGEBYSCORE', key, '-inf', min)
	
	local last = redis.call('ZRANGE', key, -1, -1)
	
	local tat = time_stamp
	
	if type(last) == 'table' and #last > 0 then
	  for key,value in pairs(last) do
		tat = tonumber(value) + 1/rps
	
		break
	  end
	end
	
	if ts > tat then
		next = ts  
	end
	
	redis.call('ZADD', key, tat, tat)
	
	
	local limited=0
	
	if tat > ts then
	  limited = 1
	end
	
	return {limited, tostring(tat-ts)}  -- return the remaining time`)

	return gcra
}

func (gcra *LeakyBucket) Allow(key string, rps int) (bool, int) {
	ts := time.Now().UnixNano()
	res, err := gcra.gcra.Run(gcra.rdb, []string{gcra.prefix + key}, ts, rps).Result()
	if err != nil {
		return false, 0
	}

	limited := res.([]interface{})[0].(int64)
	remaining := res.([]interface{})[1].(int64)

	return limited == 0, int(remaining)
}
