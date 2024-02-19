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
	"github.com/cr-mao/lori/registry/consul"
	"github.com/hashicorp/consul/api"
	"testing"
	"time"

	"github.com/cr-mao/lori/transport/grpc"
)

// 基于服务发现的client
func TestDiscoverClient(t *testing.T) {
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

// 基于直连client
func TestDirectClient(t *testing.T) {
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithClientTimeout(time.Second*5),
		grpc.WithClientEndpoint("127.0.0.1:8081"),
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
