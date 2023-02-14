package concurrency

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRateLimitedConcurrent(t *testing.T) {
	input := make([]int, 100)
	for i := range input {
		input[i] = i
	}

	rate := CPSToDuration(25)
	assert.Equal(t, time.Millisecond*40, rate)

	start := time.Now()
	output := RateLimitedConcurrent(rate, input, func(i int) (int, bool) {
		time.Sleep(time.Second)
		fmt.Printf("Thread %d finsished at second %.3f\n", i, time.Since(start).Seconds())
		return 0, true
	})
	assert.Equal(t, 100, len(output))
	assert.Less(t, time.Since(start).Seconds(), 5.0)

	output2 := RateLimitedConcurrent(0, input, func(i int) (any, bool) {
		success := i%3 != 0
		if success {
			time.Sleep(time.Second)
			fmt.Println("Success")
		} else {
			fmt.Println("Failure")
		}
		return nil, success
	})
	assert.Equal(t, 66, len(output2))
}

func TestRateLimitedSerial(t *testing.T) {
	input := make([]int, 5)
	for i := range input {
		input[i] = i
	}

	start := time.Now()
	output := RateLimitedSerial(time.Second+time.Millisecond*50, input, func(i int) (int, bool) {
		time.Sleep(time.Second)
		fmt.Printf("Thread %d finsished at second %.3f\n", i, time.Since(start).Seconds())
		return 0, true
	})
	assert.Equal(t, 5, len(output))
	assert.Greater(t, time.Since(start).Seconds(), 5.0)

	output2 := RateLimitedSerial(0, input, func(i int) (any, bool) {
		success := i%3 != 0
		if success {
			time.Sleep(time.Second)
			fmt.Println("Success")
		} else {
			fmt.Println("Failure")
		}
		return nil, success
	})
	assert.Equal(t, 3, len(output2))
}
