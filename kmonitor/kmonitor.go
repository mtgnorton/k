// kmonitor 提供用于监控和统计的工具函数
// 采样: 对输入数据进行采样处理
// 统计: 统计任务执行时间
// 超时: 监控超时
package kmonitor

import (
	"fmt"
	"time"
)

// Sampling 对输入数据进行采样处理
//
// 参数说明:
//   - duration: 采样时间间隔，如果为0则只根据数量触发
//   - amount: 采样数量，如果为0则只根据时间触发
//   - exec: 处理采样数据的函数
//
// 返回值说明:
//   - rch: 用于接收数据的通道
//   - clear: 用于关闭采样和清理资源的函数
//
// 注意事项:
//   - duration和amount不能同时为0
//   - 使用带缓冲的信号量控制并发，最大并发数为100
//   - 当达到采样条件时，会重置计数器和时间
//   - 需要调用clear函数来关闭通道和清理资源
//
// 示例:
//
//	rch, clear := Sampling(100*time.Millisecond, 10, func(item int) {
//	    fmt.Println(item)
//	})
//	defer clear()
//	rch <- 1
func Sampling[T any](duration time.Duration, amount int, exec func(T)) (rch chan<- T, clear func()) {
	ch := make(chan T)
	sem := make(chan struct{}, 100)
	if duration <= 0 && amount <= 0 {
		panic("至少需要设置 duration 或 amount 其中一个参数")
	}
	var (
		counter      int
		startTime    = time.Now()
		timeTrigger  = duration > 0
		countTrigger = amount > 0
	)

	go func() {
		defer close(sem)
		for item := range ch {
			counter++
			triggered := false
			if countTrigger && counter >= amount {
				triggered = true
			}
			if timeTrigger && time.Since(startTime) >= duration {
				triggered = true
			}

			if triggered {
				sem <- struct{}{}
				go func(item T) {
					defer func() { <-sem }()
					exec(item)
				}(item)
				counter = 0
				startTime = time.Now()
			}
		}
	}()
	return ch, func() {
		close(ch)
	}
}

// ConsumeTimeStatistics 用于统计任务执行时间
//
// 参数说明:
//   - name: 任务名称，用于标识统计结果
//
// 返回值说明:
//   - 返回一个函数，该函数接收label参数并返回统计信息字符串
//
// 注意事项:
//   - 统计从调用consumeTimeStatistic时开始
//   - 每次调用返回的函数都会更新最后一次统计时间
//   - 返回的字符串包含总时间和间隔时间
//
// 示例:
//
//	stats := ConsumeTimeStatistics("MyTask")
//	time.Sleep(100 * time.Millisecond)
//	fmt.Println(stats("Step1"))
func ConsumeTimeStatistics(name string) func(label string) string {
	startTime := time.Now()
	lastTime := startTime
	return func(label string) string {
		now := time.Now()
		totalDuration := now.Sub(startTime)
		intervalDuration := now.Sub(lastTime)
		lastTime = now
		return fmt.Sprintf("[%s] %s: Total Time: %s Interval Time: %s",
			name, label, totalDuration, intervalDuration)
	}
}
