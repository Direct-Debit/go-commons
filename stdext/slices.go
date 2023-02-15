package stdext

import "fmt"

func SafeIdx[T any](idx int, arr []T) T {
	if 0 <= idx && idx < len(arr) {
		return arr[idx]
	}
	var nothing T
	return nothing
}

func SafeSlice[T any](slice []T, start, end int) []T {
	end = Max(0, Min(len(slice), end))
	start = Min(Max(0, start), end)
	return slice[start:end]
}

func FindInSlice[T comparable](slice []T, val T) (int, bool) {
	for k, item := range slice {
		if item == val {
			return k, true
		}
	}
	return -1, false
}

func Filter[T any](slice []T, keep func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, t := range slice {
		if keep(t) {
			result = append(result, t)
		}
	}
	return result
}

func Map[T any, U any](slice []T, conv func(T) U) []U {
	result := make([]U, len(slice))
	for i, t := range slice {
		result[i] = conv(t)
	}
	return result
}

// MapErr creates a new slice with one element for every element in the given slice.
// The elements in the new slice will be the return values of the conversion function conv.
// If conv returns an error, MapErr will return an error wrapping the original error
// and a slice of the same length as the given slice containing the elements that were successfully converted before the error occurred.
// MapErr does not generate any of its own errors.
func MapErr[T any, U any](slice []T, conv func(T) (U, error)) ([]U, error) {
	var err error
	result := make([]U, len(slice))
	for i, t := range slice {
		result[i], err = conv(t)
		if err != nil {
			return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
		}
	}
	return result, nil
}

func SafeIdxStr(idx int, arr []string) string {
	return SafeIdx(idx, arr)
}

func FindInStrSlice(slice []string, val string) (int, bool) {
	return FindInSlice(slice, val)
}

// ChanToSlice listens to the given channel and appends all the messages to a slice.
// ChanToSlice blocks until the channel is closed.
// If the channel is never closed, ChanToSlice will block forever.
func ChanToSlice[T any](c chan T) []T {
	result := make([]T, 0)
	for t := range c {
		result = append(result, t)
	}
	return result
}

func Flatten[T any](slices [][]T) []T {
	size := 0
	for _, s := range slices {
		size += len(s)
	}
	result := make([]T, 0, size)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// ChunkifyByCount returns a slice containing n slices of roughly equal size.
// The slices will contain the copied elements from the given slice divided up in FIFO fashion.
func ChunkifyByCount[T any](slice []T, n int) [][]T {
	chunks := Max(n, 1)
	length := len(slice)

	chunkSize := length / chunks
	specialChunks := length - chunkSize*chunks

	result := make([][]T, 0, chunks)
	for i, start := 0, 0; start < length; i++ {
		end := Min(start+chunkSize, length)
		if i < specialChunks {
			end += 1
		}

		result = append(result, slice[start:end])
		start = end
	}
	return result
}

// ChunkifyBySize returns a slice containing slices, each with at most chunkSize elements.
// The slices will contain the copied elements from the given slice divided up in FIFO fashion.
func ChunkifyBySize[T any](slice []T, chunkSize int) [][]T {
	length := len(slice)

	chunks := length / chunkSize
	for chunkSize*chunks > length {
		chunks += 1
	}

	result := make([][]T, 0, chunks)
	for start := 0; start < length; {
		end := Min(start+chunkSize, length)
		result = append(result, slice[start:end])
		start = end
	}
	return result
}
