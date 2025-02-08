package main

import (
	"fmt"
	"go_zinx/ziface"
	"go_zinx/znet"
)

/*
	基于zinx框架的应用程序
*/

// PingRouter 用户自定义router
type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("PreHandle PingRouter"))
	if err != nil {
		fmt.Println("PreHandle write conn err ", err)
		return
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Handle PingRouter"))
	if err != nil {
		fmt.Println("Handle write conn err ", err)
		return
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("PostHandle PingRouter"))
	if err != nil {
		fmt.Println("PostHandle write conn err ", err)
		return
	}
}

func main() {
	// 创建一个server服务
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	// 运行服务
	s.Serve()
}
