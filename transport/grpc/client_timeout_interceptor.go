package grpc

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// client超时控制中间件
func clientTimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		//没有超时则忽律
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		var cancel context.CancelFunc
		//如果没设置 超时 ，才追究超时
		if _, ok := ctx.Deadline(); !ok {
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
