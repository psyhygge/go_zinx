package znet

import (
	"fmt"
	"go_zinx/ziface"
	"net"
)

// Server IServer接口实现
type Server struct {
	Name      string // 服务器名称
	IPVersion string // 服务器绑定的ip版本
	IP        string // 服务器监听的IP
	Port      int    // 服务器监听的端口号
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient is Called...")
	fmt.Printf("[Receive Client Buf] %s, cnt = %d\n", data, cnt)
	_, err := conn.Write(data[:cnt])
	if err != nil {
		fmt.Println("write back buf err: ", err)
		return err
	}
	return nil
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP:%s, Port %d, is starting\n", s.IP, s.Port)

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
		// 3.阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp err: ", err)
				continue
			}
			// 客户端正常连接，创建协程进行业务处理
			dealConn := NewConnection(conn, cid, CallBackToClient)
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
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
}
