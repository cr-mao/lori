/**
User: cr-mao
Date: 2024/2/19 12:48
Email: crmao@qq.com
Desc: test_grpc_server.go
*/
package example

import (
	"context"
	"github.com/cr-mao/lori/example/proto"
	"github.com/cr-mao/lori/log"
	"testing"

	"github.com/cr-mao/lori/transport/grpc"
)

type HelloWorldServer struct {
	proto.UnsafeGreeterServer
}

func (s *HelloWorldServer) SayHello(ctx context.Context, r *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: "hello " + r.Name,
	}, nil
}

func registerServer(server *grpc.Server) {
	proto.RegisterGreeterServer(server, &HelloWorldServer{})
}

func TestGrpcServer(t *testing.T) {
	baseCtx := context.Background()
	grpcServer := grpc.NewServer(grpc.WithAddress("0.0.0.0:8081"))
	registerServer(grpcServer)
	err := grpcServer.Start(baseCtx)
	if err != nil {
		log.Errorf("err %v", err)
	}
}
