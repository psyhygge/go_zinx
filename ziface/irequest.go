package ziface

/*
	把客户端 连接信息 和 请求的数据 包装到 Request
*/

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection

	// GetData 得到请求的数据
	GetData() []byte

	GetMsgID() uint32
}
