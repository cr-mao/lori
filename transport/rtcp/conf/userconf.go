package conf

// UserConfToGlobal, Note that if UserConf is used,
// the method should be called to synchronize with GlobalConfObject
// because other parameters are called from this structure parameter.
// (注意如果使用UserConf应该调用方法同步至 GlobalConfObject 因为其他参数是调用的此结构体参数)
func UserConfToGlobal(config *Config) {
	// Server
	if config.Name != "" {
		GlobalObject.Name = config.Name
	}
	if config.Host != "" {
		GlobalObject.Host = config.Host
	}
	if config.TCPPort != 0 {
		GlobalObject.TCPPort = config.TCPPort
	}

	// Zinx
	if config.Version != "" {
		GlobalObject.Version = config.Version
	}
	if config.MaxPacketSize != 0 {
		GlobalObject.MaxPacketSize = config.MaxPacketSize
	}
	if config.MaxConn != 0 {
		GlobalObject.MaxConn = config.MaxConn
	}
	if config.WorkerPoolSize != 0 {
		GlobalObject.WorkerPoolSize = config.WorkerPoolSize
	}
	if config.MaxWorkerTaskLen != 0 {
		GlobalObject.MaxWorkerTaskLen = config.MaxWorkerTaskLen
	}
	if config.MaxMsgChanLen != 0 {
		GlobalObject.MaxMsgChanLen = config.MaxMsgChanLen
	}
	if config.IOReadBuffSize != 0 {
		GlobalObject.IOReadBuffSize = config.IOReadBuffSize
	}

	// Keepalive
	if config.HeartbeatMax != 0 {
		GlobalObject.HeartbeatMax = config.HeartbeatMax
	}

	// TLS
	if config.CertFile != "" {
		GlobalObject.CertFile = config.CertFile
	}
	if config.PrivateKeyFile != "" {
		GlobalObject.PrivateKeyFile = config.PrivateKeyFile
	}

	if config.Mode != "" {
		GlobalObject.Mode = config.Mode
	}
	if config.WsPort != 0 {
		GlobalObject.WsPort = config.WsPort
	}

	if config.RouterSlicesMode {
		GlobalObject.RouterSlicesMode = config.RouterSlicesMode
	}
}
