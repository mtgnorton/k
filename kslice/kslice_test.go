package kslice

import (
	"fmt"
	"math/rand"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRemoveElements(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		pred     func(int) bool
		inOrder  bool
		expected []int
	}{
		{
			name:     "删除偶数-无序",
			slice:    []int{1, 2, 3, 4, 5},
			pred:     func(i int) bool { return i%2 == 0 },
			inOrder:  false,
			expected: []int{1, 5, 3},
		},
		{
			name:     "删除偶数-有序",
			slice:    []int{1, 2, 3, 4, 5},
			pred:     func(i int) bool { return i%2 == 0 },
			inOrder:  true,
			expected: []int{1, 3, 5},
		},
		{
			name:     "删除奇数-无序",
			slice:    []int{1, 2, 3, 4, 5},
			pred:     func(i int) bool { return i%2 == 1 },
			inOrder:  false,
			expected: []int{2, 4},
		},
		{
			name:     "空切片",
			slice:    []int{},
			pred:     func(i int) bool { return true },
			inOrder:  true,
			expected: []int{},
		},
		{
			name:     "全部删除",
			slice:    []int{1, 2, 3},
			pred:     func(i int) bool { return true },
			inOrder:  true,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveElements(tt.slice, tt.pred, tt.inOrder)
			assert.Equal(t, len(tt.expected), len(result), "长度不匹配")
			if tt.inOrder {
				assert.Equal(t, tt.expected, result, "有序结果不匹配")
			} else {
				for _, exp := range tt.expected {
					assert.True(t, slices.Contains(result, exp), "结果 %v 中未找到期望元素 %v", result, exp)
				}
			}
		})
	}
}

func TestRemoveElement(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		pred      func(int) bool
		inOrder   bool
		expected  []int
		checkFunc func([]int, []int) bool
	}{
		{
			name:     "删除偶数-无序",
			slice:    []int{1, 2, 3, 4, 5},
			pred:     func(i int) bool { return i%2 == 0 },
			inOrder:  false,
			expected: []int{1, 3, 5, 4},
			checkFunc: func(result, expected []int) bool {
				return len(result) == 4 && !slices.Contains(result, 2)
			},
		},
		{
			name:     "删除偶数-有序",
			slice:    []int{1, 2, 3, 4, 5},
			pred:     func(i int) bool { return i%2 == 0 },
			inOrder:  true,
			expected: []int{1, 3, 4, 5},
			checkFunc: func(result, expected []int) bool {
				return len(result) == 4 && result[0] == 1 && result[1] == 3 && result[2] == 4 && result[3] == 5
			},
		},
		{
			name:     "空切片",
			slice:    []int{},
			pred:     func(i int) bool { return true },
			inOrder:  true,
			expected: []int{},
			checkFunc: func(result, expected []int) bool {
				return len(result) == 0
			},
		},
		{
			name:     "全部删除",
			slice:    []int{1, 2, 3},
			pred:     func(i int) bool { return true },
			inOrder:  true,
			expected: []int{},
			checkFunc: func(result, expected []int) bool {
				return len(result) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveElement(tt.slice, tt.pred, tt.inOrder)
			assert.True(t, tt.checkFunc(result, tt.expected), "期望 %v, 实际得到 %v", tt.expected, result)
		})
	}
}

