package grpcmetric

import (
	"context"
	"github.com/cr-mao/lori/metric"
	"github.com/cr-mao/lori/metric/prometheus"
	"google.golang.org/grpc"
	"time"
)

type PromInstance struct {
	metricServerReqDur prometheus.HistogramVec
	serverName         string
}

var _ metric.GrpcMetric = (*PromInstance)(nil)

func NewMetricInstance(serverName string) metric.GrpcMetric {
	metricServerReqDur := prometheus.NewHistogramVec(&prometheus.HistogramVecOpts{
		Namespace: serverName + "_grpc_server",
		Subsystem: "requests",
		Name:      "grpc_server_duration_ms",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{30, 50, 100, 250, 500, 1000, 2000},
	})
	return &PromInstance{
		serverName:         serverName,
		metricServerReqDur: metricServerReqDur,
	}
}

func (p *PromInstance) GrpcMetricInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{p.serverUnaryPrometheusInterceptor}
}

//每个方法耗时的中间件
func (p *PromInstance) serverUnaryPrometheusInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	resp, err = handler(ctx, req)
	//记录了耗时
	p.metricServerReqDur.Observe(int64(time.Since(startTime)/time.Millisecond), info.FullMethod)
	return resp, err
}
