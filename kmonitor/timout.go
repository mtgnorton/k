package kmonitor

import (
	"sync"
	"time"

	"github.com/mtgnorton/k/kunique"
)

// defaultTimeoutController 默认的超时检测器实例
var defaultTimeoutController = &TimeoutController{
	callIDs: make(map[int64]struct{}),
}

// TimeoutController 超时检测器
type TimeoutController struct {
	callIDs      map[int64]struct{} // 记录活跃的调用ID
	sync.RWMutex                    // 使用读写锁提升性能
}

// NewTimeoutController 创建一个新的超时检测器
func NewTimeoutController() *TimeoutController {
	return &TimeoutController{
		callIDs: make(map[int64]struct{}),
	}
}

// Do 执行一个带超时检测的任务
//
// 参数说明:
//   - duration: 超时时间
//   - timeoutHandler: 超时处理函数
//
// 返回值说明:
//   - end: 用于提前结束任务的函数
//
// 注意事项:
//   - 使用互斥锁保证并发安全
//   - 超时后会自动清理资源
//   - 调用end函数会停止定时器并清理资源
//   - 每个任务都有唯一的callID标识
//
// 示例:
//
//	end := monitor.Do(5*time.Second, func() {
//	    fmt.Println("timeout")
//	})
//	defer end()
func (t *TimeoutController) Do(duration time.Duration, timeoutHandler func()) (end func()) {
	callID := kunique.GenerateUniqueID()

	t.Lock()
	t.callIDs[callID] = struct{}{}
	t.Unlock()

	timer := time.AfterFunc(duration, func() {
		t.Lock()
		defer t.Unlock()
		if _, ok := t.callIDs[callID]; ok {
			timeoutHandler()
			delete(t.callIDs, callID)
		}
	})

	return func() {
		timer.Stop() // 停止定时器
		t.Lock()
		delete(t.callIDs, callID)
		t.Unlock()
	}
}

// MonitorTimeout 监控超时,参见 TimeoutController.Do
func MonitorTimeout(duration time.Duration, timeoutHandler func()) (end func()) {
	return defaultTimeoutController.Do(duration, timeoutHandler)
}
