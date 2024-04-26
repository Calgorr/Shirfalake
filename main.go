package main

import (
	"fmt"
	"time"

	"github.com/Calgorr/Shirfalake/shirfalake"
	"github.com/go-redis/redis"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// burst := 5 // burst 5 for gcra
	emission_interval := 1000
	testDuration := time.Second * 15
	rl := shirfalake.NewLeakyBucket(rdb, "test", time.Second, emission_interval) // for leaky bucket you can change to gcra
	rateLimiter := shirfalake.NewRateLimiter(rl)

	acceptedCount := 0
	reqCount := 0
	reqTime := time.Duration(0)

	start := time.Now()
	t := time.NewTicker(time.Millisecond * 10)
	for {
		<-t.C

		callStart := time.Now()
		can, _ := rateLimiter.Rl.Allow("leaky-test1")

		reqTime += time.Since(callStart)
		reqCount++
		if can {
			acceptedCount++
			fmt.Println("#", acceptedCount, ": accepted > ", time.Since(start))
		}

		elapsedTime := time.Since(start)

		if elapsedTime > testDuration {
			break
		}

		// if elapsedTime > time.Second*7 && elapsedTime < time.Second*8 { // also for gcra uncomment this block if you want to test gcra
		// 	// three seconds of rest, it should open up capacity for 3 more requests
		// 	<-time.After(time.Second * 3)
		// }
	}

	fmt.Println("Redis average response time ", time.Duration(reqTime.Nanoseconds()/int64(reqCount)))
}
