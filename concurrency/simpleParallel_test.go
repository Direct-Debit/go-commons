package concurrency

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
