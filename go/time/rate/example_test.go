package rate_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/searKing/golang/go/time/rate"
)

func ExampleNewFullBurstLimiter() {
	const (
		burst = 3
	)
	limiter := rate.NewFullBurstLimiter(burst)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// expect dropped, as limiter is inited with full tokens(3)
	limiter.PutToken()

	for i := 0; ; i++ {
		//fmt.Printf("%03d %s\n", i, time.Now().Format(time.RFC3339))
		fmt.Printf("Wait %03d\n", i)
		err := limiter.Wait(ctx)
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}

		fmt.Printf("Got %03d\n", i)
		if i == 0 {
			// refill one token
			limiter.PutToken()
		}
	}
	// Output:
	// Wait 000
	// Got 000
	// Wait 001
	// Got 001
	// Wait 002
	// Got 002
	// Wait 003
	// Got 003
	// Wait 004
	// err: context deadline exceeded

}

func ExampleNewEmptyBurstLimiter() {
	const (
		burst       = 3
		concurrency = 2
	)
	limiter := rate.NewEmptyBurstLimiter(burst)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	fmt.Printf("tokens: %d\n", limiter.Tokens())

	// expect not allowed, as limiter is inited with empty tokens(0)
	if limiter.Allow() {
		fmt.Printf("allow passed\n")
	} else {
		fmt.Printf("allow refused\n")
	}
	// fill one token
	limiter.PutToken()
	fmt.Printf("tokens: %d\n", limiter.Tokens())

	// expect allowed, as limiter is filled with one token(1)
	if limiter.Allow() {
		fmt.Printf("allow passed\n")
	} else {
		fmt.Printf("allow refused\n")
	}
	fmt.Printf("tokens: %d\n", limiter.Tokens())

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			//fmt.Printf("%03d %s\n", i, time.Now().Format(time.RFC3339))
			fmt.Printf("Wait %03d\n", i)
			err := limiter.Wait(ctx)
			if err != nil {
				fmt.Printf("err: %s\n", err.Error())
				return
			}

			fmt.Printf("Got %03d\n", i)
			if i == 0 {
				// refill one token
				limiter.PutToken()
			}
		}()
	}

	for i := 0; i < concurrency; i++ {
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("PutToken #%d: before tokens: %d\n", i, limiter.Tokens())
		// fill one token
		limiter.PutToken()
		fmt.Printf("PutToken #%d: after tokens: %d\n", i, limiter.Tokens())
	}
	wg.Wait()
	fmt.Printf("tokens: %d\n", limiter.Tokens())

	// expect allowed, as limiter is filled with one token(1)
	if limiter.Allow() {
		fmt.Printf("allow passed\n")
	} else {
		fmt.Printf("allow refused\n")
	}
	fmt.Printf("tokens: %d\n", limiter.Tokens())

	// expect not allowed, as limiter is inited with empty tokens(0)
	if limiter.Allow() {
		fmt.Printf("allow passed\n")
	} else {
		fmt.Printf("allow refused\n")
	}
	// Output:
	// tokens: 0
	// allow refused
	// tokens: 1
	// allow passed
	// tokens: 0
	// Wait 001
	// Wait 000
	// PutToken #0: before tokens: 0
	// PutToken #0: after tokens: 1
	// Got 000
	// Got 001
	// PutToken #1: before tokens: 0
	// PutToken #1: after tokens: 1
	// tokens: 1
	// allow passed
	// tokens: 0
	// allow refused
}
