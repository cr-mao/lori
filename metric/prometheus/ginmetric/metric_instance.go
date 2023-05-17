package ginmetric

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cr-mao/lori/metric"
	"github.com/cr-mao/lori/metric/prometheus"
)

const defaultMetricPath = "/metric"

type Instance struct {
	metricPath   string
	serverName   string
	metricReqDur prometheus.HistogramVec //每个方法请求耗时
}

var _ metric.GinMetric = (*Instance)(nil)

func NewMetricInstance(serverName string) metric.GinMetric {
	metricReqDur := prometheus.NewHistogramVec(&prometheus.HistogramVecOpts{
		Namespace: serverName + "_gin_server",
		Subsystem: "requests",
		Name:      "gin_http_server_duration_ms",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{30, 50, 100, 250, 500, 1000, 2000},
	})
	return &Instance{
		metricPath:   defaultMetricPath,
		serverName:   serverName,
		metricReqDur: metricReqDur,
	}
}

func (s *Instance) SetMetricPath(path string) {
	s.metricPath = path
}

func (s *Instance) Use(r gin.IRouter) {
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		//记录了耗时
		s.metricReqDur.Observe(int64(time.Since(start)/time.Millisecond), c.Request.URL.Path+":"+c.Request.Method)
	})
	r.GET(s.metricPath, func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})
}
