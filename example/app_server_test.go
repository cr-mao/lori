/**
User: cr-mao
Date: 2024/2/19 12:48
Email: crmao@qq.com
Desc: test_grpc_server.go
*/
package example

import (
	"context"
	"github.com/cr-mao/lori"
	"github.com/cr-mao/lori/example/proto"
	"github.com/cr-mao/lori/internal/netlib"
	"github.com/cr-mao/lori/log"
	"github.com/cr-mao/lori/pprof_server"
	"github.com/cr-mao/lori/registry/consul"
	"github.com/hashicorp/consul/api"
	"strconv"
	"testing"
	"time"

	"github.com/cr-mao/lori/transport/grpc"
)

type HelloWorldServer struct {
	proto.UnsafeGreeterServer
}

func (s *HelloWorldServer) SayHello(ctx context.Context, r *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Info("in")
	return &proto.HelloResponse{
		Message: "hello " + r.Name,
	}, nil
}

func registerServer(server *grpc.Server) {
	proto.RegisterGreeterServer(server, &HelloWorldServer{})
}

func TestAppServer(t *testing.T) {
	c := api.DefaultConfig()
	c.Address = "127.0.0.1:8500"
	c.Scheme = "http"
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	grpcServer := grpc.NewServer(grpc.WithAddress("0.0.0.0:8081"))
	pprofPort, _ := netlib.AssignRandPort()
	pprofServer := pprof_server.NewPProf("0.0.0.0:" + strconv.Itoa(pprofPort))
	registerServer(grpcServer)
	app := lori.New(lori.WithName("lori-app"),
		lori.WithServer(grpcServer, pprofServer),
		lori.WithRegistrar(r),
		lori.WithRegistrarTimeout(time.Second*5),
	)
	err = app.Run()
	if err != nil {
		log.Errorf("err %v", err)
	}
}

func TestAppServer2(t *testing.T) {
	c := api.DefaultConfig()
	c.Address = "127.0.0.1:8500"
	c.Scheme = "http"
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	grpcServer := grpc.NewServer(grpc.WithAddress("0.0.0.0:8082"))
	pprofPort, _ := netlib.AssignRandPort()
	pprofServer := pprof_server.NewPProf("0.0.0.0:" + strconv.Itoa(pprofPort))
	registerServer(grpcServer)
	app := lori.New(lori.WithName("lori-app"),
		lori.WithServer(grpcServer, pprofServer),
		lori.WithRegistrar(r),
		lori.WithRegistrarTimeout(time.Second*5),
	)
	err = app.Run()
	if err != nil {
		log.Errorf("err %v", err)
	}
}
