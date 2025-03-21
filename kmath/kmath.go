// kmath 提供了一些对数字类型的一些常用操作
//
// 主要功能:
//   - Max: 返回两个可比较类型值中的较大值
//   - Min: 返回两个可比较类型值中的较小值
//   - Round: 四舍五入保留n位小数
//   - Floor: 向下取整
//   - Ceil: 向上取整
//   - Abs: 返回一个数的绝对值
//   - Pow: 返回一个数的n次方
//   - Sqrt: 返回一个数的平方根
//   - RandInt: 返回一个随机整数
//   - RandFloat: 返回一个随机浮点数
package kmath

import (
	"cmp"
	"math"
	"math/rand"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

// Max 返回两个可比较类型值中的较大值
//
// 参数说明:
//   - a: 第一个比较值
//   - b: 第二个比较值
//
// 返回值:
//   - 两个值中的较大者
//
// 示例:
//
//	max := Max(10, 20)
//	// max = 20
//
//	max := Max(3.14, 2.71)
//	// max = 3.14
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min 返回两个可比较类型值中的较小值
//
// 参数说明:
//   - a: 第一个比较值
//   - b: 第二个比较值
//
// 返回值:
//   - 两个值中的较小者
//
// 示例:
//
//	min := Min(10, 20)
//	// min = 10
//
//	min := Min(3.14, 2.71)
//	// min = 2.71
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Round 四舍五入保留n位小数
//
// 参数说明:
//   - f: 需要四舍五入的浮点数
//   - n: 保留的小数位数
//
// 返回值:
//   - 四舍五入后的浮点数
//
// 示例:
//
//	rounded := Round(3.14159, 2)
//	// rounded = 3.14
//
//	rounded := Round(2.718, 1)
//	// rounded = 2.7
func Round[T ~float32 | ~float64](f T, n int) T {
	pow := math.Pow(10, float64(n))
	return T(math.Round(float64(f)*pow) / pow)
}

// Floor 向下取整
//
// 参数说明:
//   - f: 需要向下取整的浮点数
//
// 返回值:
//   - 向下取整后的浮点数
//
// 示例:
//
//	floored := Floor(3.75)
//	// floored = 3.0
//
//	floored := Floor(-2.3)
//	// floored = -3.0
func Floor[T ~float32 | ~float64](f T) T {
	return T(math.Floor(float64(f)))
}

// Ceil 向上取整
//
// 参数说明:
//   - f: 需要向上取整的浮点数
//   - n: 保留的小数位数
//
// 返回值:
//   - 向上取整后的浮点数
//
// 示例:
//
//	ceiled := Ceil(3.14, 1)
//	// ceiled = 3.2
//
//	ceiled := Ceil(2.01, 0)
//	// ceiled = 3.0
func Ceil[T ~float32 | ~float64](f T, n int) T {
	pow := math.Pow(10, float64(n))
	return T(math.Ceil(float64(f)*pow) / pow)
}

// Abs 返回一个数的绝对值
//
// 参数说明:
//   - f: 需要取绝对值的数
//
// 返回值:
//   - 绝对值结果
//
// 示例:
//
//	abs := Abs(-10)
//	// abs = 10
//
//	abs := Abs(3.14)
//	// abs = 3.14
func Abs[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](f T) T {
	return T(math.Abs(float64(f)))
}

// Pow 返回一个数的n次方
//
// 参数说明:
//   - f: 底数
//   - n: 指数
//
// 返回值:
//   - f的n次方结果
//
// 示例:
//
//	pow := Pow(2, 3)
//	// pow = 8
//
//	pow := Pow(10, 2)
//	// pow = 100
func Pow[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](f T, n int) T {
	return T(math.Pow(float64(f), float64(n)))
}

// Sqrt 返回一个数的平方根
//
// 参数说明:
//   - f: 需要计算平方根的数
//
// 返回值:
//   - f的平方根
//
// 示例:
//
//	sqrt := Sqrt(9)
//	// sqrt = 3
//
//	sqrt := Sqrt(2)
//	// sqrt = 1.4142135623730951
func Sqrt[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](f T) T {
	return T(math.Sqrt(float64(f)))
}

// RandInt 返回一个随机整数
//
// 参数说明:
//   - min: 随机数的最小值（包含）
//   - max: 随机数的最大值（包含）
//
// 返回值:
//   - 介于min和max之间的随机整数
//
// 示例:
//
//	rand := RandInt(1, 10)
//	// rand 是1到10之间的随机整数
//
//	rand := RandInt(-5, 5)
//	// rand 是-5到5之间的随机整数
func RandInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](min, max T) T {
	return T(rand.Intn(int(max-min+1)) + int(min))
}

// RandFloat 返回一个随机浮点数
//
// 参数说明:
//   - min: 随机数的最小值（包含）
//   - max: 随机数的最大值（不包含）
//
// 返回值:
//   - 介于min和max之间的随机浮点数
//
// 示例:
//
//	rand := RandFloat(0.0, 1.0)
//	// rand 是0.0到1.0之间的随机浮点数
//
//	rand := RandFloat(1.5, 3.5)
//	// rand 是1.5到3.5之间的随机浮点数
func RandFloat[T ~float32 | ~float64](min, max T) T {
	return T(rand.Float64()*float64(max-min) + float64(min))
}
