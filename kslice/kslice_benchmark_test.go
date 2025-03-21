package kslice

import "testing"

func BenchmarkRemoveElements(b *testing.B) {

	b.Run("OutOfOrder", func(b *testing.B) {
		original := make([]int, 400)
		for i := 0; i < 400; i++ {
			original[i] = i
		}
		for i := 0; i < b.N; i++ {
			slice := make([]int, 400)
			copy(slice, original)

			RemoveElements(slice, func(item int) bool {
				return item%2 == 0
			})
		}
	})

	b.Run("InOrder", func(b *testing.B) {
		original := make([]int, 400)
		for i := 0; i < 400; i++ {
			original[i] = i
		}
		for i := 0; i < b.N; i++ {
			slice := make([]int, 400)
			copy(slice, original)
			RemoveElements(slice, func(item int) bool {
				return item%2 == 0
			}, true)
		}
	})

	var simple = func(slice []int) []int {
		result := make([]int, 0, len(slice))
		for _, item := range slice {
			if item%2 != 0 {
				result = append(result, item)
			}
		}
		return result
	}

	b.Run("Simple", func(b *testing.B) {
		original := make([]int, 400)
		for i := 0; i < 400; i++ {
			original[i] = i
		}
		for i := 0; i < b.N; i++ {
			slice := make([]int, 400)
			copy(slice, original)
			simple(slice)
		}
	})

}

func BenchmarkRemoveElement(b *testing.B) {
	b.Run("OutOfOrder", func(b *testing.B) {
		original := make([]int, 400)
		for i := 0; i < 400; i++ {
			original[i] = i
		}
		for i := 0; i < b.N; i++ {
			slice := make([]int, 400)
			copy(slice, original)
			RemoveElement(slice, func(item int) bool {
				return item%2 == 0
			}, false)
		}
	})

	b.Run("InOrder", func(b *testing.B) {
		original := make([]int, 400)
		for i := 0; i < 400; i++ {
			original[i] = i
		}
		for i := 0; i < b.N; i++ {
			slice := make([]int, 400)
			copy(slice, original)
			RemoveElement(slice, func(item int) bool {
				return item%2 == 0
			}, true)
		}
	})

	var simple = func(slice []int) []int {
		result := make([]int, 0, len(slice))
		for _, item := range slice {
			if item%2 != 0 {
				result = append(result, item)
			}
		}
		return result
	}

	b.Run("Simple", func(b *testing.B) {
		original := make([]int, 400)
		for i := 0; i < 400; i++ {
			original[i] = i
		}
		for i := 0; i < b.N; i++ {
			slice := make([]int, 400)
			copy(slice, original)
			simple(slice)
		}
	})
}