func TestLoopConcAsync(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		// 准备测试数据
		data := []int{1, 2, 3, 4, 5}

		// 测试函数,将数字乘以2并返回
		fn := func(n int) (int, error) {
			if n == 3 {
				return 0, fmt.Errorf("处理数字3时发生错误")
			}
			return n * 2, nil
		}

		// 启动并发处理
		resultCh, cancel := LoopConcAsync(data, fn, 2)
		defer cancel()

		// 收集结果
		results := make(map[int]int)
		errors := make(map[int]error)

		for result := range resultCh {
			fmt.Printf("result: %+v\n", result)
			if result.Error != nil {
				errors[result.Key] = result.Error
			} else {
				results[result.Key] = result.Result
			}
		}

		// 验证结果
		assert.Equal(t, 4, len(results), "期望成功结果数量为4")
		assert.Equal(t, 1, len(errors), "期望错误数量为1")

		// 验证具体的结果值
		expectedResults := map[int]int{
			0: 2,  // 1*2
			1: 4,  // 2*2
			3: 8,  // 4*2
			4: 10, // 5*2
		}

		for k, v := range expectedResults {
			assert.Equal(t, v, results[k], "索引%d的结果期望为%d", k, v)
		}

		// 验证错误
		assert.Error(t, errors[2], "期望索引2处理失败")
		assert.Equal(t, "处理数字3时发生错误", errors[2].Error(), "错误消息不匹配")
	})

	t.Run("get first result then cancel", func(t *testing.T) {
		var data []int
		for i := 0; i < 100000; i++ {
			data = append(data, i)
		}
		fn := func(n int) (int, error) {
			// 随机休眠50ms-1s
			time.Sleep(time.Duration(50+rand.Intn(1000)) * time.Millisecond)
			return n, nil
		}

		resultCh, cancel := LoopConcAsync(data, fn, 100)

		result := <-resultCh

		assert.NoError(t, result.Error, "期望没有错误")
		fmt.Printf("goroutine: %d\n", runtime.NumGoroutine())
		cancel()
		time.Sleep(time.Second * 10)
	})
	t.Run("测试并发连接数量", func(t *testing.T) {
		// 创建一个大数据集
		var data []int
		for i := 0; i < 1000; i++ {
			data = append(data, i)
		}

		// 创建一个计数器来跟踪并发执行的数量
		var concurrentCount int32
		var maxConcurrentCount int32
		var mutex sync.Mutex

		fn := func(n int) (int, error) {
			// 增加当前并发计数
			atomic.AddInt32(&concurrentCount, 1)

			// 更新最大并发数
			mutex.Lock()
			if concurrentCount > maxConcurrentCount {
				maxConcurrentCount = concurrentCount
			}
			mutex.Unlock()

			// 模拟工作负载
			time.Sleep(time.Duration(10+rand.Intn(50)) * time.Millisecond)

			// 减少当前并发计数
			atomic.AddInt32(&concurrentCount, -1)

			return n * 2, nil
		}

		// 设置期望的并发数
		expectedConcurrency := 50

		// 执行并发处理
		resultCh, cancel := LoopConcAsync(data, fn, expectedConcurrency)
		defer cancel()
		// 消费所有结果
		count := 0
		for range resultCh {
			count++
		}

		// 验证结果
		assert.Equal(t, len(data), count, "应处理所有数据项")
		assert.LessOrEqual(t, maxConcurrentCount, int32(expectedConcurrency), "最大并发数不应超过设定值")
		assert.GreaterOrEqual(t, maxConcurrentCount, int32(expectedConcurrency/10*8), "最大并发数应接近设定值")
		t.Logf("最大并发数: %d, 设定并发数: %d", maxConcurrentCount, expectedConcurrency)
	})
}

func TestLoopConcAsyncFirstResult(t *testing.T) {
	t.Run("成功获取第一个结果", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		result, err := LoopConcAsyncFirstSuccess(data, func(n int) (int, error) {
			time.Sleep(time.Duration(100+rand.Intn(500)) * time.Millisecond)
			return n * 2, nil
		}, 3)

		assert.NoError(t, err, "期望没有错误")
		assert.NotEmpty(t, result, "期望有结果返回")
		assert.Equal(t, result.Result, result.Item*2, "结果值应该是输入值的2倍")
	})

	t.Run("所有处理都失败", func(t *testing.T) {
		data := []int{1, 2, 3}
		result, err := LoopConcAsyncFirstSuccess(data, func(n int) (int, error) {
			return 0, fmt.Errorf("处理 %d 失败", n)
		}, 2)

		assert.Error(t, err, "期望返回错误")
		assert.Empty(t, result.Result, "期望没有结果")
		assert.Contains(t, err.Error(), "处理", "错误信息应包含处理失败说明")
	})

	t.Run("部分成功部分失败", func(t *testing.T) {
		data := []int{1, 2, 3}
		result, err := LoopConcAsyncFirstSuccess(data, func(n int) (int, error) {
			if n == 2 {
				return n * 2, nil
			}
			return 0, fmt.Errorf("处理 %d 失败", n)
		}, 2)

		assert.NoError(t, err, "期望没有错误")
		assert.Equal(t, 4, result.Result, "期望结果为4")
		assert.Equal(t, 2, result.Item, "期望原始值为2")
	})
}
