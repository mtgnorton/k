package kmonitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/mtgnorton/k/kcollection"
	"github.com/stretchr/testify/assert"
)

func TestNewRollingResultCounter(t *testing.T) {
	// 测试默认创建
	counter := NewRollingResultCounter[int64]()
	assert.NotNil(t, counter)
	assert.NotNil(t, counter.successWindow)
	assert.NotNil(t, counter.failWindow)

	// 测试带选项创建
	opts := []kcollection.RollingWindowOption[int64, *kcollection.Bucket[int64]]{
		kcollection.WithSize[int64, *kcollection.Bucket[int64]](10),
		kcollection.WithInterval[int64, *kcollection.Bucket[int64]](time.Second),
	}
	counter = NewRollingResultCounter(opts...)
	assert.NotNil(t, counter)

	// 测试添加成功和失败记录
	counter.AddSuccess(100)
	counter.AddFail(50)

	// 测试Reduce方法
	var successCount, failCount int64
	var successTime, failTime int64
	counter.Reduce(
		func(sc int64, st int64) {
			successCount = sc
			successTime = st
		},
		func(fc int64, ft int64) {
			failCount = fc
			failTime = ft
		},
	)
	info := counter.Info("ms")
	fmt.Println(info)
	assert.Equal(t, int64(1), successCount)
	assert.Equal(t, int64(100), successTime)
	assert.Equal(t, int64(1), failCount)
	assert.Equal(t, int64(50), failTime)
}
