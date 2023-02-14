package concurrency

import (
	"github.com/Direct-Debit/go-commons/stdext"
	"sync"
	"time"
)

func CPSToDuration(callsPerSecond int) time.Duration {
	return time.Duration(int64(time.Second) / int64(callsPerSecond))
}

// RateLimitedConcurrent will concurrently call the transform function on the given input once per tick at the given rate.
// If rate is less than or equal to 0, it will be set to one nanosecond.
// If the input slice is of size N, RateLimitedConcurrent spawns at least N goroutines.
// The bool return value for transform indicates whether the transform executed successfully.
// Only successful transform results will be included in the result slice.
// No ordering is guaranteed in the result slice.
func RateLimitedConcurrent[I any, O any](rate time.Duration, input []I, transform func(I) (O, bool)) []O {
	if rate <= 0 {
		rate = time.Nanosecond
	}
	results := make(chan O)

	go func() {
		wg := sync.WaitGroup{}
		ticker := time.NewTicker(rate)

		for _, i := range input {
			copiedI := i

			wg.Add(1)
			go func() {
				defer wg.Done()

				o, ok := transform(copiedI)
				if ok {
					results <- o
				}
			}()
			<-ticker.C
		}

		wg.Wait()
		ticker.Stop()
		close(results)
	}()

	return stdext.ChanToSlice(results)
}

// RateLimitedSerial will concurrently call the transform function on the given input at most once per tick at the given rate.
// If rate is less than or equal to 0, it will be set to one nanosecond.
// The calls to the transform functions aren't concurrent,
// the next call to transform will only be called once the previous call has finished.
// The bool return value for transform indicates whether the transform executed successfully.
// Only successful transform results will be included in the result slice.
// The order of the output array will correspond to the order of the input array.
func RateLimitedSerial[I any, O any](rate time.Duration, input []I, transform func(I) (O, bool)) []O {
	if rate <= 0 {
		rate = time.Nanosecond
	}
	result := make([]O, 0)
	ticker := time.NewTicker(rate)

	for _, i := range input {
		o, ok := transform(i)
		if ok {
			result = append(result, o)
		}

		<-ticker.C
	}
	ticker.Stop()
	return result
}
