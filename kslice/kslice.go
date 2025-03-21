// kslice 提供了一些对slice的一些常用操作
//
// 包含并发遍历,分割并发遍历,转换,过滤等
package kslice

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"k/kmath"
	"k/kreflect"
)

type Result[T any, V any] struct {
	Key    int
	Item   T
	Result V
	Error  error
}

// LoopConc 并发遍历slice中的每个元素
//
// 参数说明:
//   - s: 需要遍历的slice
//   - fn: 处理每个元素的函数，接收元素索引和元素值作为参数
//   - concurrency: 可选参数，控制并发数，默认为1
//
// 返回值说明:
//
//	无返回值
//
// 注意事项:
//   - 该函数会阻塞直到所有并发任务完成
//   - 如果concurrency参数小于等于0，并发数会被设置为1
//   - 每个元素都会在一个独立的goroutine中处理
//   - 使用sync.WaitGroup和channel来控制并发数
func LoopConc[T any](s []T, fn func(index int, item T), concurrency ...int) {
	wg := sync.WaitGroup{}
	if len(concurrency) == 0 {
		concurrency = []int{1}
	}
	ch := make(chan struct{}, concurrency[0])
	for i, item := range s {
		wg.Add(1)
		ch <- struct{}{}
		go func(i int, item T) {
			defer func() {
				wg.Done()
				<-ch
			}()
			fn(i, item)
		}(i, item)
	}
	wg.Wait()
}

// LoopConcAsyncFirstSuccess 异步并发处理切片中的每个元素,返回第一个成功的结果
//
// 参数说明:
//   - s: 需要处理的切片
//   - exec: 处理每个元素的函数,接收元素值并返回结果和可能的错误
//   - concurrency: 可选参数,控制并发数,默认为1
//
// 返回值说明:
//   - Result[T, V]: 第一个成功的处理结果
//   - error: 处理过程中的错误
//
// 注意事项:
//   - 该函数会阻塞直到找到第一个成功的结果或所有元素处理完成
//   - 如果所有元素处理都失败,返回空结果和所有错误,错误通过errors.Join返回
//   - 找到第一个成功结果后会自动取消其他正在进行的任务
//   - 内部使用LoopConcAsync实现,继承其并发控制特性
//
// 示例:
//
//	data := []int{1, 2, 3}
//	result, err := LoopConcAsyncFirstSuccess(data, func(n int) (int, error) {
//	    return n * 2, nil
//	})
//	if err != nil {
//	    fmt.Println("错误:", err)
//	} else {
//	    fmt.Printf("结果: %d\n", result.Result)
//	}
func LoopConcAsyncFirstSuccess[T any, V any](
	s []T,
	exec func(T) (V, error),
	concurrency ...int,
) (Result[T, V], error) {
	var errs []error
	ch, cancel := LoopConcAsync(s, exec, concurrency...)
	defer cancel()
	for result := range ch {
		if result.Error == nil {
			return result, nil
		}
		errs = append(errs, result.Error)
	}
	return Result[T, V]{}, errors.Join(errs...)
}

// LoopConcAsync 异步并发处理切片中的每个元素并返回结果
//
// 参数说明:
//   - s: 需要处理的切片
//   - exec: 处理每个元素的函数，接收元素值并返回结果和可能的错误
//   - concurrency: 可选参数，控制并发数，默认为1
//
// 返回值说明:
//   - <-chan Result[T, V]: 结果通道，包含处理结果和可能的错误
//   - func(): 取消函数，用于提前终止所有并发任务
//
// 注意事项:
//   - 该函数不会阻塞，而是立即返回结果通道和取消函数
//   - 如果concurrency参数小于等于0，并发数会被设置为1
//   - 每个元素都会在一个独立的goroutine中处理
//   - 处理过程中的panic会被捕获并作为错误返回
//   - 调用取消函数后，所有正在进行的任务会被终止,如果exec函数一直阻塞,无法完成,会导致goroutine泄露
//   - 结果通道会在所有任务完成后自动关闭
//
// 示例:
//
//	data := []int{1, 2, 3}
//	resultCh, cancel := LoopConcAsync(data, func(n int) (int, error) {
//	    return n * 2, nil
//	}, 2)
//	defer cancel() // 确保资源释放
//
//	for result := range resultCh {
//	    if result.Error != nil {
//	        fmt.Println("错误:", result.Error)
//	    } else {
//	        fmt.Printf("索引 %d 的结果: %d\n", result.Key, result.Result)
//	    }
//	}
func LoopConcAsync[T any, V any](
	s []T,
	exec func(T) (V, error),
	concurrency ...int,
) (<-chan Result[T, V], func()) {
	conc := 1
	if len(concurrency) > 0 && concurrency[0] > 0 {
		conc = concurrency[0]
	}

	concCh := make(chan struct{}, conc)
	resultCh := make(chan Result[T, V])
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	var isCancel atomic.Bool // 新增
	go func() {
		defer func() {
			close(resultCh)
		}()
		for idx, item := range s {
			concCh <- struct{}{}
			if isCancel.Load() {
				return
			}
			wg.Add(1)
			go func(item T, index int) {
				defer func() {
					<-concCh
					wg.Done()
				}()
				var result Result[T, V]
				func() {
					defer func() {
						if r := recover(); r != nil {
							result.Error = fmt.Errorf("panic: %v, item: %+v, index: %d", r, item, index)
						}
					}()
					v, err := exec(item)
					result = Result[T, V]{
						Key:    index,
						Item:   item,
						Result: v,
						Error:  err,
					}
				}()
				if isCancel.Load() {
					return
				}
				select {
				case resultCh <- result:
				case <-ctx.Done():
				}
			}(item, idx)
		}
		wg.Wait()
	}()

	return resultCh, func() {
		isCancel.Store(true)
		cancel()
	}
}

