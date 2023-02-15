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
