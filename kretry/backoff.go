package kretry

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

const maxInt64 = float64(math.MaxInt64 - 512)

// Backoff 实现指数退避算法
type Backoff struct {
	attempt atomic.Uint64 // 当前尝试次数
	opts    *BackOffOptions
}

// NewBackoff 创建新的Backoff实例
//
// 参数说明:
//   - opts: 可选的配置函数，用于自定义Backoff参数
//
// 返回值说明:
//   - *Backoff: 返回初始化后的Backoff实例
//
// 注意事项:
//   - 默认参数: factor=2, jitter=false, min=100ms, max=10s
//   - 默认参数下的回退时间为: 10次序列为100ms 200ms 400ms 800ms 1.6s 3.2s 6.4s 10s 10s 10s
//   - 可以通过WithFactor, WithJitter, WithMin, WithMax等函数自定义参数
//
// 示例:
//
//	b := NewBackoff(WithFactor(1.5), WithJitter(true))
func NewBackoff(opts ...BackoffOption) *Backoff {
	options := NewBackOffOptions()
	for _, opt := range opts {
		opt(options)
	}
	b := &Backoff{
		opts: options,
	}
	return b
}

// Duration 计算并返回当前尝试的退避时间
//
// 返回值说明:
//   - time.Duration: 返回计算后的退避时间
//
// 注意事项:
//   - 每次调用都会增加尝试次数
//   - 返回的时间会在min和max之间
//
// 示例:
//
//	d := b.Duration() // 获取当前退避时间
func (b *Backoff) Duration() time.Duration {
	d := b.ForAttempt(float64(b.attempt.Add(1) - 1))
	return d
}

// ForAttempt 根据尝试次数计算退避时间
//
// 参数说明:
//   - attempt: 尝试次数
//
// 返回值说明:
//   - time.Duration: 返回计算后的退避时间
//
// 注意事项:
//   - 如果启用了jitter，返回的时间会有随机波动
//   - 返回的时间不会超过maxInt64
//
// 示例:
//
//	d := b.ForAttempt(3) // 计算第3次尝试的退避时间
func (b *Backoff) ForAttempt(attempt float64) time.Duration {
	min := b.opts.min
	if min <= 0 {
		min = 100 * time.Millisecond
	}
	max := b.opts.max
	if max <= 0 {
		max = 10 * time.Second
	}
	if min >= max {
		return max
	}
	factor := b.opts.factor
	if factor <= 0 {
		factor = 2
	}
	minTime := float64(min)
	duration := minTime * math.Pow(factor, attempt)
	if b.opts.jitter {
		duration = rand.Float64()*(duration-minTime) + minTime
	}
	if duration > maxInt64 {
		return max
	}
	dur := time.Duration(duration)

	if dur < min {
		return min
	}
	if dur > max {
		return max
	}
	return dur
}

// Reset 重置尝试次数
//
// 注意事项:
//   - 重置后下次调用Duration将从第一次尝试开始计算
func (b *Backoff) Reset() {
	b.attempt.Store(0)
}

// Attempt 获取当前尝试次数
//
// 返回值说明:
//   - float64: 返回当前尝试次数
func (b *Backoff) Attempt() float64 {
	return float64(b.attempt.Load())
}

// Copy 创建当前Backoff的副本
//
// 返回值说明:
//   - *Backoff: 返回新的Backoff实例
//
// 注意事项:
//   - 新实例的尝试次数会被重置为0
func (b *Backoff) Copy() *Backoff {
	return &Backoff{
		opts: b.opts,
	}
}
