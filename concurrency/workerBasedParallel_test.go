package concurrency

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWorkers(t *testing.T) {
	input := make([]int, 100)
	for i := range input {
		input[i] = i
	}

	output := Workers(25, input, func(i int) (int, bool) {
		fmt.Printf("Executing job %d...\n", i)
		time.Sleep(time.Second)
		return 0, true
	})
	assert.Equal(t, 100, len(output))
	fmt.Println("-----------------------------------------------------")

	output2 := Workers(10, input, func(i int) (any, bool) {
		fmt.Printf("Executing job %d...\n", i)
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
	fmt.Println("-----------------------------------------------------")

	output3 := Workers(3, input, func(i int) (struct{}, bool) {
		fmt.Printf("Executing job %d...\n", i)
		time.Sleep(time.Millisecond * time.Duration(10*i))
		fmt.Printf("Job %d DONE!\n", i)
		return struct{}{}, false
	})
	assert.Equal(t, 0, len(output3))
	fmt.Println("-----------------------------------------------------")
}

func TestRateLimitedWorkers(t *testing.T) {
	input := make([]int, 100)
	for i := range input {
		input[i] = i
	}

	var jobsRunning int
	start := time.Now()

	rate := time.Millisecond * 100
	output := RateLimitedWorkers(rate, 5, input, func(i int) (int, bool) {
		jobsRunning += 1
		fmt.Printf("Executing job %d...\n", i)
		fmt.Printf("Jobs running: %d\n", jobsRunning)
		fmt.Printf("Milliseconds since start: %d\n", int(time.Since(start).Milliseconds()))
		time.Sleep(time.Second)
		fmt.Printf("Job %d DONE!\n", i)
		jobsRunning -= 1
		return 0, true
	})
	assert.Equal(t, 100, len(output))
	fmt.Println("-----------------------------------------------------")

	start = time.Now()
	output2 := RateLimitedWorkers(rate, 10, input, func(i int) (any, bool) {
		jobsRunning += 1
		fmt.Printf("Executing job %d...\n", i)
		fmt.Printf("Jobs running: %d\n", jobsRunning)
		fmt.Printf("Milliseconds since start: %d\n", int(time.Since(start).Milliseconds()))
		success := i%3 != 0
		if success {
			time.Sleep(time.Second)
			fmt.Println("Success")
		} else {
			fmt.Println("Failure")
		}
		jobsRunning -= 1
		return nil, success
	})
	assert.Equal(t, 66, len(output2))
	fmt.Println("-----------------------------------------------------")

	start = time.Now()
	output3 := RateLimitedWorkers(rate, 3, input, func(i int) (struct{}, bool) {
		jobsRunning += 1
		fmt.Printf("Executing job %d...\n", i)
		fmt.Printf("Jobs running: %d\n", jobsRunning)
		fmt.Printf("Milliseconds since start: %d\n", int(time.Since(start).Milliseconds()))
		time.Sleep(time.Millisecond * time.Duration(10*i))
		fmt.Printf("Job %d DONE!\n", i)
		jobsRunning -= 1
		return struct{}{}, false
	})
	assert.Equal(t, 0, len(output3))
	fmt.Println("-----------------------------------------------------")
}
