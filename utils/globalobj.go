package utils

import (
	"encoding/json"
	"go_zinx/ziface"
	"os"
)

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer // 当前全局Server对象
	Host      string         // 当前服务器主机监听的ip
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	/*
		Zinx
	*/
	Version                 string // 当前zinx版本号
	MaxConn                 int    // 当前服务器主机允许的最大链接个数
	MaxPackageSize          uint32 // 当前zinx框架数据包的最大值
	WorkerPoolSize          uint32 // 当前业务工作Worker池的数量
	MaxWorkerTaskNum        uint32 // 当前每个worker对应的最大任务队列数量
	WorkerGoroutineNum      uint32 // 当前worker对应的最大任务队列数量
	MaxDynamicGoroutines    uint32 // 最大动态协程数量
	DynamicGoroutineTimeout uint32 // 协程超时时间
}

// 定义全局变量
var GlobalObject *GlobalObj

// Reload 从conf/zinx.json加载对应的参数到结构体中
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("zinx_app/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:                    "ZinxServerApp",
		Version:                 "V0.4",
		TcpPort:                 8999,
		Host:                    "0.0.0.0",
		MaxConn:                 1000,
		MaxPackageSize:          4096,
		WorkerPoolSize:          10,
		MaxWorkerTaskNum:        0,
		WorkerGoroutineNum:      2,
		MaxDynamicGoroutines:    2,
		DynamicGoroutineTimeout: 120,
	}

	// 读取conf/zinx.json文件，根据配置文件赋值GlobalObject
	GlobalObject.Reload()
}
