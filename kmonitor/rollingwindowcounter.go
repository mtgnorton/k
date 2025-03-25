package kmonitor

import (
	"fmt"
	"time"

	"github.com/mtgnorton/k/kcollection"
	"github.com/mtgnorton/k/kmath"
)

// RollingResultCounter 滚动结果计数器,用于统计成功和失败的请求及其消耗时间
// 支持泛型,可以统计任意数字类型
type RollingResultCounter[T kmath.Number] struct {
	successWindow *kcollection.RollingWindow[T, *kcollection.Bucket[T]]
	failWindow    *kcollection.RollingWindow[T, *kcollection.Bucket[T]]
}

// NewRollingResultCounter 创建一个新的滚动结果计数器
// 参数:
//   - opts: 可选配置项,包括窗口大小、时间间隔等
//
// 返回:
//   - *RollingResultCounter[T]: 新创建的滚动结果计数器
//
// 注意:
//   - T 必须是数字类型
//
// 示例:
//
//	counter := NewRollingResultCounter[int64]()
//	counter.AddSuccess(100)
func NewRollingResultCounter[T kmath.Number](opts ...kcollection.RollingWindowOption[T, *kcollection.Bucket[T]]) *RollingResultCounter[T] {
	opt := kcollection.NewRollingWindowOptions[T, *kcollection.Bucket[T]]()
	for _, o := range opts {
		o(opt)
	}
	r := &RollingResultCounter[T]{
		successWindow: kcollection.NewRollingWindow(func() *kcollection.Bucket[T] {
			return &kcollection.Bucket[T]{}
		}, opts...),
		failWindow: kcollection.NewRollingWindow(func() *kcollection.Bucket[T] {
			return &kcollection.Bucket[T]{}
		}, opts...),
	}
	return r
}

// AddSuccess 添加一个成功请求的记录
// 参数:
//   - consumeTime: 请求消耗的时间
func (r *RollingResultCounter[T]) AddSuccess(consumeTime T) {
	r.successWindow.Add(consumeTime)
}

// AddFail 添加一个失败请求的记录
// 参数:
//   - consumeTime: 请求消耗的时间
func (r *RollingResultCounter[T]) AddFail(consumeTime T) {
	r.failWindow.Add(consumeTime)
}

// Reduce 遍历所有有效的桶并执行回调函数
// 参数:
//   - successFn: 处理成功请求桶的函数,接收成功请求数量和总消耗时间
//   - failFn: 处理失败请求桶的函数,接收失败请求数量和总消耗时间
//
// 注意:
//   - 如果设置了ignoreCurrent为true,则不会处理当前桶
//
// 示例:
//
//	counter.Reduce(
//	  func(sc int64, st int64) { fmt.Println("成功:", sc, st) },
//	  func(fc int64, ft int64) { fmt.Println("失败:", fc, ft) },
//	)
func (r *RollingResultCounter[T]) Reduce(successFn func(successCount int64, successConsumeTime T), failFn func(failCount int64, failConsumeTime T)) {
	r.successWindow.Reduce(func(b *kcollection.Bucket[T]) {
		successFn(b.Count, b.Sum)
	})
	r.failWindow.Reduce(func(b *kcollection.Bucket[T]) {
		failFn(b.Count, b.Sum)
	})
}

// Info 获取计数器的详细信息
// 返回:
//   - string: 包含成功和失败请求的详细统计信息
func (r *RollingResultCounter[T]) Info() string {
	info := "successInfo:\n"
	size := r.successWindow.Opts.Size
	interval := r.successWindow.Opts.Interval
	// size = 5  (5-1)*interval -> 5*interval
	// ...
	// size = 1  (1-1)*interval -> 1 *interval
	temp := make([]struct {
		successCount          int64
		avgSuccessConsumeTime string
		failCount             int64
		avgFailConsumeTime    string
	}, size)

	r.successWindow.Reduce(func(b *kcollection.Bucket[T]) {
		d := "-"
		if b.Count > 0 {
			d = fmt.Sprintf("%v", float64(b.Sum)/float64(b.Count))
		}
		size--
		temp[size] = struct {
			successCount          int64
			avgSuccessConsumeTime string
			failCount             int64
			avgFailConsumeTime    string
		}{
			successCount:          b.Count,
			avgSuccessConsumeTime: d,
		}
	})
	size = r.failWindow.Opts.Size
	r.failWindow.Reduce(func(b *kcollection.Bucket[T]) {
		size--
		d := "-"
		if b.Count > 0 {
			d = fmt.Sprintf("%v", float64(b.Sum)/float64(b.Count))
		}
		temp[size].failCount = b.Count
		temp[size].avgFailConsumeTime = d
	})

	for i := 0; i < len(temp); i++ {
		info += fmt.Sprintf(" time:%v-%v,successCount: %v, successAvgConsumeTime: %v,failCount: %v, failAvgConsumeTime: %v\n", time.Duration(i)*interval, time.Duration(i+1)*interval, temp[i].successCount, temp[i].avgSuccessConsumeTime, temp[i].failCount, temp[i].avgFailConsumeTime)
	}
	return info
}
