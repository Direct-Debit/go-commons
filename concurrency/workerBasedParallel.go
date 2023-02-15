package concurrency

import (
	"github.com/Direct-Debit/go-commons/stdext"
	"sync"
	"time"
)

// Workers executes the transform function in a go routine for every input.
// If the input slice is of size N, Workers spawns N goroutines, but no more than workerCount at a time.
//
// The bool return value for transform indicates whether the transform executed successfully.
// Only successful transform results will be included in the result slice.
// No ordering is guaranteed in the result slice.
func Workers[I any, O any](workerCount int, input []I, transform func(I) (O, bool)) []O {
	if workerCount <= 0 {
		workerCount = 1
	}

	results := make(chan O, len(input))
	wg := sync.WaitGroup{}
	workerQueue := make(chan struct{}, workerCount)

	for _, i := range input {
		workerQueue <- struct{}{}
		copiedI := i

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-workerQueue
			}()

			o, ok := transform(copiedI)
			if ok {
				results <- o
			}
		}()
	}

	wg.Wait()
	close(results)

	return stdext.ChanToSlice(results)
}

// RateLimitedWorkers will concurrently call the transform function on the given input.
// If rate is less than or equal to 0, it will be set to one nanosecond.
// If the input slice is of size N, RateLimitedWorkers spawns N goroutines, but no more than workerCount at a time.
// The go routines will  also only be spawned once per tick at the given rate.
//
// The bool return value for transform indicates whether the transform executed successfully.
// Only successful transform results will be included in the result slice.
// No ordering is guaranteed in the result slice.
func RateLimitedWorkers[I any, O any](rate time.Duration, workerCount int, input []I, transform func(I) (O, bool)) []O {
	if rate <= 0 {
		rate = time.Nanosecond
	}
	if workerCount <= 0 {
		workerCount = 1
	}

	results := make(chan O, len(input))
	wg := sync.WaitGroup{}
	workerQueue := make(chan struct{}, workerCount)
	ticker := time.NewTicker(rate)

	for _, i := range input {
		workerQueue <- struct{}{}
		copiedI := i

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-workerQueue
			}()

			o, ok := transform(copiedI)
			if ok {
				results <- o
			}
		}()
		<-ticker.C
	}

	wg.Wait()
	close(results)

	return stdext.ChanToSlice(results)
}
