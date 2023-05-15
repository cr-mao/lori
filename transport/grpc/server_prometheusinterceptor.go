package grpc

import (
	"context"
	"time"

	"github.com/cr-mao/lori/metric"
	"google.golang.org/grpc"
)

/*
基本指标。 1. 每个请求的耗时(histogram)
*/

//  名称暂时写死, 可以不用这个，自己外面传中间件即可。
var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "grpc_server",
		Subsystem: "requests",
		Name:      "grpc_server_duration_ms",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{30, 50, 100, 250, 500, 1000, 2000},
	})
)

//每个方法耗时的中间件
func serverUnaryPrometheusInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	resp, err = handler(ctx, req)
	//记录了耗时
	metricServerReqDur.Observe(int64(time.Since(startTime)/time.Millisecond), info.FullMethod)
	return resp, err
}
