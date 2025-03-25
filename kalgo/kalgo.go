// Package kalgo 提供一些常用的算法实现
package kalgo

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

const (
	SortAsc  Sort = "asc"
	SortDesc Sort = "desc"
)

type Sort string

// QuickSort 快速排序算法实现
//
// 参数说明:
//   - arr: 待排序的数组
//   - l: 排序的起始位置
//   - r: 排序的结束位置
//   - sort: 可选的排序方式,默认为升序(SortAsc)
//
// 注意事项:
//   - 该函数会直接修改原数组
//   - 支持任意可比较类型
//   - 时间复杂度为O(nlogn),空间复杂度为O(logn)
//   - 当l >= r时会直接返回
//
// 示例:
//
//	arr := []int{3, 1, 4, 1, 5}
//	QuickSort(arr, 0, len(arr)-1) // 升序排序
//	QuickSort(arr, 0, len(arr)-1, SortDesc) // 降序排序
func QuickSort[T constraints.Ordered](arr []T, l, r int, sort ...Sort) {
	if l >= r {
		return
	}

	s := SortAsc
	if len(sort) > 0 {
		s = sort[0]
	}
	q := partition(arr, l, r, s)

	QuickSort(arr, l, q-1, s)
	QuickSort(arr, q+1, r, s)
}

func partition[T constraints.Ordered](arr []T, l, r int, sort Sort) int {
	var (
		i = l
		j = l
	)
	randomIndex := l + rand.Intn(r-l+1)
	arr[r], arr[randomIndex] = arr[randomIndex], arr[r]
	for ; j < r; j++ {
		if sort == SortAsc {
			if arr[j] <= arr[r] {
				arr[i], arr[j] = arr[j], arr[i]
				i++
			}
		} else {
			if arr[j] >= arr[r] {
				arr[i], arr[j] = arr[j], arr[i]
				i++
			}
		}
	}
	arr[i], arr[r] = arr[r], arr[i]
	return i
}
