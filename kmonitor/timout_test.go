package kmonitor

import (
	"testing"
	"time"
)

func TestMonitorTimeout(t *testing.T) {
	// 测试正常情况下不会触发超时
	triggered := false
	end := MonitorTimeout(100*time.Millisecond, func() {
		triggered = true
	})
	time.Sleep(50 * time.Millisecond)
	end()
	if triggered {
		t.Error("不应该触发超时")
	}

	// 测试超时会被正确触发
	triggered = false
	end = MonitorTimeout(50*time.Millisecond, func() {
		triggered = true
	})
	time.Sleep(100 * time.Millisecond)
	if !triggered {
		t.Error("应该触发超时")
	}
	end()

	// 测试多个并发的超时检测
	const concurrency = 10
	done := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			triggered := false
			end := MonitorTimeout(50*time.Millisecond, func() {
				triggered = true
			})
			time.Sleep(100 * time.Millisecond)
			if !triggered {
				t.Error("并发测试中应该触发超时")
			}
			end()
			done <- true
		}()
	}

	// 等待所有并发测试完成
	for i := 0; i < concurrency; i++ {
		<-done
	}
}

func TestTimeoutController(t *testing.T) {
	controller := NewTimeoutController()

	// 测试提前结束不会触发超时
	triggered := false
	end := controller.Do(100*time.Millisecond, func() {
		triggered = true
	})
	time.Sleep(50 * time.Millisecond)
	end()
	if triggered {
		t.Error("提前结束不应该触发超时")
	}

	// 测试超时处理器会被正确调用
	triggered = false
	end = controller.Do(50*time.Millisecond, func() {
		triggered = true
	})
	time.Sleep(100 * time.Millisecond)
	if !triggered {
		t.Error("应该触发超时处理器")
	}
	end()
}
