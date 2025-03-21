package kmonitor

import (
	"fmt"
	"time"
)

// 示例1: 使用实时计数器统计QPS
func Example_qps() {
	// 创建一个每秒回调的实时计数器
	counter := NewRealtimeCounter[int64](time.Second, func(count int64) {
		fmt.Printf("当前QPS: %d\n", count)
	})
	defer counter.Stop()

	// 模拟请求
	counter.Add(1)
	counter.Add(1)
	counter.Add(1)

	time.Sleep(2 * time.Second)

	// Output:
	// 当前QPS: 3
	// 当前QPS: 3
}

// 示例2: 使用实时计数器统计内存使用
func Example_memoryUsage() {
	// 创建一个每5秒回调的实时计数器
	counter := NewRealtimeCounter(5*time.Second, func(bytes float64) {
		fmt.Printf("当前内存使用: %.2f MB\n", bytes/1024/1024)
	})
	defer counter.Stop()

	// 模拟内存分配和释放
	counter.Add(1024 * 1024 * 100) // 分配100MB
	counter.Dec(1024 * 1024 * 30)  // 释放30MB

	time.Sleep(6 * time.Second)

	// Output:
	// 当前内存使用: 70.00 MB
}
