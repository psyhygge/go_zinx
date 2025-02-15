package main

import (
	"fmt"
	"go_zinx/znet"
	"io"
	"net"
	"time"
)

// 发送数据的 goroutine
func writeLoop(conn net.Conn) {
	defer conn.Close()
	dp := znet.NewDataPack()

	for {
		// 调用write写数据
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx psyhygge client Test Message")))
		if err != nil {
			fmt.Println("client pack err ", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("client write err ", err)
			return
		}

		// CPU 阻塞，避免写过快
		time.Sleep(1 * time.Second)
	}
}

// 接收数据的 goroutine
func readLoop(conn net.Conn) {
	defer conn.Close()
	dp := znet.NewDataPack()

	for {
		// 接收服务器的数据
		binaryHead := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(conn, binaryHead)
		if err != nil {
			fmt.Println("client read head err ", err)
			return
		}

		// 解包消息头
		msg, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack err ", err)
			return
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			// 读取消息体数据
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("client read msg data err ", err)
				return
			}
		}
		msg.SetData(data)
		fmt.Printf("client recv server msg: %s\n", msg.GetData())
	}
}

func main() {
	fmt.Println("client start")

	time.Sleep(time.Second * 1)

	// 连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	// 启动写数据的 goroutine
	go writeLoop(conn)

	// 启动读数据的 goroutine
	go readLoop(conn)

	select {}
	//for {
	//	// 调用write写数据
	//	dp := znet.NewDataPack()
	//	binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx psyhygge client Test Message")))
	//	if err != nil {
	//		fmt.Println("client pack err ", err)
	//		return
	//	}
	//
	//	if _, err = conn.Write(binaryMsg); err != nil {
	//		fmt.Println("client write err ", err)
	//		return
	//	}
	//
	//	// 接收服务器的数据
	//	binaryHead := make([]byte, dp.GetHeadLen())
	//	_, err = io.ReadFull(conn, binaryHead)
	//	if err != nil {
	//		fmt.Println("client read head err ", err)
	//		return
	//	}
	//	msg, err := dp.Unpack(binaryHead)
	//	if err != nil {
	//		fmt.Println("client unpack err ", err)
	//		return
	//	}
	//	var data []byte
	//	if msg.GetMsgLen() > 0 {
	//		data = make([]byte, msg.GetMsgLen())
	//		if _, err := io.ReadFull(conn, data); err != nil {
	//			fmt.Println("client read msg data err ", err)
	//			return
	//		}
	//	}
	//	msg.SetData(data)
	//	fmt.Printf("client recv server msg: %s\n", msg.GetData())
	//
	//	// cpu 阻塞
	//	time.Sleep(1 * time.Second)
	//}
}
