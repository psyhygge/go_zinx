package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"go_zinx/utils"
	"go_zinx/ziface"
)

type DataPack struct {
}

// GetHeadLen 获取消息头长度
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

// Pack 封包方法，创建一个msg包
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuffer := bytes.NewBuffer([]byte{})
	// 将dataLen写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将id写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}

// Unpack 拆包方法，将包的head信息处理完成之后，得到msg消息
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}
	// 只解压head的信息，得到dataLen和msgID

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen的长度是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data received")
	}

	return msg, nil
}

// NewDataPack 拆包封包实例的初试化
func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}
