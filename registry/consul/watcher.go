package consul

import (
	"context"

	"github.com/cr-mao/lori/registry"
)

type watcher struct {
	event chan struct{}
	set   *serviceSet

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

// 返回新的注册实例
func (w *watcher) Next() (services []*registry.ServiceInstance, err error) {
	select {
	case <-w.ctx.Done():
		err = w.ctx.Err()
	case <-w.event:
	}

	// 有变化的适合，event事件会过来， 就有新数据
	ss, ok := w.set.services.Load().([]*registry.ServiceInstance)

	if ok {
		services = append(services, ss...)
	}
	return
}

func (w *watcher) Stop() error {
	w.cancel()
	w.set.lock.Lock()
	defer w.set.lock.Unlock()
	delete(w.set.watcher, w)
	return nil
}
