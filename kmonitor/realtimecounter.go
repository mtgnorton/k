package kmonitor

import (
	"sync"
	"time"

	"github.com/mtgnorton/k/kmath"
)

// RealtimeCounter 实时计数器,用于统计一段时间内的计数值
// 支持泛型,可以统计任意数字类型
type RealtimeCounter[T kmath.Number] struct {
	counter          T
	mu               sync.Mutex
	callbackDuration time.Duration
	callback         func(T)
	stopCh           chan struct{}
}

// NewRealtimeCounter 创建一个新的实时计数器
// 参数:
//   - duration: 回调时间间隔
//   - callback: 回调函数,用于处理计数值
//
// 返回:
//   - *RealtimeCounter[T]: 实时计数器实例
//
// 示例:
//
//	counter := NewRealtimeCounter[int64](time.Second, func(count int64) {
//	    fmt.Printf("当前计数: %d\n", count)
//	})
func NewRealtimeCounter[T kmath.Number](duration time.Duration, callback func(T)) *RealtimeCounter[T] {
	r := &RealtimeCounter[T]{
		callbackDuration: duration,
		callback:         callback,
		stopCh:           make(chan struct{}),
	}
	go r.InfiniteRun()
	return r
}

// Add 增加计数值
// 参数:
//   - v: 要增加的值
func (r *RealtimeCounter[T]) Add(v T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter += v
}

// Dec 减少计数值
// 参数:
//   - v: 要减少的值
func (r *RealtimeCounter[T]) Dec(v T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter -= v
}

// Get 获取当前计数值
// 返回:
//   - T: 当前计数值
func (r *RealtimeCounter[T]) Get() T {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.counter
}

// Reset 重置计数值为0
func (r *RealtimeCounter[T]) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	var v T
	r.counter = v
}

// InfiniteRun 启动计数器的无限循环
// 注意:
//   - 该函数会在新的goroutine中运行
//   - 通过stopCh来控制停止
func (r *RealtimeCounter[T]) InfiniteRun() {
	ticker := time.NewTicker(r.callbackDuration)
	defer ticker.Stop()
	for {
		select {
		case <-r.stopCh:
			return
		case <-ticker.C:
			r.mu.Lock()
			r.callback(r.counter)
			r.mu.Unlock()
		}
	}
}

// Stop 停止计数器
// 注意:
//   - 停止后不能重新启动
//   - 需要新建计数器才能继续使用
func (r *RealtimeCounter[T]) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	close(r.stopCh)
}
