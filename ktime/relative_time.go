package ktime

import "time"

// initTime 是系统启动时间, 用于计算相对时间,使相对时间的计算不依赖系统时间,避免出现时间差为0或负数
// 解决以下3个问题
// 1. 系统时间被重置的问题
// 2. 避免高频调用下时间差为零的问题
// 3. 依赖单调时钟（Monotonic Clock）
//
// go中的单调时钟
//
//	time.Now() 返回的值是单调时钟,所以initTime 包含单调时钟
//	当调用 time.Since(t) 时：
//	如果 t 包含单调时钟读数（即通过 time.Now() 获取的时间），则优先使用单调时钟计算时间差。
//	如果 t 不包含单调时钟（例如通过 time.Parse 解析的时间），则使用系统时钟计算时间差。
//
// 高频调用 Now()：
//
//	t1 := ktime.Now() // 返回 time.Since(initTime) 的差值（例如 1年1个月 + 1ms）
//	t2 := ktime.Now() // 返回 1年1个月 + 2ms
//
// 即使两次调用发生在同一系统时间点（例如系统时钟未更新），单调时钟仍会确保 t2 > t1。
var initTime = time.Now().AddDate(-1, -1, -1)

// Now 返回相对于系统启动时间的时间
func Now() time.Duration {
	return time.Since(initTime)
}

// Since 返回相对于系统启动时间的时间减去d
func Since(d time.Duration) time.Duration {
	return time.Since(initTime) - d
}
