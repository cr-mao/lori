package http

import (
	"context"
	"crypto/tls"
	"errors"

	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/cr-mao/lori/internal/endpoint"
	"github.com/cr-mao/lori/internal/host"
	"github.com/cr-mao/lori/log"
	"github.com/cr-mao/lori/transport"
	mids "github.com/cr-mao/lori/transport/http/middlewares"
	"github.com/cr-mao/lori/transport/http/pprof"
	"github.com/gin-gonic/gin"
)

var _ transport.Endpointer = (*Server)(nil)
var _ transport.Server = (*Server)(nil)

// wrapper for gin.Engine
type Server struct {
	*gin.Engine
	server      *http.Server
	lis         net.Listener
	endpoint    *url.URL
	serviceName string //服务名
	network     string
	address     string

	//开发模式， 默认值 debug
	mode string

	//是否开启pprof接口， 默认开启， 如果开启会自动添加 /debug/pprof 接口
	enableProfiling bool

	//是否开启metrics接口， 默认开启， 如果开启会自动添加 /metrics 接口, prometheus
	enableMetrics bool

	// 指标path  默认/metrics
	metricsPath string

	//中间件
	middlewares []string //传字符串进来， 顺序需要自己定义好

	//请求超时
	timeout time.Duration

	err error

	tlsConf *tls.Config
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network:         "tcp",
		address:         ":0",
		mode:            "debug",
		enableProfiling: true,
		Engine:          gin.New(), //纯的，没有logger，和default 。
		serviceName:     "lori-gin-http",
		timeout:         time.Second * 5, //默认5秒
		metricsPath:     "/metrics",
	}
	for _, o := range opts {
		o(srv)
	}
	for _, m := range srv.middlewares {
		mw, ok := mids.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
		}
		log.Infof("install middleware: %s", m)
		srv.Use(mw)
	}
	//超时中间件
	if srv.timeout > 0 {
		srv.Use(mids.TimeoutMiddleware(srv.timeout))
		log.Infof("install middleware: %s", "timeout")
	}

	//设置开发模式，打印路由信息
	if srv.mode != gin.DebugMode && srv.mode != gin.ReleaseMode && srv.mode != gin.TestMode {
		srv.mode = gin.ReleaseMode
	}
	//设置开发模式
	gin.SetMode(srv.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s(%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	return srv
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
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("http", s.tlsConf != nil), addr)
	}
	return s.err
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	https://127.0.0.1:8000
//	Legacy: http://127.0.0.1:8000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

// start rest server
func (s *Server) Start(ctx context.Context) error {
	//根据配置初始化pprof路由
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}
	if s.enableMetrics {

	}

	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.server = &http.Server{
		Addr:      s.address,
		Handler:   s.Engine,
		TLSConfig: s.tlsConf,
	}
	log.Infof("[HTTP] server listening on: %s", s.lis.Addr().String())
	var err error
	if s.tlsConf != nil {
		err = s.server.ServeTLS(s.lis, "", "")
	} else {
		err = s.server.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Infof("rest server is stopping")
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("rest server shutdown error: %s", err.Error())
		return err
	}
	log.Info("rest server stopped")
	return nil
}
