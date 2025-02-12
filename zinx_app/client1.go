package main

import (
	"fmt"
	"go_zinx/znet"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client1 start")

	time.Sleep(time.Second * 1)

	// 连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client1 start err, exit!")
		return
	}

	for {
		// 调用write写数据
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx Artist client1 Test Message")))
		if err != nil {
			fmt.Println("client1 pack err ", err)
			return
		}

		if _, err = conn.Write(binaryMsg); err != nil {
			fmt.Println("client1 write err ", err)
			return
		}

		// 接收服务器的数据
		binaryHead := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, binaryHead)
		if err != nil {
			fmt.Println("client1 read head err ", err)
			return
		}
		msg, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client1 unpack err ", err)
			return
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("client1 read msg data err ", err)
				return
			}
		}
		msg.SetData(data)
		fmt.Printf("client1 recv server msg: %s\n", msg.GetData())

		// cpu 阻塞
		time.Sleep(1 * time.Second)
	}
}
