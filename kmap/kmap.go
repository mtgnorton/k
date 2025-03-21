package kmap

import (
	"sync"

	"github.com/mtgnorton/k/kalgo"
	"golang.org/x/exp/constraints"
)

// Copy 浅拷贝一个map并返回副本,
//
// 参数说明:
//   - src: 源map,需要被复制的map
//
// 返回值说明:
//   - map[K]V: 返回一个新的map,包含src中所有的键值对
//
// 注意事项:
//   - 仅支持可比较类型的key
//   - 对于值为引用类型的情况,只会复制引用而不是深度复制
//
// 示例:
//
//	srcMap := map[string]int{"a": 1, "b": 2}
//	dstMap := Copy(srcMap) // dstMap: map[string]int{"a": 1, "b": 2}
func Copy[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// RangeInOrder 按照key的顺序遍历map
//
// 参数说明:
//   - m: 要遍历的map
//   - fn: 遍历时的回调函数,接收value和key作为参数
//   - sort: 可选的排序方式,默认为升序可选值:kalgo.SortAsc,kalgo.SortDesc
//
// 返回值说明:
//   - 无
//
// 注意事项:
//   - key必须是可排序类型
//   - 当map长度小于等于1时会直接返回
//   - 遍历顺序由sort参数决定,默认升序
//
// 示例:
//
//	m := map[int]string{1: "a", 2: "b", 3: "c"}
//	RangeInOrder(m, func(v string, k int) {
//	    fmt.Println(k, v) // 按key升序打印: 1 a, 2 b, 3 c
//	})
func RangeInOrder[K constraints.Ordered, V any](m map[K]V, fn func(v V, k K), sort ...kalgo.Sort) {
	if len(m) <= 1 {
		return
	}
	// 获取所有key
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	kalgo.QuickSort(keys, 0, len(keys)-1, sort...)

	// 按顺序遍历
	for _, k := range keys {
		fn(m[k], k)
	}
}

// LoopConc 并发遍历map中的每个键值对
//
// 参数说明:
//   - m: 需要遍历的map
//   - fn: 处理每个键值对的函数,接收key和value作为参数
//   - concurrency: 可选参数,控制并发数,默认为1
//
// 返回值说明:
//   - 无返回值
//
// 注意事项:
//   - 该函数会阻塞直到所有并发任务完成
//   - 如果concurrency参数小于等于0,并发数会被设置为1
//   - 每个键值对都会在一个独立的goroutine中处理
//   - 使用sync.WaitGroup和channel来控制并发数
//
// 示例:
//
//	m := map[string]int{"a": 1, "b": 2}
//	LoopConc(m, func(k string, v int) {
//	    fmt.Println(k, v)
//	})
func LoopConc[K comparable, V any](m map[K]V, fn func(key K, value V), concurrency ...int) {
	wg := sync.WaitGroup{}
	if len(concurrency) == 0 {
		concurrency = []int{1}
	}
	ch := make(chan struct{}, concurrency[0])
	for key, value := range m {
		wg.Add(1)
		ch <- struct{}{}
		go func(key K, value V) {
			defer func() {
				wg.Done()
				<-ch
			}()
			fn(key, value)
		}(key, value)
	}
	wg.Wait()
}

// ChunkConc 将map分块并发处理
//
// 参数说明:
//   - m: 需要处理的map
//   - size: 每个分块的大小
//   - fn: 处理每个分块的函数,接收分块作为参数
//   - concurrent: 可选参数,控制并发数,默认为1
//
// 返回值说明:
//   - 无返回值
//
// 注意事项:
//   - 该函数会阻塞直到所有并发任务完成
//   - 如果size参数小于等于0,函数直接返回
//   - 如果map为空,函数直接返回
//   - 如果concurrent参数小于等于0,并发数会被设置为1
//   - 使用sync.WaitGroup和channel来控制并发数
//   - 每个分块都会在一个独立的goroutine中处理
//
// 示例:
//
//	m := map[string]int{"a": 1, "b": 2, "c": 3}
//	ChunkConc(m, 2, func(chunk map[string]int) {
//	    fmt.Println(chunk)
//	})

func ChunkConc[K comparable, V any](m map[K]V, size int, fn func(chunk map[K]V), concurrent ...int) {
	if len(concurrent) == 0 {
		concurrent = []int{1}
	}
	if size <= 0 {
		return
	}
	if len(m) == 0 {
		return
	}
	length := len(m)

	wg := sync.WaitGroup{}
	ch := make(chan struct{}, concurrent[0])

	chunk := make(map[K]V)
	i := 0
	for key, value := range m {
		chunk[key] = value
		i++
		if i%size == 0 || i == length {
			wg.Add(1)
			ch <- struct{}{}
			go func(chunk map[K]V) {
				defer func() {
					wg.Done()
					<-ch
				}()
				fn(chunk)
			}(chunk)
			chunk = make(map[K]V)
		}
	}
	wg.Wait()
}
