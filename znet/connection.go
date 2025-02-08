package znet

import (
	"fmt"
	"go_zinx/ziface"
	"net"
)

type Connection struct {
	Conn      *net.TCPConn      // tcp socket
	ConnID    uint32            // 连接id
	IsClosed  bool              // 连接是否关闭
	handleApi ziface.HandleFunc // 该连接所绑定的处理api方法
	ExitChan  chan bool         // 告知当前连接停止的 channel
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read from client failed, ", err)
			continue
		}

		// 调用当前连接所绑定的handleAPI执行
		if err := c.handleApi(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID: ", c.ConnID, " handle msg err: ", err)
			break
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("connection start... ConnID:", c.ConnID)
	// TODO 启动从当前连接读数据业务
	go c.StartReader()

	// TODO 启动从当前连接写数据业务
}

func (c *Connection) Stop() {
	fmt.Println("connection stop... ConnID:", c.ConnID)

	if c.IsClosed {
		return
	}
	c.IsClosed = true

	c.Conn.Close()
	close(c.ExitChan)
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

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}

func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandleFunc) ziface.IConnection {
	return &Connection{
		Conn:      conn,
		ConnID:    connID,
		IsClosed:  false,
		handleApi: callbackApi,
		ExitChan:  make(chan bool, 1),
	}
}
