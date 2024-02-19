/**
User: cr-mao
Date: 2024/2/19 12:52
Email: crmao@qq.com
Desc: grpc_server_client.go
*/
package example

import (
	"context"
	"github.com/cr-mao/lori/example/proto"
	"github.com/cr-mao/lori/log"
	"testing"
	"time"

	"github.com/cr-mao/lori/transport/grpc"
)

func TestGrpcClient(t *testing.T) {
	baseCtx := context.Background()
	conn, err := grpc.DialInsecure(baseCtx,
		grpc.WithClientEndpoint("127.0.0.1:8081"),
		grpc.WithClientTimeout(time.Second*5))
	if err != nil {
		log.Fatalf("err:%v", err)
	}
	client := proto.NewGreeterClient(conn)
	resp, err := client.SayHello(context.Background(), &proto.HelloRequest{
		Name: "cr-mao",
	})
	if err != nil {
		log.Errorf("SayHello err:%v", err)
		return
	}
	log.Info(resp.Message)
}
