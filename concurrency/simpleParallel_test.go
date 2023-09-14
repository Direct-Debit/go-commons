package concurrency

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleParallel(t *testing.T) {
	input := make([]int, 100)
	for i := range input {
		input[i] = i
	}

	output := SimpleParallel(input, func(int) (int, bool) {
		fmt.Println("Execute...")
		time.Sleep(time.Second)
		return 0, true
	})
	assert.Equal(t, 100, len(output))

	output2 := SimpleParallel(input, func(i int) (any, bool) {
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

func TestSimpleParallelMap(t *testing.T) {
	input := make(map[string]int)
	for i := 0; i < 100; i++ {
		input[fmt.Sprint(i)] = i
	}

	output := SimpleParallelMap(input, func(v int) (int, bool) {
		fmt.Println("Execute...")
		time.Sleep(time.Second)
		return v * 2, true
	})
	assert.Equal(t, 100, len(output))

	output2 := SimpleParallelMap(input, func(i int) (any, bool) {
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
