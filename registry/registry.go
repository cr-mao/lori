package registry

import "context"

// 服务注册接口
type Registrar interface {
	//注册
	Register(ctx context.Context, service *ServiceInstance) error
	//注销
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// 服务发现接口
type Discovery interface {
	//获取服务实例
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	//创建服务监听器
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

// // Watcher is service watcher.
type Watcher interface {
	//获取服务实例, next在下面的情况下会返回服务
	//1. 第一次监听时，如果服务实例列表不为空，则返回服务实例列表
	//2. 如果服务实例发生变化，则返回服务实例列表
	// 如果上面两种情况都不满足，则会阻塞到context deadline或者cancel
	Next() ([]*ServiceInstance, error)
	//主动放弃监听
	Stop() error
}

type ServiceInstance struct {
	//注册到注册中心的服务id
	ID string `json:"id"`

	//服务名称
	Name string `json:"name"`

	//服务版本
	Version string `json:"version"`

	//服务元数据
	Metadata map[string]string `json:"metadata"`

	//http://127.0.0.1:8000
	//grpc://127.0.0.1:9000
	Endpoints []string `json:"endpoints"`
}
