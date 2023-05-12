package grpc

import (
	"context"
	"google.golang.org/grpc"
	"runtime/debug"

	"github.com/cr-mao/lori/log"
)

// 防止panic crash 中间件

func streamCrashInterceptor(svr interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v\n \n %s", r, debug.Stack())
	})

	return handler(svr, stream)
}

func unaryCrashInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v\n \n %s", r, debug.Stack())
	})
	return handler(ctx, req)
}

func handleCrash(hanlder func(interface{})) {
	if r := recover(); r != nil {
		hanlder(r)
	}
}
