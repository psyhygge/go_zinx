package ziface

type IDataPack interface {
	// GetHeadLen 获取消息头长度
	GetHeadLen() uint32

	// Pack 封包方法，创建一个msg包
	Pack(msg IMessage) ([]byte, error)

	// Unpack 拆包方法，将包的head信息处理完成之后，得到msg消息
	Unpack(binaryData []byte) (IMessage, error)
}
