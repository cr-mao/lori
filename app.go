package lori

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cr-mao/lori/log"
	"github.com/cr-mao/lori/registry"
	"github.com/cr-mao/lori/transport"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// AppInfo is application context value.
type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

type App struct {
	opts   options
	ctx    context.Context
	locker sync.Mutex

	instance *registry.ServiceInstance

	cancel func()
}

func New(opts ...Option) *App {
	o := options{
		ctx:              context.Background(),
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
		stopTimeout:      10 * time.Second, //注意长链接，这个时间要给足够长。
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.logger != nil {
		log.SetLogger(o.logger)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {
	//app 实例信息
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	a.locker.Lock()
	a.instance = instance
	a.locker.Unlock()
	sctx := NewContext(a.ctx, a)
	eg, ctx := errgroup.WithContext(sctx)
	wg := sync.WaitGroup{}

	for _, fn := range a.opts.beforeStart {
		if err = fn(sctx); err != nil {
			return err
		}
	}
	for _, srv := range a.opts.servers {
		// 复制一个，防止..
		server := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal , 收到信号后走Stop方法，然后把app的context cancel掉，那么它的子context就会关闭。
			stopCtx, cancel := context.WithTimeout(NewContext(a.opts.ctx, a), a.opts.stopTimeout)
			defer cancel()
			return server.Stop(stopCtx) //执行server的stop 下面的start 方法会停止阻塞。
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done() // here is to ensure server start has begun running before register, so defer is not needed
			return server.Start(sctx)
		})
	}
	//等待所有server启动
	wg.Wait()
	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(ctx, a.opts.registrarTimeout)
		defer rcancel()
		if err = a.opts.registrar.Register(rctx, instance); err != nil {
			return err
		}
	}
	for _, fn := range a.opts.afterStart {
		if err = fn(sctx); err != nil {
			return err
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})
	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	for _, fn := range a.opts.afterStop {
		err = fn(sctx)
	}
	return err
}

// Stop gracefully stops the application.
func (a *App) Stop() (err error) {
	sctx := NewContext(a.ctx, a)
	for _, fn := range a.opts.beforeStop {
		err = fn(sctx)
	}
	a.locker.Lock()
	instance := a.instance
	a.locker.Unlock()
	if a.opts.registrar != nil && instance != nil {
		ctx, cancel := context.WithTimeout(NewContext(a.ctx, a), a.opts.registrarTimeout)
		defer cancel()
		if err = a.opts.registrar.Deregister(ctx, instance); err != nil {
			return err
		}
	}
	if a.cancel != nil {
		a.cancel()
	}
	return err
}

func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0, len(a.opts.endpoints))
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if len(endpoints) == 0 {
		for _, srv := range a.opts.servers {
			if r, ok := srv.(transport.Endpointer); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, e.String())
			}
		}
	}
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Metadata:  a.opts.metadata,
		Endpoints: endpoints,
	}, nil
}

// ID returns app instance id. app实例id
func (a *App) ID() string { return a.opts.id }

// Name returns service name. app名
func (a *App) Name() string { return a.opts.name }

// Version returns app version. app版本
func (a *App) Version() string { return a.opts.version }

// Metadata returns service metadata.
func (a *App) Metadata() map[string]string { return a.opts.metadata }

// Endpoint returns endpoints.  app 服务的端点
func (a *App) Endpoint() []string {
	if a.instance != nil {
		return a.instance.Endpoints
	}
	return nil
}

type appKey struct{}

// NewContext returns a new Context that carries value. app 上下文信息
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any. 获得app 上下文信息
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInfo)
	return
}
