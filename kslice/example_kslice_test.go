package kslice

import (
	"fmt"
	"runtime"
	"time"
)

func Example_loopConcAsyncWaitFirst() {
	var data []int
	for i := 1; i < 1000; i++ {
		data = append(data, i)
	}

	processFunc := func(n int) (int, error) {
		delay := 1 * time.Millisecond
		time.Sleep(delay)
		return n, nil
	}
	resultCh, cancel := LoopConcAsync(data, processFunc, 100)

	firstResult := <-resultCh
	fmt.Printf("cancel前协程数量: %d\n", runtime.NumGoroutine())
	cancel()
	fmt.Printf("cancel后协程数量: %d\n", runtime.NumGoroutine())

	fmt.Printf("firstResult: %v\n", firstResult)
	// Output:
}
