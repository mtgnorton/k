package kmonitor

import (
	"fmt"
)

// 示例1: 使用实时计数器统计QPS
func Example_qps() {
	// 创建一个每秒回调的实时计数器
	counter := NewRealtimeCounter[int64]()

	// 模拟请求
	counter.Add(1)
	counter.Add(1)
	counter.Add(1)

	fmt.Println(counter.Get())
	// Output:
	// 3
}

// 示例2: 使用实时计数器统计内存使用
func Example_memoryUsage() {
	// 创建一个每5秒回调的实时计数器
	counter := NewRealtimeCounter[int64]()

	// 模拟内存分配和释放
	counter.Add(1024 * 1024 * 100) // 分配100MB
	counter.Dec(1024 * 1024 * 30)  // 释放30MB

	fmt.Println(counter.Get())
	// Output:
	// 73400320
}
