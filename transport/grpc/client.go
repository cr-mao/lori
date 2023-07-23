package grpc

import (
	"context"
	"crypto/tls"
	"github.com/cr-mao/lori/metric"
	"github.com/cr-mao/lori/registry"
	"github.com/cr-mao/lori/transport/grpc/resolver/discovery"
	"google.golang.org/grpc/credentials"

	grpcinsecure "google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
	"time"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	endpoint string
	timeout  time.Duration
	// discovery接口
	discovery  registry.Discovery
	unaryInts  []grpc.UnaryClientInterceptor
	streamInts []grpc.StreamClientInterceptor
	metric     metric.GrpcClientMetric //metric 接口，可以传可不传
	rpcOpts    []grpc.DialOption
	tlsConf    *tls.Config

	balancerName  string
	enableTracing bool
	enableMetrics bool
}

func WithMetric(metric metric.GrpcClientMetric) ClientOption {
	return func(o *clientOptions) {
		o.metric = metric
	}
}

func WithEnableTracing(enable bool) ClientOption {
	return func(o *clientOptions) {
		o.enableTracing = enable
	}
}

// 设置地址
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// 设置超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// 设置服务发现
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// 设置拦截器
func WithClientUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.unaryInts = in
	}
}

// 设置stream拦截器
func WithClientStreamInterceptor(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.streamInts = in
	}
}

// 设置grpc的dial选项
func WithClientOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.rpcOpts = opts
	}
}

// 设置负载均衡器
func WithBalancerName(name string) ClientOption {
	return func(o *clientOptions) {
		o.balancerName = name
	}
}

func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

/*

 */

func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout:       2000 * time.Millisecond,
		balancerName:  "round_robin",
		enableTracing: true,
	}

	for _, o := range opts {
		o(&options)
	}

	// 超时中间件
	ints := []grpc.UnaryClientInterceptor{
		clientTimeoutInterceptor(options.timeout),
	}
	//if options.enableTracing {
	//	ints = append(ints, otelgrpc.UnaryClientInterceptor())
	//}

	streamInts := []grpc.StreamClientInterceptor{}

	if len(options.unaryInts) > 0 {
		ints = append(ints, options.unaryInts...)
	}

	if options.metric != nil {
		ints = append(ints, options.metric.GrpcClientMetricInterceptors()...)
	}

	if len(options.streamInts) > 0 {
		streamInts = append(streamInts, options.streamInts...)
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "` + options.balancerName + `"}`),
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithChainStreamInterceptor(streamInts...),
	}

	// 服务发现的选项
	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(
			discovery.NewBuilder(
				options.discovery,
				discovery.WithInsecure(insecure),
			),
		))
	}
	if options.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
	}
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}
	if len(options.rpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.rpcOpts...)
	}
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
