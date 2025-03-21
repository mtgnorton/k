package kbase

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsDeadlineError 判断错误是否为超时错误
//
// 参数说明:
//   - err: 需要判断的错误
//
// 返回值说明:
//   - bool: 如果是超时错误返回true,否则返回false
//
// 注意事项:
//   - 支持判断context.DeadlineExceeded错误
//   - 支持判断grpc的DeadlineExceeded错误码
//
// 示例:
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	err := doSomething(ctx)
//	if IsDeadlineError(err) {
//	    // 处理超时错误
//	}
func IsDeadlineError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var (
		statusErr *status.Status
		ok        bool
	)
	if statusErr, ok = status.FromError(err); !ok {
		return false
	}
	return statusErr.Code() == codes.DeadlineExceeded
}
