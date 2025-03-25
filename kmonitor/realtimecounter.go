package kmonitor

import (
	"sync"

	"github.com/mtgnorton/k/kmath"
)

// RealtimeCounter 实时计数器,用于统计一段时间内的计数值
// 支持泛型,可以统计任意数字类型
type RealtimeCounter[T kmath.Number] struct {
	counter T
	mu      sync.Mutex
}

// NewRealtimeCounter 创建一个新的实时计数器
// 参数:
//   - 无
//
// 返回:
//   - *RealtimeCounter[T]: 新创建的实时计数器
//
// 注意:
//   - T 必须是数字类型
//
// 示例:
//
//	counter := NewRealtimeCounter[int64]()
//	counter.Add(10)
func NewRealtimeCounter[T kmath.Number]() *RealtimeCounter[T] {
	r := &RealtimeCounter[T]{}
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
