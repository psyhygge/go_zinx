package znet

import (
	"fmt"
	"go_zinx/utils"
	"go_zinx/ziface"
	"net"
)

// Server IServer接口实现
type Server struct {
	Name        string                        // 服务器名称
	IPVersion   string                        // 服务器绑定的ip版本
	IP          string                        // 服务器监听的IP
	Port        int                           // 服务器监听的端口号
	MsgHandler  ziface.IMsgHandler            // 消息管理模块
	ConnManager ziface.IConnectionManager     // 连接管理模块
	OnConnStart func(conn ziface.IConnection) // 连接创建时调用Hook函数
	OnConnStop  func(conn ziface.IConnection) // 连接断开时调用Hook函数
}

// SetOnConnStart 设置OnConnStart
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置OnConnStop
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用OnConnStart
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("[Zinx]=======> CallOnConnStart")
		s.OnConnStart(conn)
	} else {
		fmt.Println("[Zinx]=======> CallOnConnStart is nil")
	}
}

// CallOnConnStop 调用OnConnStop
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("[Zinx]=======> CallOnConnStop")
		s.OnConnStop(conn)
	} else {
		fmt.Println("[Zinx]=======> CallOnConnStop is nil")
	}
}

func (s *Server) GetConnMgr() ziface.IConnectionManager {
	return s.ConnManager
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name:%s, Host:%s, Port:%d\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	go func() {
		// 1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		// 2.监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("listen %s err: %s\n", s.IPVersion, err)
			return
		}
		fmt.Println("start zinx server success")

		// 启动worker工作池
		s.MsgHandler.StartWorkerPool()

		var cid uint32
		cid = 0
		// 3.阻塞等待客户端连接，处理客户端连接业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp err: ", err)
				continue
			}

			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				// TODO 链接数满了，关闭客户端连接
				fmt.Println("========>[Too many connections]<========")
				conn.Close()
				continue
			}

			// 客户端正常连接，创建协程进行业务处理
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// TODO 做服务器的资源回收或停止
	fmt.Println("[Zinx] Server Stop")

	// 先关闭工作池
	s.MsgHandler.StopWorkerPool()
	// 再清除连接
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	// 启动服务
	s.Start()

	// TODO 可以在启动之后做一些额外的业务，如服务注册

	// 阻塞
	select {}
}

// NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	return &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
}
