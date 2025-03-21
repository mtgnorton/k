package kalgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuickSort(t *testing.T) {
	t.Run("整数升序排序", func(t *testing.T) {
		arr := []int{3, 1, 4, 1, 5, 9, 2, 6}
		QuickSort(arr, 0, len(arr)-1)
		assert.Equal(t, []int{1, 1, 2, 3, 4, 5, 6, 9}, arr)
	})

	t.Run("整数降序排序", func(t *testing.T) {
		arr := []int{3, 1, 4, 1, 5, 9, 2, 6}
		QuickSort(arr, 0, len(arr)-1, SortDesc)
		assert.Equal(t, []int{9, 6, 5, 4, 3, 2, 1, 1}, arr)
	})

	t.Run("浮点数升序排序", func(t *testing.T) {
		arr := []float64{3.14, 1.41, 2.71, 0.58}
		QuickSort(arr, 0, len(arr)-1)
		assert.Equal(t, []float64{0.58, 1.41, 2.71, 3.14}, arr)
	})

	t.Run("字符串升序排序", func(t *testing.T) {
		arr := []string{"banana", "apple", "orange", "grape"}
		QuickSort(arr, 0, len(arr)-1)
		assert.Equal(t, []string{"apple", "banana", "grape", "orange"}, arr)
	})

	t.Run("空数组", func(t *testing.T) {
		var arr []int
		QuickSort(arr, 0, -1)
		assert.Empty(t, arr)
	})

	t.Run("单元素数组", func(t *testing.T) {
		arr := []int{1}
		QuickSort(arr, 0, 0)
		assert.Equal(t, []int{1}, arr)
	})
}
