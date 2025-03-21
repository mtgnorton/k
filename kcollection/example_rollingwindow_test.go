package kcollection

import (
	"fmt"
	"time"
)

// 示例1: 使用滑动窗口统计每分钟的请求数量
func Example_requestCounter() {
	// 创建一个窗口大小为5,每个桶统计1分钟数据的滑动窗口
	rw := NewRollingWindow(func() *Bucket[int64] {
		return new(Bucket[int64])
	}, WithSize[int64, *Bucket[int64]](5), WithInterval[int64, *Bucket[int64]](time.Minute))

	// 模拟请求
	rw.Add(1) // 第1分钟有1个请求
	rw.Add(1) // 第1分钟有1个请求

	var i = 5
	// 统计结果
	rw.Reduce(func(b *Bucket[int64]) {
		fmt.Printf("距离当前时间%d分钟,请求数量: %d\n", i, b.Count)
		i--
	})

	// Output:
	// 距离当前时间5分钟,请求数量: 0
	// 距离当前时间4分钟,请求数量: 0
	// 距离当前时间3分钟,请求数量: 0
	// 距离当前时间2分钟,请求数量: 0
	// 距离当前时间1分钟,请求数量: 2
}

// 示例2: 使用滑动窗口计算最近5秒的平均响应时间
func Example_responseTimeAvg() {
	// 创建一个窗口大小为5,每个桶统计1秒数据的滑动窗口
	rw := NewRollingWindow[float64, *Bucket[float64]](func() *Bucket[float64] {
		return new(Bucket[float64])
	}, WithSize[float64, *Bucket[float64]](5), WithInterval[float64, *Bucket[float64]](time.Second))

	// 记录响应时间(ms)
	rw.Add(100) // 第1秒响应时间100ms
	rw.Add(200) // 第1秒响应时间200ms
	rw.Add(150) // 第1秒响应时间150ms

	// 计算平均响应时间
	rw.Reduce(func(b *Bucket[float64]) {
		if b.Count > 0 {
			fmt.Printf("平均响应时间: %.2fms\n", b.Sum/float64(b.Count))
		}
	})

	// Output:
	// 平均响应时间: 150.00ms
}

// 示例3: 使用滑动窗口实现限流器
func Example_rateLimiter() {
	// 创建一个窗口大小为1,每个桶统计1秒数据的滑动窗口
	rw := NewRollingWindow[int64, *Bucket[int64]](func() *Bucket[int64] {
		return new(Bucket[int64])
	}, WithSize[int64, *Bucket[int64]](1), WithInterval[int64, *Bucket[int64]](time.Second))

	// 限流规则: 每秒最多允许100个请求
	const limit int64 = 100

	// 判断是否允许请求通过
	isAllowed := func() bool {
		var count int64
		rw.Reduce(func(b *Bucket[int64]) {
			count = b.Count
		})
		return count < limit
	}

	// 模拟请求
	if isAllowed() {
		rw.Add(1)
		fmt.Println("请求通过")
	} else {
		fmt.Println("请求被限流")
	}

	// Output:
	// 请求通过
}
