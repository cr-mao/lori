package grpc

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"net/url"
	"time"

	//apimd "github.com/cr-mao/lori/api/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

// Server is a gRPC server wrapper.
type Server struct {
	Server     *grpc.Server
	tlsConf    *tls.Config // TLS configuration ,https
	baseCtx    context.Context
	network    string                         // tcp
	address    string                         //地址
	unaryInts  []grpc.UnaryServerInterceptor  //一次元拦截器
	streamInts []grpc.StreamServerInterceptor //流式拦截器
	grpcOpts   []grpc.ServerOption            //
	lis        net.Listener
	timeout    time.Duration
	health     *health.Server // 健康检测server
	//metadata      *apimd.Server
	endpoint      *url.URL // url
	enableMetrics bool     //是否开启链路追踪
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx: context.Background(),
		network: "tcp",
		address: ":0",
		timeout: 2 * time.Second,
		health:  health.NewServer(),
	}
	for _, o := range opts {
		o(srv)
	}
	unaryInts := []grpc.UnaryServerInterceptor{
		unaryCrashInterceptor,        //防止panic crash 中间件
		srv.unaryServerInterceptor(), //metadata 方便获取， 请求超时控制中间件
	}
	streamInts := []grpc.StreamServerInterceptor{
		streamCrashInterceptor,
		srv.streamServerInterceptor(),
	}
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}
	if len(srv.streamInts) > 0 {
		streamInts = append(streamInts, srv.streamInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	//srv.metadata = apimd.NewServer(srv.Server)
	// internal register ,健康检测， 框架自带的server
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	//apimd.RegisterMetadataServer(srv.Server, srv.metadata)

	// Register reflection service on gRPC server.
	reflection.Register(srv.Server)
	return srv
}

// ServerOption is gRPC server option.
type ServerOption func(o *Server)

func WithNetwork(network string) ServerOption {
	return func(o *Server) {
		o.network = network
	}
}

func WithAddress(address string) ServerOption {
	return func(o *Server) {
		o.address = address
	}
}

func WithEnableMetrics(enableMetrics bool) ServerOption {
	return func(o *Server) {
		o.enableMetrics = enableMetrics
	}
}

// Listener with server lis
func WithListener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// WithUnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInts = in
	}
}

// WithStreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the server.
func WithStreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInts = in
	}
}

// WithGrpcOpts with grpc options.
func WithGrpcOpts(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}
