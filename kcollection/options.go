package kcollection

import (
	"time"

	"github.com/mtgnorton/k/kmath"
)

type RollingWindowOptions[T kmath.Number, B BucketInterface[T]] struct {
	Size          int           // 窗口大小(桶的数量)
	Interval      time.Duration // 每个桶的时间间隔
	IgnoreCurrent bool          // 是否忽略当前桶
}

type RollingWindowOption[T kmath.Number, B BucketInterface[T]] func(opts *RollingWindowOptions[T, B])

func NewRollingWindowOptions[T kmath.Number, B BucketInterface[T]]() *RollingWindowOptions[T, B] {
	return &RollingWindowOptions[T, B]{
		Size:     10,
		Interval: time.Minute,
	}
}

func WithSize[T kmath.Number, B BucketInterface[T]](size int) RollingWindowOption[T, B] {
	return func(opts *RollingWindowOptions[T, B]) {
		opts.Size = size
	}
}

func WithInterval[T kmath.Number, B BucketInterface[T]](interval time.Duration) RollingWindowOption[T, B] {
	return func(opts *RollingWindowOptions[T, B]) {
		opts.Interval = interval
	}
}

func WithIgnoreCurrent[T kmath.Number, B BucketInterface[T]](ignore bool) RollingWindowOption[T, B] {
	return func(opts *RollingWindowOptions[T, B]) {
		opts.IgnoreCurrent = ignore
	}
}
