package http

import (
	"crypto/tls"
	"github.com/cr-mao/lori/metric"
	"time"
)

type ServerOption func(*Server)

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}
func WithEnableProfiling(profiling bool) ServerOption {
	return func(s *Server) {
		s.enableProfiling = profiling
	}
}

func WithMode(mode string) ServerOption {
	return func(s *Server) {
		s.mode = mode
	}
}

func WithServiceName(srvName string) ServerOption {
	return func(s *Server) {
		s.serviceName = srvName
	}
}

func WithMiddlewares(middlewares []string) ServerOption {
	return func(s *Server) {
		s.middlewares = middlewares
	}
}

func WithMetric(metric metric.GinMetric) ServerOption {
	return func(s *Server) {
		s.metric = metric
	}
}
