// Package kcollection 提供了一些常用的集合类型和数据结构
//
// RollingWindow(from go-zero) 是一个基于时间的滑动窗口实现,主要用于限流、统计等场景。
// 它将时间划分为多个桶,每个桶记录一个时间段内的统计数据。随着时间推移,
// 旧的桶会被重置并重新使用,形成一个环形缓冲区。
//
// 主要特点:
//   - 支持泛型,可以统计任意数字类型
//   - 线程安全
//   - 自动对齐时间窗口,避免统计偏差
//   - 支持自定义桶的实现
//   - 支持配置窗口大小、时间间隔等参数
//
// 典型使用场景:
//   - 限流器: 统计单位时间内的请求次数
//   - 统计器: 计算最近一段时间内的平均值、总和等
//   - 监控指标: 收集系统性能数据
package kcollection

import (
	"sync"
	"time"

	"github.com/mtgnorton/k/kmath"
	"github.com/mtgnorton/k/ktime"
)

type (
	// BucketInterface 定义了桶的接口
	// T 为数字类型
	BucketInterface[T kmath.Number] interface {
		// Add 向桶中添加一个值
		Add(v T)
		// Reset 重置桶
		Reset()
	}

	// RollingWindow 滑动窗口
	// T 为数字类型
	// B 为实现了 BucketInterface 接口的类型
	RollingWindow[T kmath.Number, B BucketInterface[T]] struct {
		lock     sync.RWMutex
		win      *window[T, B]
		lastTime time.Duration
		offset   int // 当前桶的位置
		opts     *RollingWindowOptions[T, B]
	}
)

// NewRollingWindow 创建一个新的滑动窗口
// 参数:
//   - newBucket: 创建新桶的函数
//   - opts: 可选配置项,包括窗口大小、时间间隔等
//
// 返回:
//   - *RollingWindow: 新创建的滑动窗口
//
// 注意:
//   - 窗口大小必须大于0,否则会panic
//   - 默认窗口大小为10,时间间隔为1分钟,不忽略当前桶
//
// 示例:
//
//	rw := NewRollingWindow(func() *Bucket[int64] {
//	  return new(Bucket[int64])
//	}, WithSize[int64](10))
func NewRollingWindow[T kmath.Number, B BucketInterface[T]](newBucket func() B, opts ...RollingWindowOption[T, B]) *RollingWindow[T, B] {
	options := NewRollingWindowOptions[T, B]()
	for _, opt := range opts {
		opt(options)
	}
	if options.size < 1 {
		panic("size must be greater than 0")
	}
	w := &RollingWindow[T, B]{
		win:      newWindow(newBucket, options.size),
		lastTime: ktime.Now(),
		opts:     options,
	}
	return w
}

// Add 向当前桶中添加一个值
// 参数:
//   - v: 要添加的值
func (rw *RollingWindow[T, B]) Add(v T) {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	rw.updateOffset()
	rw.win.add(rw.offset, v)
}

// Reduce 遍历所有有效的桶
// 参数:
//   - fn: 处理每个桶的函数
//
// 注意:
//   - 如果设置了ignoreCurrent为true,则不会处理当前桶
//   - 遍历顺序为从旧到新,即最早的桶到最近的桶
func (rw *RollingWindow[T, B]) Reduce(fn func(b B)) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	var diff int
	span := rw.span()

	if span == 0 && rw.opts.ignoreCurrent {
		diff = rw.opts.size - 1
	} else {
		diff = rw.opts.size - span
	}
	if diff > 0 {
		offset := (rw.offset + span + 1) % rw.opts.size
		rw.win.reduce(offset, diff, fn)
	}
}

// GetLastValidBucket 获取最后一个有效的桶
// 返回:
//   - bucket: 最后一个有效的桶
//   - ok: 是否存在有效的桶
func (rw *RollingWindow[T, B]) GetLastValidBucket() (B, bool) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	span := rw.span()
	var diff int

	if span == 0 && rw.opts.ignoreCurrent {
		diff = rw.opts.size - 1
	} else {
		diff = rw.opts.size - span
	}

	if diff <= 0 {
		var zero B
		return zero, false // 无有效桶
	}
	// 计算最后一个有效桶的位置
	offset := (rw.offset + span + 1) % rw.opts.size
	lastPos := (offset + diff - 1) % rw.opts.size
	return rw.win.buckets[lastPos], true
}

// span 计算从上次更新到现在经过了多少个时间间隔
// 返回:
//   - int: 经过的时间间隔数
func (rw *RollingWindow[T, B]) span() int {
	offset := int(ktime.Since(rw.lastTime) / rw.opts.interval)
	if 0 <= offset && offset < rw.opts.size {
		return offset
	}

	return rw.opts.size
}

// updateOffset 更新窗口的偏移量
func (rw *RollingWindow[T, B]) updateOffset() {
	span := rw.span()
	if span <= 0 {
		return
	}

	offset := rw.offset

	for i := 0; i < span; i++ {
		rw.win.resetBucket((offset + i + 1) % rw.opts.size)
	}

	rw.offset = (offset + span) % rw.opts.size
	now := ktime.Now()

	rw.lastTime = now - (now-rw.lastTime)%rw.opts.interval
}

// Bucket 实现了BucketInterface接口的基础桶类型
type Bucket[T kmath.Number] struct {
	Sum   T     // 桶中所有值的和
	Count int64 // 桶中值的数量
}

// Add 向桶中添加一个值
func (b *Bucket[T]) Add(v T) {
	b.Sum += v
	b.Count++
}

// Reset 重置桶
func (b *Bucket[T]) Reset() {
	b.Sum = 0
	b.Count = 0
}

// window 窗口的内部实现
type window[T kmath.Number, B BucketInterface[T]] struct {
	buckets []B // 所有桶的切片
	size    int // 窗口大小
}

// newWindow 创建一个新的窗口
func newWindow[T kmath.Number, B BucketInterface[T]](newBucket func() B, size int) *window[T, B] {
	buckets := make([]B, size)
	for i := 0; i < size; i++ {
		buckets[i] = newBucket()
	}
	return &window[T, B]{
		buckets: buckets,
		size:    size,
	}
}

// add 向指定位置的桶中添加值
func (w *window[T, B]) add(offset int, v T) {
	w.buckets[offset%w.size].Add(v)
}

// reduce 从指定位置开始遍历指定数量的桶
func (w *window[T, B]) reduce(start, count int, fn func(b B)) {
	for i := 0; i < count; i++ {
		fn(w.buckets[(start+i)%w.size])
	}
}

// resetBucket 重置指定位置的桶
func (w *window[T, B]) resetBucket(offset int) {
	w.buckets[offset%w.size].Reset()
}
