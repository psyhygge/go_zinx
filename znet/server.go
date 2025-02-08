package znet

import (
	"fmt"
	"go_zinx/utils"
	"go_zinx/ziface"
	"net"
)

// Server IServer接口实现
type Server struct {
	Name      string         // 服务器名称
	IPVersion string         // 服务器绑定的ip版本
	IP        string         // 服务器监听的IP
	Port      int            // 服务器监听的端口号
	Router    ziface.IRouter // router
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
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
		var cid uint32
		cid = 0
		// 3.阻塞等待客户端连接，处理客户端连接业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp err: ", err)
				continue
			}
			// 客户端正常连接，创建协程进行业务处理
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// TODO 做服务器的资源回收或停止
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
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
}
