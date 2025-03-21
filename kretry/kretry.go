// kretry 提供了一种基于上下文的超时和重试机制
package kretry

import (
	"context"
	"time"

	"errors"
)

var DefaultRetryTimes = 3

// ErrorFunc 错误处理函数类型
// 参数说明:
//   - error: 需要处理的错误
//
// 返回值说明:
//   - shouldStop: 是否停止重试,true表示停止重试,false表示继续重试
type ErrorFunc func(error) (shouldStop bool)

// RetryFunc 重试回调函数类型
// 参数说明:
//   - attempt: 当前重试次数
//   - err: 本次执行的错误
type RetryFunc func(attempt int, err error)

// ExecFunc 执行函数类型
// 参数说明:
//   - ctx: 上下文对象,用于控制超时和取消
//
// 返回值说明:
//   - T: 执行结果
//   - error: 执行过程中的错误
type ExecFunc[T any] func(ctx context.Context) (T, error)

type retry[T any] struct {
	opts *Options
}

// New 创建一个新的重试器
// 参数说明:
//   - opts: 可选的配置选项
//
// 返回值说明:
//   - *retry[T]: 重试器实例
//
// 举例:
//
//	retry := New[string](WithTimes(3), WithBackoff(NewBackoff()))
func New[T any](opts ...Option) *retry[T] {
	options := NewOptions()
	for _, opt := range opts {
		opt(options)
	}
	if len(options.CustomDelay) > 0 {
		if len(options.CustomDelay) != options.AttemptTimes {
			panic("CustomRetryDelay must be equal to AttemptTimes")
		}
	}
	return &retry[T]{
		opts: options,
	}
}

// Do 执行带重试的操作
// 参数说明:
//   - exec: 需要执行的函数
//
// 返回值说明:
//   - T: 执行成功时的结果
//   - error: 执行失败时的错误,包含所有重试过程中的错误信息
//
// 注意事项:
//   - 默认情况下,重试次数为3次,重试间隔为100ms 200ms 400ms
//   - 可以通过WithCustomRetryDelay设置自定义重试间隔,如果设置,则必须和重试次数一致,否则会panic
//   - 如果成功,即使之前有失败也不会返回错误
//   - ctx超时控制是不精确的,只会在重试间隔内生效,如果执行一次成功,但是该次执行时间大于ctx的超时时间,则认为成功
//   - 当ErrorHandler返回true时会立即停止重试
//   - 当重试一直失败,所有的错误会通过 errors.Join 合并返回
//
// 举例:
//
//	result, err := retry.Do(func(ctx context.Context) (string, error) {
//	    return "hello", nil
//	})
func (r *retry[T]) Do(exec ExecFunc[T]) (T, error) {
	var result T
	var errs []error
	if r.opts.Ctx.Err() != nil {
		return result, r.opts.Ctx.Err()
	}
	for attempt := 0; attempt < r.opts.AttemptTimes; attempt++ {
		result, err := exec(r.opts.Ctx)
		if err == nil {
			return result, nil // 成功立即返回
		}
		// 错误处理流程
		if r.opts.ErrorHandler != nil && r.opts.ErrorHandler(err) {
			return result, err
		}
		errs = append(errs, err)

		// 执行重试回调
		if r.opts.RetryHandler != nil {
			r.opts.RetryHandler(attempt, err)
		}

		// 使用可取消的定时器避免资源泄漏
		var delay time.Duration
		if len(r.opts.CustomDelay) > 0 {
			delay = r.opts.CustomDelay[attempt]
		} else {
			delay = r.opts.Backoff.Duration()
		}
		timer := time.NewTimer(delay)
		select {
		case <-r.opts.Ctx.Done():
			timer.Stop()
			errs = append(errs, r.opts.Ctx.Err())
			return result, mergeErrors(errs)
		case <-timer.C:
			timer.Stop()
		}
	}

	return result, mergeErrors(errs)
}

// Do 执行带重试的函数调用
//
// 参数说明:
//   - exec: 需要执行的函数
//   - opts: 重试选项配置
//
// 返回值说明:
//   - T: 执行成功时的返回值
//   - error: 执行失败时的错误信息
//
// 参见 retry.Do
// 举例:
//
//	result, err := Do(func(ctx context.Context) (int, error) {
//	    return 42, nil
//	})
func Do[T any](exec ExecFunc[T], opts ...Option) (T, error) {
	r := New[T](opts...)
	return r.Do(exec)
}

// mergeErrors 合并多个错误信息
// 参数说明:
//   - errs: 错误列表
//
// 返回值说明:
//   - error: 合并后的错误
//
// 注意事项:
//   - 如果错误列表为空,返回nil
//   - 使用errors.Wrap逐层包装错误信息
func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}
