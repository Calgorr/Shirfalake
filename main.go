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

	burst := 5
	rps := 10
	testDuration := time.Second * 15
	rl := shirfalake.NewGCRA(rdb, "test", time.Second, rps, burst)
	rateLimiter := shirfalake.NewRateLimiter(rl)

	acceptedCount := 0
	reqCount := 0
	reqTime := time.Duration(0)

	start := time.Now()
	t := time.NewTicker(time.Millisecond * 10)
	for {
		<-t.C

		callStart := time.Now()
		can, remaining := rateLimiter.Rl.Allow("gcra-test")

		reqTime += time.Since(callStart)
		reqCount++
		if can {
			acceptedCount++
			fmt.Println("#", acceptedCount, ": accepted > ", time.Since(start))
		} else {
			fmt.Println("#", acceptedCount, ": rejected > ", time.Since(start), " remaining: ", remaining)
		}

		elapsedTime := time.Since(start)

		if elapsedTime > testDuration {
			break
		}

		if elapsedTime > time.Second*7 && elapsedTime < time.Second*8 {
			// three seconds of rest, it should open up capacity for 3 more requests
			<-time.After(time.Second * 3)
		}
	}

	fmt.Println("Redis average response time ", time.Duration(reqTime.Nanoseconds()/int64(reqCount)))
}
