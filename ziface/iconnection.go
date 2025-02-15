package ziface

import "net"

type IConnection interface {
	// Start 启动连接
	Start()
	// Stop 停止连接
	Stop()
	// GetTCPConnection 获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接模块的ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的TCP状态 IP:Port
	RemoteAddr() net.Addr
	// SendMsg 发送数据，将数据发给远程客户端
	SendMsg(msgId uint32, data []byte) error
	// ReadMsg 读取客户端发送的消息
	ReadMsg() (IMessage, error)
	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除连接属性
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
