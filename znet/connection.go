package znet

import (
	"errors"
	"fmt"
	"go_zinx/utils"
	"go_zinx/ziface"
	"io"
	"net"
)

type Connection struct {
	Conn       *net.TCPConn       // tcp socket
	ConnID     uint32             // 连接id
	IsClosed   bool               // 连接是否关闭
	ExitChan   chan bool          // 告知当前连接停止的 channel
	MsgChan    chan []byte        // 用于读写消息的 channel
	MsgHandler ziface.IMsgHandler // 消息管理模块
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " conn writer exit!")

	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err)
				return
			}
		case <-c.ExitChan:
			// 表示Reader已经退出
			return
		}
	}
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		msg, err := c.ReadMsg()
		if err != nil {
			fmt.Println("read msg error: ", err)
			break
		}

		// 得到Request数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 开启了工作池机制，将消息发送给TaskQueue, 由Worker Pool进行处理
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			// 从路由中找到注册绑定的Conn对应的router调用
			go c.MsgHandler.DoMsgHandler(req)
		}

	}
}

func (c *Connection) Start() {
	fmt.Println("connection start... ConnID:", c.ConnID)
	// TODO 启动从当前连接读数据业务
	go c.StartReader()

	// TODO 启动从当前连接写数据业务
	go c.StartWriter()
}

func (c *Connection) Stop() {
	fmt.Println("connection stop... ConnID:", c.ConnID)

	if c.IsClosed {
		return
	}
	c.IsClosed = true
	// 连接关闭
	c.Conn.Close()
	// 告知Writer退出
	c.ExitChan <- true
	close(c.ExitChan)
	close(c.MsgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("connection closed when send msg")
	}

	// 将data进行封包 MsgData包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}

	c.MsgChan <- binaryMsg
	return nil
}

func (c *Connection) ReadMsg() (ziface.IMessage, error) {
	// 创建一个拆包解包的对象
	dp := NewDataPack()
	headData := make([]byte, dp.GetHeadLen())
	// 读取客户端的MsgHead 8个字节
	if _, err := io.ReadFull(c.Conn, headData); err != nil {
		fmt.Println("read buf error: ", err)
		return nil, err
	}
	msg, err := dp.Unpack(headData)
	if err != nil {
		fmt.Println("unpack err ", err)
		return nil, err
	}
	var data []byte
	if msg.GetMsgLen() > 0 {
		data = make([]byte, msg.GetMsgLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
			fmt.Println("read data buf error: ", err)
			return nil, err
		}
	}
	msg.SetData(data)
	return msg, nil
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) ziface.IConnection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
		MsgHandler: msgHandler,
	}
}
