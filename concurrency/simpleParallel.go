package concurrency

import (
	"sync"

	"github.com/Direct-Debit/go-commons/stdext"
)

// SimpleParallel executes the transform function in a go routine for every input.
// If the input slice is of size N, SimpleParallel spawns at least N goroutines.
//
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

// SimpleParallelMap executes the transform function in a go routine for every key-value pair in the input map.
// If the input map has N key-value paids, SimpleParallel spawns at least N goroutines.
//
// The bool return value for transform indicates whether the transform executed successfully.
// Only successful transform results will be included in the result map.
func SimpleParallelMap[K comparable, V any, O any](input map[K]V, transform func(V) (O, bool)) map[K]O {
	type transformResult struct {
		key    K
		output O
	}

	resultChan := make(chan transformResult)
	go func() {
		wg := sync.WaitGroup{}

		for k, v := range input {
			wg.Add(1)
			go func(k K, v V) {
				defer wg.Done()

				o, ok := transform(v)
				if ok {
					resultChan <- transformResult{k, o}
				}
			}(k, v)
		}

		wg.Wait()
		close(resultChan)
	}()

	resultMap := make(map[K]O)
	for result := range resultChan {
		resultMap[result.key] = result.output
	}
	return resultMap
}
