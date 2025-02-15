package ziface

type IConnectionManager interface {
	// Add 添加一个链接
	Add(conn IConnection)
	// Remove 删除一个链接
	Remove(conn IConnection)
	// Get 获取链接
	Get(connID uint32) (IConnection, error)
	// Len 获取当前连接总数
	Len() int
	// ClearConn 清空链接
	ClearConn()
}
