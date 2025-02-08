package main

import (
	"go_zinx/znet"
)

/*
	基于zinx框架的应用程序
*/

func main() {
	// 创建一个server服务
	s := znet.NewServer("***zinx***")
	// 运行服务
	s.Serve()
}
