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

		/**
		md :=metadata.New(map[string]string{"crmao":"crmaoclient"})
		newCtx :=metadata.NewOutgoingContext(ctx,md)
		链接的时候把ctx 传进去，那么就带上这个header头了
		**/

		//追加metadata
		//metadata.AppendToOutgoingContext(ctx,"crmao","go编程")
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
