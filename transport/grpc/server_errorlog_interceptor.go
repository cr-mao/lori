package grpc

import (
	"context"
	"github.com/cr-mao/lori/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func unaryErrorLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	if gstatus, ok := status.FromError(err); ok {
		errLog := "grpc error:method:%s, code:%v,message:%v"
		log.Errorf(errLog, info.FullMethod, gstatus.Code(), err.Error())
	} else {
		errLog := "not grpc error:method:%s,message:%v"
		log.Errorf(errLog, info.FullMethod, err.Error())
	}
	return resp, err
}