// ChunkConc 将slice分块并发处理
//
// 参数说明:
//   - s: 需要处理的slice
//   - size: 每个分块的大小
//   - fn: 处理每个分块的函数，接收分块作为参数
//   - concNumber: 可选参数，控制并发数，默认为1
//
// 返回值说明:
//
//	无返回值
//
// 注意事项:
//   - 该函数会阻塞直到所有并发任务完成
//   - 如果size参数小于等于0，函数直接返回
//   - 如果slice为空，函数直接返回
//   - 如果concNumber参数小于等于0，并发数会被设置为1
//   - 使用sync.WaitGroup和channel来控制并发数
//   - 每个分块都会在一个独立的goroutine中处理
func ChunkConc[T any](s []T, size int, fn func(chunk []T), concNumber ...int) {
	if len(concNumber) == 0 {
		concNumber = []int{1}
	}
	if size <= 0 {
		return
	}
	if len(s) == 0 {
		return
	}
	length := len(s)

	wg := sync.WaitGroup{}
	ch := make(chan struct{}, concNumber[0])

	for i := 0; i < length; i += size {
		end := kmath.Min(i+size, length)
		chunk := s[i:end]
		wg.Add(1)
		ch <- struct{}{}
		go func(chunk []T) {
			defer func() {
				wg.Done()
				<-ch
			}()
			fn(chunk)
		}(chunk)
	}
	wg.Wait()
}

// ToMap 将slice转换为map
//
// 参数说明:
//   - s: 需要转换的slice
//   - fn: 转换函数，接收索引和元素作为参数，返回key和value
//
// 返回值说明:
//   - 返回转换后的map
//
// 注意事项:
//   - 如果slice中有重复的key，后面的value会覆盖前面的value
//   - 如果slice为空，返回空的map
//   - key的类型必须是可比较的(comparable)
func ToMap[T any, K comparable](s []T, fn func(index int, item T) (key K, value T)) map[K]T {
	m := make(map[K]T)
	for i, item := range s {
		key, value := fn(i, item)
		m[key] = value
	}
	return m
}

// ItemToSlice 将切片中的每个元素转换为新类型的切片
//
// 参数说明:
//   - s: 原始切片
//   - fn: 转换函数，接收元素索引和元素值，返回转换后的值
//
// 返回值说明:
//   - []V: 转换后的新切片
//
// 注意事项:
//   - 返回的新切片长度与原始切片相同
//   - 转换函数fn不能为nil
//
// 示例:
//
//	// 将[]int转换为[]string
//	nums := []int{1, 2, 3}
//	strs := ItemToSlice(nums, func(i int, n int) string {
//	    return fmt.Sprintf("num%d", n)
//	})
//	// strs = []string{"num1", "num2", "num3"}
func ItemToSlice[T any, V any](s []T, fn func(index int, item T) V) []V {
	result := make([]V, 0, len(s))
	for i, item := range s {
		result = append(result, fn(i, item))
	}
	return result
}

