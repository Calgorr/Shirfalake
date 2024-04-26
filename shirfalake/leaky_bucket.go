package shirfalake

import (
	"time"

	"github.com/go-redis/redis"
)

type LeakyBucket struct {
	rdb         *redis.Client
	prefix      string
	leakyBucket *redis.Script
	timeout     time.Duration
	rps         int
}

func NewLeakyBucket(rdb *redis.Client, prefix string, timeout time.Duration, rps int) *LeakyBucket {
	leakyBucket := &LeakyBucket{
		rdb:     rdb,
		prefix:  prefix,
		timeout: timeout,
		rps:     rps,
	}

	leakyBucket.leakyBucket = redis.NewScript(`local time_stamp  = tonumber(ARGV[1])
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

	return leakyBucket
}

func (leakyBucket *LeakyBucket) Allow(key string) (bool, int) {
	ts := time.Now().UnixNano()
	res, err := leakyBucket.leakyBucket.Run(leakyBucket.rdb, []string{leakyBucket.prefix + key}, ts, leakyBucket.rps).Result()
	if err != nil {
		return false, 0
	}

	limited := res.([]interface{})[0].(int64)
	remaining := res.([]interface{})[1].(int64)

	return limited == 0, int(remaining)
}
