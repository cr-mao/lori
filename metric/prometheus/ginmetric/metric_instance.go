package ginmetric

import (
	"github.com/cr-mao/lori/metric"
	"github.com/gin-gonic/gin"
)

const defaultMetricPath = "/metric"

type Instance struct {
}

func NewMetricInstance() metric.GinMetric {
	return &Instance{}
}

func (s *Instance) MiddleWares() []gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (s *Instance) SetMetricPath(path string) {
	//TODO implement me
	panic("implement me")
}

func (s *Instance) GetMetricPath() string {
	//TODO implement me
	panic("implement me")
}

func (s *Instance) Use(router gin.IRouter) {
	//TODO implement me
	panic("implement me")
}
