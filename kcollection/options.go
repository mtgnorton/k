package kcollection

import (
	"time"

	"github.com/mtgnorton/k/kmath"
)

type RollingWindowOptions[T kmath.Number, B BucketInterface[T]] struct {
	size          int           // 窗口大小(桶的数量)
	interval      time.Duration // 每个桶的时间间隔
	ignoreCurrent bool          // 是否忽略当前桶
}

type RollingWindowOption[T kmath.Number, B BucketInterface[T]] func(opts *RollingWindowOptions[T, B])

func NewRollingWindowOptions[T kmath.Number, B BucketInterface[T]]() *RollingWindowOptions[T, B] {
	return &RollingWindowOptions[T, B]{
		size:     10,
		interval: time.Minute,
	}
}

func WithSize[T kmath.Number, B BucketInterface[T]](size int) RollingWindowOption[T, B] {
	return func(opts *RollingWindowOptions[T, B]) {
		opts.size = size
	}
}

func WithInterval[T kmath.Number, B BucketInterface[T]](interval time.Duration) RollingWindowOption[T, B] {
	return func(opts *RollingWindowOptions[T, B]) {
		opts.interval = interval
	}
}

func WithIgnoreCurrent[T kmath.Number, B BucketInterface[T]](ignore bool) RollingWindowOption[T, B] {
	return func(opts *RollingWindowOptions[T, B]) {
		opts.ignoreCurrent = ignore
	}
}