// Filter 根据条件过滤切片中的元素
//
// 参数说明:
//   - s: 需要过滤的切片
//   - fn: 可选参数，过滤条件函数，接收元素索引和元素值，返回bool值
//
// 返回值说明:
//   - []T: 过滤后的新切片
//
// 注意事项:
//   - 如果未提供过滤函数，则默认过滤掉nil值
//   - 返回的新切片长度可能小于原切片
//
// 示例:
//
//	nums := []int{1, 2, 3, 4}
//	evens := Filter(nums, func(i int, n int) bool {
//	    return n%2 == 0
//	})
//	// evens = []int{2, 4}
func Filter[T any](s []T, fn ...func(index int, item T) bool) []T {
	result := make([]T, 0, len(s))
	for i, item := range s {
		if len(fn) == 0 && !kreflect.IsNil(item) {
			result = append(result, item)
		} else {
			if fn[0](i, item) {
				result = append(result, item)
			}
		}
	}
	return result
}

// FilterRepeat 去除切片中的重复元素
//
// 参数说明:
//   - s: 需要去重的切片
//
// 返回值说明:
//   - []T: 去重后的新切片
//
// 注意事项:
//   - 元素类型必须实现comparable接口
//   - 返回的新切片长度可能小于原切片
//   - 保留第一个出现的元素
//
// 示例:
//
//	nums := []int{1, 2, 2, 3}
//	unique := FilterRepeat(nums)
//	// unique = []int{1, 2, 3}
func FilterRepeat[T comparable](s []T) []T {
	m := make(map[T]struct{})
	result := make([]T, 0, len(s))
	for _, item := range s {
		if _, ok := m[item]; !ok {
			m[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// RemoveElements 根据条件移除切片中的多个元素
//
// 参数说明:
//   - slice: 需要处理的切片
//   - condition: 移除条件函数，接收元素值，返回bool值
//   - keepOrder: 可选参数，是否保持元素顺序，默认为false
//
// 返回值说明:
//   - []T: 移除元素后的新切片
//
// 注意事项:
//   - 如果keepOrder为true，会保持剩余元素的原始顺序
//   - 如果keepOrder为false，会使用更高效的交换方式
//   - 返回的新切片长度可能小于原切片
//
// 示例:
//
//	nums := []int{1, 2, 3, 4}
//	result := RemoveElements(nums, func(n int) bool {
//	    return n%2 == 0
//	}, true)
//	// result = []int{1, 3}
func RemoveElements[T any](slice []T, condition func(T) bool, keepOrder ...bool) []T {
	keepOrderFlag := len(keepOrder) > 0 && keepOrder[0]

	if keepOrderFlag {
		writeIdx := 0
		for _, v := range slice {
			if !condition(v) {
				slice[writeIdx] = v
				writeIdx++
			}
		}
		tail := slice[writeIdx:]
		for i := range tail {
			var zero T
			tail[i] = zero
		}

		return slice[:writeIdx]
	} else {
		for i := len(slice) - 1; i >= 0; i-- {
			if condition(slice[i]) {
				slice[i] = slice[len(slice)-1]
				slice = slice[:len(slice)-1]
			}
		}
		return slice
	}
}

// RemoveElement 根据条件移除切片中的第一个匹配元素
//
// 参数说明:
//   - slice: 需要处理的切片
//   - condition: 移除条件函数，接收元素值，返回bool值
//   - keepOrder: 可选参数，是否保持元素顺序，默认为false
//
// 返回值说明:
//   - []T: 移除元素后的新切片
//
// 注意事项:
//   - 如果keepOrder为true，会保持剩余元素的原始顺序
//   - 如果keepOrder为false，会使用更高效的交换方式
//   - 只移除第一个匹配的元素
//   - 返回的新切片长度比原切片小1
//
// 示例:
//
//	nums := []int{1, 2, 3, 2}
//	result := RemoveElement(nums, func(n int) bool {
//	    return n == 2
//	})
//	// result = []int{1, 3, 2}
func RemoveElement[T any](slice []T, condition func(T) bool, keepOrder ...bool) []T {
	keepOrderFlag := len(keepOrder) > 0 && keepOrder[0]
	if len(slice) == 0 {
		return slice
	}
	if keepOrderFlag {
		for i, v := range slice {
			if condition(v) {
				return append(slice[:i], slice[i+1:]...)
			}
		}
	} else {
		for i, v := range slice {
			if condition(v) {
				slice[i] = slice[len(slice)-1]
				var zero T
				slice[len(slice)-1] = zero
				return slice[:len(slice)-1]
			}
		}
	}
	return slice
}
