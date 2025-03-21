package kmap

import (
	"testing"

	"github.com/mtgnorton/k/kalgo"

	"github.com/stretchr/testify/assert"
)

func TestRangeInOrder(t *testing.T) {
	t.Run("整数键升序遍历", func(t *testing.T) {
		m := map[int]string{3: "c", 1: "a", 2: "b"}
		result := make([]string, 0, len(m))
		keys := make([]int, 0, len(m))

		RangeInOrder(m, func(v string, k int) {
			result = append(result, v)
			keys = append(keys, k)
		})

		assert.Equal(t, []string{"a", "b", "c"}, result)
		assert.Equal(t, []int{1, 2, 3}, keys)
	})

	t.Run("整数键降序遍历", func(t *testing.T) {
		m := map[int]string{3: "c", 1: "a", 2: "b"}
		result := make([]string, 0, len(m))
		keys := make([]int, 0, len(m))

		RangeInOrder(m, func(v string, k int) {
			result = append(result, v)
			keys = append(keys, k)
		}, kalgo.SortDesc)

		assert.Equal(t, []string{"c", "b", "a"}, result)
		assert.Equal(t, []int{3, 2, 1}, keys)
	})

	t.Run("字符串键升序遍历", func(t *testing.T) {
		m := map[string]int{"c": 3, "a": 1, "b": 2}
		result := make([]int, 0, len(m))
		keys := make([]string, 0, len(m))

		RangeInOrder(m, func(v int, k string) {
			result = append(result, v)
			keys = append(keys, k)
		})

		assert.Equal(t, []int{1, 2, 3}, result)
		assert.Equal(t, []string{"a", "b", "c"}, keys)
	})

	t.Run("浮点数键升序遍历", func(t *testing.T) {
		m := map[float64]string{3.14: "pi", 2.71: "e", 1.41: "sqrt2"}
		result := make([]string, 0, len(m))
		keys := make([]float64, 0, len(m))

		RangeInOrder(m, func(v string, k float64) {
			result = append(result, v)
			keys = append(keys, k)
		})

		assert.Equal(t, []string{"sqrt2", "e", "pi"}, result)
		assert.Equal(t, []float64{1.41, 2.71, 3.14}, keys)
	})

	t.Run("空map", func(t *testing.T) {
		m := map[int]string{}
		called := false

		RangeInOrder(m, func(v string, k int) {
			called = true
		})

		assert.False(t, called, "回调函数不应被调用")
	})

	t.Run("单元素map", func(t *testing.T) {
		m := map[int]string{1: "a"}
		called := false

		RangeInOrder(m, func(v string, k int) {
			called = true
		})

		assert.True(t, called)
	})
}
