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
	// Send 发送数据，将数据发给远程客户端
	SendMsg(msgId uint32, data []byte) error

	ReadMsg() (IMessage, error)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
