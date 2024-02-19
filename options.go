package lori

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/cr-mao/lori/log"
	"github.com/cr-mao/lori/registry"
	"github.com/cr-mao/lori/transport"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL // 暴露的地址

	ctx  context.Context // 上下文，可以传入进来，一般是不需要的。应该就是background context 出发的
	sigs []os.Signal     // 注册信号

	logger           log.Logger
	registrar        registry.Registrar // 服务注册
	registrarTimeout time.Duration      // 服务注册超时
	stopTimeout      time.Duration      // 停止超时时间，可以给很大的值，看情况自己定
	servers          []transport.Server //  有哪些server

	// Before and After funcs  钩子函数
	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// ID with service id.
func WithID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name with service name.
func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func WithVersion(version string) Option {
	return func(o *options) { o.version = version }
}

// Metadata with service metadata.
func WithMetadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

// Endpoint with service endpoint.
func WithEndpoint(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

// Context with service context.
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Logger with service logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// Server with transport servers.
func WithServer(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Registrar with service registry.
func WithRegistrar(r registry.Registrar) Option {
	return func(o *options) { o.registrar = r }
}

// RegistrarTimeout with registrar timeout.
func WithRegistrarTimeout(t time.Duration) Option {
	return func(o *options) { o.registrarTimeout = t }
}

// StopTimeout with app stop timeout.
func WithStopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

// Before and Afters

// BeforeStart run funcs before app starts
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

// BeforeStop run funcs before app stops
func BeforeStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn)
	}
}

// AfterStart run funcs after app starts
func AfterStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

// AfterStop run funcs after app stops
func AfterStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}
