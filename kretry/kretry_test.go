package kretry

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		result, err := Do(func(ctx context.Context) (string, error) {
			return "hello", nil
		})
		assert.NoError(t, err)
		assert.Equal(t, "hello", result)
	})
	t.Run("with context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		var attempt int
		result, err := Do(func(ctx context.Context) (string, error) {
			time.Sleep(40 * time.Millisecond)
			attempt++ //
			if attempt > 1 {
				return "hello", nil
			}
			return "", errors.Errorf("error: %d", attempt)
		}, WithContext(ctx))
		assert.Error(t, err)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected error to be context.DeadlineExceeded, got %v", err)
		}
		assert.Equal(t, "", result)
	})
	t.Run("with context cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()
		var attempt int
		result, err := Do(func(ctx context.Context) (string, error) {
			time.Sleep(40 * time.Millisecond)
			attempt++
			if attempt > 1 {
				return "hello", nil
			}
			return "", errors.Errorf("error: %d", attempt)
		}, WithContext(ctx))
		assert.Error(t, err)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected error to be context.Canceled, got %v", err)
		}
		assert.Equal(t, "", result)
	})

	t.Run("retry twice then success", func(t *testing.T) {
		var attempt int
		result, err := Do(func(ctx context.Context) (string, error) {
			attempt++
			if attempt > 2 {
				return "success", nil
			}
			return "", errors.Errorf("error attempt: %d", attempt)
		})
		assert.NoError(t, err)
		assert.Equal(t, "success", result)
		assert.Equal(t, 3, attempt)
	})

	t.Run("skip retry when error match", func(t *testing.T) {
		var attempt int
		result, err := Do(func(ctx context.Context) (string, error) {
			attempt++
			return "", errors.New("skip retry")
		}, WithErrHandler(func(err error) bool {
			return err.Error() == "skip retry"
		}))
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Equal(t, 1, attempt) // 只尝试一次就停止
	})

	t.Run("retry three times with custom delay", func(t *testing.T) {
		var attempt int
		// 第一次重试完成 attempt=1
		// 第二次重试完成 attempt=2
		// 第三次重试完成 attempt=3
		result, err := Do(func(ctx context.Context) (string, error) {
			attempt++
			if attempt < 3 {
				return "", errors.New("error")
			}
			return "success", nil
		}, WithCustomDelay([]time.Duration{50 * time.Millisecond, 100 * time.Millisecond, 150 * time.Millisecond}))
		assert.NoError(t, err)
		assert.Equal(t, "success", result)
		assert.Equal(t, 3, attempt)
	})
}
