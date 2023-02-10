package concurrency

import (
	"github.com/Direct-Debit/go-commons/stdext"
	"sync"
)

// SimpleParallel executes the transform function in a go routine for every input.
// If the input slice is of size N, SimpleParallel spawns at least N goroutines.
// The bool return value for transform indicates whether the transform executed successfully.
// Only successful transform results will be included in the result slice.
// No ordering is guaranteed in the result slice.
func SimpleParallel[I any, O any](input []I, transform func(I) (O, bool)) []O {
	results := make(chan O)
	go func() {
		wg := sync.WaitGroup{}

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
		}

		wg.Wait()
		close(results)
	}()

	return stdext.ChanToSlice(results)
}
