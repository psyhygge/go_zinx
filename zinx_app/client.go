package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start")

	time.Sleep(time.Second * 1)

	// 连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 调用write写数据
		_, err := conn.Write([]byte("hello zinx!!!"))
		if err != nil {
			fmt.Println("write conn err ", err)
			return
		}

		// 接收服务端数据
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read conn err ", err)
			return
		}

		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)

		// cpu 阻塞
		time.Sleep(1 * time.Second)
	}
}
