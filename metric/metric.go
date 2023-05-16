package metric

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type GrpcMetric interface {
	GrpcMetricInterceptors() []grpc.UnaryServerInterceptor //grpc 中间件
}

type GinMetric interface {
	MiddleWares() []gin.HandlerFunc
	SetMetricPath(path string)
	GetMetricPath() string
	Use(router gin.IRouter)
}
