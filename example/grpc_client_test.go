/*
*
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
	"github.com/cr-mao/lori/registry/consul"
	"github.com/cr-mao/lori/trace"
	"github.com/hashicorp/consul/api"
	"testing"
	"time"

	"github.com/cr-mao/lori/transport/grpc"
)

// 基于服务发现的client,并集成trace基于jaeger
func TestDiscoverTraceClient(t *testing.T) {

	trace.InitAgent(trace.Options{
		Name:     "lori_example",
		Endpoint: "http://127.0.0.1:14268/api/traces",
		Sampler:  1.0,
		Batcher:  "jaeger",
	})
	conf := api.DefaultConfig()
	conf.Address = "127.0.0.1:8500"
	conf.Scheme = "http"
	cli, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithClientDiscovery(r),
		grpc.WithClientTimeout(time.Second*5),
		grpc.WithClientEndpoint("discovery:///lori-app"),
		grpc.WithClientEnableTracing(true),
		// 自定义trace中间件
		//grpc.WithClientUnaryInterceptor(grpc.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
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

	//  确保jager上报成功，不然会报 父span 找不到。
	time.Sleep(time.Second * 10)
}

// 基于直连client
func TestDirectClient(t *testing.T) {
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithClientTimeout(time.Second*5),
		grpc.WithClientEndpoint("127.0.0.1:8081"),
		grpc.WithClientEnableTracing(true),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
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
