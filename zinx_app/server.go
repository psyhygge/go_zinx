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

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Handle...")
	// 读取客户端数据
	// 根据MsgID做不同的业务处理
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println("write back err ", err)
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
