package kretry

import (
	"context"
	"time"
)

type Options struct {
	Ctx          context.Context // 当Ctx设置了超时时间, 则当Ctx超时后, 会停止重试
	ErrorHandler ErrorFunc       // 错误处理回调函数
	RetryHandler RetryFunc       // 重试时调用的函数
	AttemptTimes int             // 重试次数
	CustomDelay  []time.Duration // 自定义重试间隔时间,必须和重试次数一致
	Backoff      *Backoff        // 退避策略

}

type Option func(o *Options)

func NewOptions() *Options {
	return &Options{
		Ctx:          context.Background(),
		AttemptTimes: DefaultRetryTimes,
		Backoff:      NewBackoff(),
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Ctx = ctx
	}
}

func WithErrHandler(errHandle func(error) (stop bool)) Option {
	return func(o *Options) {
		o.ErrorHandler = errHandle
	}
}

func WithRetryHandler(retryHandler func(attempt int, err error)) Option {
	return func(o *Options) {
		o.RetryHandler = retryHandler
	}
}

func WithTimes(times int) Option {
	return func(o *Options) {
		o.AttemptTimes = times
	}
}

func WithCustomDelay(delay []time.Duration) Option {
	return func(o *Options) {
		o.CustomDelay = delay
	}
}

func WithBackoff(backoff *Backoff) Option {
	return func(o *Options) {
		o.Backoff = backoff
	}
}

type BackOffOptions struct {
	factor float64       // 指数因子
	jitter bool          // 是否添加随机抖动
	min    time.Duration // 最小退避时间
	max    time.Duration // 最大退避时间
}

// BackoffOption 用于配置Backoff的选项函数类型
type BackoffOption func(b *BackOffOptions)

func NewBackOffOptions() *BackOffOptions {
	return &BackOffOptions{
		factor: 2,
		jitter: false,
		min:    100 * time.Millisecond,
		max:    10 * time.Second,
	}
}

// WithFactor 设置指数因子
//
// 参数说明:
//   - factor: 指数因子，必须大于0
func WithFactor(factor float64) BackoffOption {
	return func(b *BackOffOptions) {
		b.factor = factor
	}
}

// WithJitter 设置是否添加随机抖动
//
// 参数说明:
//   - jitter: 是否启用随机抖动
func WithJitter(jitter bool) BackoffOption {
	return func(b *BackOffOptions) {
		b.jitter = jitter
	}
}

// WithMin 设置最小退避时间
//
// 参数说明:
//   - min: 最小退避时间，必须大于0
func WithMin(min time.Duration) BackoffOption {
	return func(b *BackOffOptions) {
		b.min = min
	}
}

// WithMax 设置最大退避时间
//
// 参数说明:
//   - max: 最大退避时间，必须大于0
func WithMax(max time.Duration) BackoffOption {
	return func(b *BackOffOptions) {
		b.max = max
	}
}
