package tcp

import (
	"time"
)

const (
	defaultServerAddr                   = ":3553"
	defaultServerMaxMsgLen              = 1024
	defaultServerMaxConnNum             = 5000
	defaultServerHeartbeatCheck         = false
	defaultServerHeartbeatCheckInterval = 10
)

type ServerOption func(o *serverOptions)

type serverOptions struct {
	addr                   string        // 监听地址，默认0.0.0.0:3553
	maxMsgLen              int           // 最大消息长度，默认1K
	maxConnNum             int           // 最大连接数，默认5000
	enableHeartbeatCheck   bool          // 是否启用心跳检测，默认不启用
	heartbeatCheckInterval time.Duration // 心跳检测间隔时间，默认10s
}

func defaultServerOptions() *serverOptions {
	return &serverOptions{
		addr:                   defaultServerAddr,
		maxMsgLen:              defaultServerMaxMsgLen,
		maxConnNum:             defaultServerMaxConnNum,
		enableHeartbeatCheck:   defaultServerHeartbeatCheck,
		heartbeatCheckInterval: defaultServerHeartbeatCheckInterval * time.Second,
	}
}

// WithServerListenAddr 设置监听地址
func WithServerListenAddr(addr string) ServerOption {
	return func(o *serverOptions) { o.addr = addr }
}

// WithServerMaxMsgLen 设置消息最大长度
func WithServerMaxMsgLen(maxMsgLen int) ServerOption {
	return func(o *serverOptions) { o.maxMsgLen = maxMsgLen }
}

// WithServerMaxConnNum 设置连接的最大连接数
func WithServerMaxConnNum(maxConnNum int) ServerOption {
	return func(o *serverOptions) { o.maxConnNum = maxConnNum }
}

// WithServerEnableHeartbeatCheck 是否启用心跳检测
func WithServerEnableHeartbeatCheck(enable bool) ServerOption {
	return func(o *serverOptions) { o.enableHeartbeatCheck = enable }
}

// WithServerHeartbeatInterval 设置心跳检测间隔时间
func WithServerHeartbeatInterval(heartbeatInterval time.Duration) ServerOption {
	return func(o *serverOptions) { o.heartbeatCheckInterval = heartbeatInterval }
}
