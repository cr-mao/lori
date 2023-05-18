package tcp

import (
	"time"
)

const (
	defaultClientDialAddr          = "127.0.0.1:3553"
	defaultClientMaxMsgLen         = 1024
	defaultClientHeartbeat         = false
	defaultClientHeartbeatInterval = 10
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	addr              string        // 地址
	maxMsgLen         int           // 最大消息长度
	enableHeartbeat   bool          // 是否启用心跳，默认不启用
	heartbeatInterval time.Duration // 心跳间隔时间，默认10s
}

func defaultClientOptions() *clientOptions {
	return &clientOptions{
		addr:              defaultClientDialAddr,
		maxMsgLen:         defaultClientMaxMsgLen,
		enableHeartbeat:   defaultClientHeartbeat,
		heartbeatInterval: defaultClientHeartbeatInterval * time.Second,
	}
}

// WithClientDialAddr 设置拨号地址
func WithClientDialAddr(addr string) ClientOption {
	return func(o *clientOptions) { o.addr = addr }
}

// WithClientMaxMsgLen 设置消息最大长度
func WithClientMaxMsgLen(maxMsgLen int) ClientOption {
	return func(o *clientOptions) { o.maxMsgLen = maxMsgLen }
}

// WithClientEnableHeartbeat 设置是否启用心跳间隔时间
func WithClientEnableHeartbeat(enable bool) ClientOption {
	return func(o *clientOptions) { o.enableHeartbeat = enable }
}

// WithClientHeartbeatInterval 设置心跳间隔时间
func WithClientHeartbeatInterval(heartbeatInterval time.Duration) ClientOption {
	return func(o *clientOptions) { o.heartbeatInterval = heartbeatInterval }
}
