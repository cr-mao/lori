package grpc

import (
	"context"
	"crypto/tls"
	"github.com/cr-mao/lori/metric"
	"net"
	"net/url"
	"time"

	"github.com/cr-mao/lori/internal/endpoint"
	"github.com/cr-mao/lori/internal/host"
	"github.com/cr-mao/lori/transport"
	//apimd "github.com/cr-mao/lori/api/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/cr-mao/lori/log"
)

var _ transport.Endpointer = (*Server)(nil)
var _ transport.Server = (*Server)(nil)

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server
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
	endpoint *url.URL          // url
	metric   metric.GrpcMetric //metric 接口，可以传可不传

	err error
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

	if srv.metric != nil {
		unaryInts = append(unaryInts, srv.metric.GrpcMetricInterceptors()...)
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

	//todo
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	//用户传的ServerOption
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

// Endpoint return a real address to registry endpoint.
// examples:
//
//	grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return s.err
	}
	s.baseCtx = ctx
	log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())
	//设置serving 状态
	s.health.Resume()
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(_ context.Context) error {
	//if s.adminClean != nil {
	//	s.adminClean()
	//}
	s.health.Shutdown()
	s.GracefulStop()
	log.Info("[gRPC] server stopping")
	return nil
}

func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.tlsConf != nil), addr)
	}
	return s.err
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

//func WithEnableMetrics(enableMetrics bool) ServerOption {
//	return func(o *Server) {
//		o.enableMetrics = enableMetrics
//	}
//}

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
