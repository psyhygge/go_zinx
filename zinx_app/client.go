package main

import (
	"fmt"
	"go_zinx/znet"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

// 发送数据的 goroutine
func writeLoop(conn net.Conn, stopChan chan struct{}) {
	defer wg.Done()
	dp := znet.NewDataPack()

	for {
		select {
		case <-stopChan:
			fmt.Println("writeLoop received stop signal, exiting...")
			return
		default:
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
}

// 接收数据的 goroutine
func readLoop(conn net.Conn, stopChan chan struct{}) {
	defer wg.Done()
	dp := znet.NewDataPack()

	for {
		select {
		case <-stopChan:
			fmt.Println("readLoop received stop signal, exiting...")
			return
		default:
			// 接收服务器的数据
			binaryHead := make([]byte, dp.GetHeadLen())
			_, err := io.ReadFull(conn, binaryHead)
			if err != nil {
				if err == io.EOF {
					fmt.Println("server closed")
					close(stopChan) // 通知其他 goroutine 退出
					return
				}
				fmt.Println("client read head err ", err)
				conn.Close()
				close(stopChan) // 通知其他 goroutine 退出
				return
			}

			// 解包消息头
			msg, err := dp.Unpack(binaryHead)
			if err != nil {
				fmt.Println("client unpack err ", err)
				close(stopChan) // 通知其他 goroutine 退出
				return
			}

			var data []byte
			if msg.GetMsgLen() > 0 {
				// 读取消息体数据
				data = make([]byte, msg.GetMsgLen())
				if _, err := io.ReadFull(conn, data); err != nil {
					fmt.Println("client read msg data err ", err)
					close(stopChan) // 通知其他 goroutine 退出
					return
				}
			}
			msg.SetData(data)
			fmt.Printf("client recv server msg: %s\n", msg.GetData())
		}
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

	// 创建一个 channel 用于通知 goroutine 退出
	stopChan := make(chan struct{})

	// 启动写数据的 goroutine
	wg.Add(1)
	go writeLoop(conn, stopChan)

	// 启动读数据的 goroutine
	wg.Add(1)
	go readLoop(conn, stopChan)

	// 设置一个信号捕获器，用于优雅地停止服务
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) // 监听终止信号

	// 等待中断信号或 stopChan 信号
	select {
	case <-signalChan:
		fmt.Println("\nReceived stop signal, shutting down the client...")
		// 通知 goroutine 退出
		close(stopChan)
	case <-stopChan:
		fmt.Println("\nServer closed, shutting down the client...")
	}

	// 关闭连接
	conn.Close()

	// 等待所有 goroutine 退出
	wg.Wait()

	// 等待一段时间确保服务完全停止
	time.Sleep(2 * time.Second)
	fmt.Println("Client has been stopped.")
}
