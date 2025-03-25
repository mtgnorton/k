package kmonitor

import (
	"fmt"
	"time"
)

func ExampleSampling() {

	rch, clear := Sampling(100*time.Millisecond, 10, func(item int) {
		fmt.Println(item)
	})

	for i := 0; i < 10; i++ {
		rch <- i
	}
	time.Sleep(time.Second)
	clear()
	// Output:
	// 9
}

func ExampleConsumeTimeStatistics() {
	stats := ConsumeTimeStatistics("MyProcess")

	time.Sleep(100 * time.Millisecond)
	fmt.Println(stats("步骤1"))

	time.Sleep(50 * time.Millisecond)
	fmt.Println(stats("步骤2"))

	// Output:
	// 输出为随机值,如:
	// [MyProcess] 步骤1: Total Time: 101.059333ms Interval Time: 101.059333ms
	// [MyProcess] 步骤2: Total Time: 151.1685ms Interval Time: 50.109167ms
}
