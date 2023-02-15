package concurrency

import (
	"github.com/Direct-Debit/go-commons/stdext"
	"sync"
)

// Workers executes the transform function in a go routine for every input.
// If the input slice is of size N, Workers spawns N goroutines, but no more than workerCount at a time.
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
