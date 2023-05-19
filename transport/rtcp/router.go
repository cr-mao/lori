package rtcp

import (
	"strconv"
	"sync"

	"github.com/cr-mao/lori/transport/rtcp/iface"
)

// BaseRouter is used as the base class when implementing a router.
// Depending on the needs, the methods of this base class can be overridden.
// (实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写)
type BaseRouter struct{}

// Here, all of BaseRouter's methods are empty, because some routers may not want to have PreHandle or PostHandle.
// Therefore, inheriting all routers from BaseRouter has the advantage that PreHandle and PostHandle do not need to be
// implemented to instantiate a router.
// (这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化)

// PreHandle -
func (br *BaseRouter) PreHandle(req iface.IRequest) {}

// Handle -
func (br *BaseRouter) Handle(req iface.IRequest) {}

// PostHandle -
func (br *BaseRouter) PostHandle(req iface.IRequest) {}

// New slice-based router
// The new version of the router has basic logic that allows users to pass in varying numbers of router handlers.
// The router will save all of these router handler functions and find them when a request comes in, then execute them using IRequest.
// The router can set globally shared components using the Use method.
// The router can be grouped using Group, and groups also have their own Use method for setting group-shared components.
// (新切片集合式路由
// 新版本路由基本逻辑,用户可以传入不等数量的路由路由处理器
// 路由本体会讲这些路由处理器函数全部保存,在请求来的时候找到，并交由IRequest去执行
// 路由可以设置全局的共用组件通过Use方法
// 路由可以分组,通过Group,分组也有自己对应Use方法设置组共有组件)

type RouterSlices struct {
	Apis     map[uint32][]iface.RouterHandler
	Handlers []iface.RouterHandler
	sync.RWMutex
}

func NewRouterSlices() *RouterSlices {
	return &RouterSlices{
		Apis:     make(map[uint32][]iface.RouterHandler, 10),
		Handlers: make([]iface.RouterHandler, 0, 6),
	}
}

func (r *RouterSlices) Use(handles ...iface.RouterHandler) {
	r.Handlers = append(r.Handlers, handles...)
}

func (r *RouterSlices) AddHandler(msgId uint32, Handlers ...iface.RouterHandler) {
	// 1. Check if the API handler method bound to the current msg already exists
	if _, ok := r.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}

	finalSize := len(r.Handlers) + len(Handlers)
	mergedHandlers := make([]iface.RouterHandler, finalSize)
	copy(mergedHandlers, r.Handlers)
	copy(mergedHandlers[len(r.Handlers):], Handlers)
	r.Apis[msgId] = append(r.Apis[msgId], mergedHandlers...)
}

func (r *RouterSlices) GetHandlers(MsgId uint32) ([]iface.RouterHandler, bool) {
	r.RLock()
	defer r.RUnlock()
	handlers, ok := r.Apis[MsgId]
	return handlers, ok
}

func (r *RouterSlices) Group(start, end uint32, Handlers ...iface.RouterHandler) iface.IGroupRouterSlices {
	return NewGroup(start, end, r, Handlers...)
}

type GroupRouter struct {
	start    uint32
	end      uint32
	Handlers []iface.RouterHandler
	router   iface.IRouterSlices
}

func NewGroup(start, end uint32, router *RouterSlices, Handlers ...iface.RouterHandler) *GroupRouter {
	g := &GroupRouter{
		start:    start,
		end:      end,
		Handlers: make([]iface.RouterHandler, 0, len(Handlers)),
		router:   router,
	}
	g.Handlers = append(g.Handlers, Handlers...)
	return g
}

func (g *GroupRouter) Use(Handlers ...iface.RouterHandler) {
	g.Handlers = append(g.Handlers, Handlers...)
}

func (g *GroupRouter) AddHandler(MsgId uint32, Handlers ...iface.RouterHandler) {
	if MsgId < g.start || MsgId > g.end {
		panic("add router to group err in msgId:" + strconv.Itoa(int(MsgId)))
	}

	finalSize := len(g.Handlers) + len(Handlers)
	mergedHandlers := make([]iface.RouterHandler, finalSize)
	copy(mergedHandlers, g.Handlers)
	copy(mergedHandlers[len(g.Handlers):], Handlers)

	g.router.AddHandler(MsgId, mergedHandlers...)
}
