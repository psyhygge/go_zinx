package ziface

type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 添加router
	AddRouter(msgId uint32, router IRouter)
	// GetConnMgr 获取连接管理器
	GetConnMgr() IConnectionManager
	// SetOnConnStart 设置OnConnStart
	SetOnConnStart(func(conn IConnection))
	// SetOnConnStop 设置OnConnStop
	SetOnConnStop(func(conn IConnection))
	// CallOnConnStart 调用OnConnStart
	CallOnConnStart(conn IConnection)
	// CallOnConnStop 调用OnConnStop
	CallOnConnStop(conn IConnection)
}
