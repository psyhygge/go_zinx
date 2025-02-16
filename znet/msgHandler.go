package znet

import (
	"fmt"
	"go_zinx/utils"
	"go_zinx/ziface"
	"strconv"
	"sync"
	"time"
)

type MsgHandler struct {
	Apis                  map[uint32]ziface.IRouter // 存放每个MsgID对应的处理方法
	TaskQueue             []chan ziface.IRequest    // 负责Worker取任务的消息队列, TaskQueue[0]对应worker0, TaskQueue[1]对应worker1...
	WorkerPoolSize        uint32                    // 业务工作Worker池的worker数量
	WorkerGoroutineNum    uint32                    // 每个Worker队列的协程数量
	stopChan              chan struct{}             // 停止信号
	wg                    sync.WaitGroup            // 总等待组
	queueWg               []sync.WaitGroup          // 每个队列的等待组
	dynamicGoroutineCount uint32                    // 当前动态协程数量
	dynamicMutex          sync.Mutex                // 保护动态协程计数器的锁
	dynamicWg             sync.WaitGroup            // 动态协程的等待组
}

// DoMsgHandler 执行对应Router的消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; !ok {
		mh.Apis[msgId] = router
		fmt.Println("Add api msgId = ", msgId, " success!")
	} else {
		panic("repeated api, msgId = " + strconv.Itoa(int(msgId)))
	}
}

func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		Apis:                  make(map[uint32]ziface.IRouter),
		WorkerPoolSize:        utils.GlobalObject.WorkerPoolSize,
		WorkerGoroutineNum:    utils.GlobalObject.WorkerGoroutineNum,
		TaskQueue:             make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		stopChan:              make(chan struct{}),
		queueWg:               make([]sync.WaitGroup, utils.GlobalObject.WorkerPoolSize),
		dynamicGoroutineCount: 0,
		dynamicWg:             sync.WaitGroup{},
	}
}

// StartWorkerPool 启动一个Worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	// 初始化带缓冲的任务队列
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(
			chan ziface.IRequest,
			utils.GlobalObject.MaxWorkerTaskNum, // 从配置获取缓冲区大小
		)
	}

	// 为每个任务队列启动多个消费者协程
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		for j := 0; j < int(mh.WorkerGoroutineNum); j++ {
			mh.queueWg[i].Add(1)
			go mh.StartWorker(i, j, mh.TaskQueue[i])
		}
	}
}

// StartWorker 工作协程
func (mh *MsgHandler) StartWorker(queueID, workerID int, taskChan <-chan ziface.IRequest) {
	defer mh.queueWg[queueID].Done()

	fmt.Printf("Worker Queue[%d]-Goroutine[%d] started\n", queueID, workerID)

	for {
		select {
		case <-mh.stopChan: // 监听停止信号
			fmt.Printf("Worker Queue[%d]-Goroutine[%d] stopped\n", queueID, workerID)
			return
		case req, ok := <-taskChan:
			if !ok {
				fmt.Printf("Worker Queue[%d]-Goroutine[%d] exited\n", queueID, workerID)
				return
			}
			fmt.Printf("Worker Queue[%d]-Goroutine[%d] working\n", queueID, workerID)
			mh.wg.Add(1)
			mh.DoMsgHandler(req)
			mh.wg.Done()
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue, 由worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := mh.loadBalance(request.GetConnection().GetConnID())

	select {
	case mh.TaskQueue[workerID] <- request:
		// 成功写入队列
	case <-mh.stopChan:
		fmt.Println("SendMsgToTaskQueue stopped")
		return
	default:
		// 队列已满时的降级处理
		fmt.Printf("Worker Queue[%d] is full! ReqID=%d\n", workerID, request.GetMsgID())

		// 尝试动态增加协程处理请求
		if mh.tryStartDynamicGoroutine(request, mh.TaskQueue[workerID]) {
			return
		}

		// 如果动态协程数已达上限，直接在当前goroutine处理
		//mh.DoMsgHandler(request)
		// 暂时先不做降级处理
		fmt.Printf("Worker Queue is full! ReqID=%d\n", request.GetMsgID())
	}
}

// 负载均衡算法
func (mh *MsgHandler) loadBalance(connID uint32) int {
	// 一致性哈希算法，相同connID总是分配到同一个队列
	return int(connID) % int(mh.WorkerPoolSize)
}

// tryStartDynamicGoroutine 尝试启动动态协程处理请求
func (mh *MsgHandler) tryStartDynamicGoroutine(request ziface.IRequest, taskChan <-chan ziface.IRequest) bool {
	mh.dynamicMutex.Lock()
	defer mh.dynamicMutex.Unlock()

	// 检查是否超过最大动态协程数
	if mh.dynamicGoroutineCount >= utils.GlobalObject.MaxDynamicGoroutines {
		return false
	}

	// 增加动态协程计数
	mh.dynamicGoroutineCount++
	fmt.Println("Starting dynamic goroutine... Num=", mh.dynamicGoroutineCount)

	// 增加动态协程的WaitGroup计数
	mh.dynamicWg.Add(1)

	// 启动动态协程
	go mh.handleWithDynamicGoroutine(request, taskChan)

	return true
}

// handleWithDynamicGoroutine 动态协程处理请求
func (mh *MsgHandler) handleWithDynamicGoroutine(request ziface.IRequest, taskChan <-chan ziface.IRequest) {
	defer func() {
		// 减少动态协程计数
		mh.dynamicMutex.Lock()
		mh.dynamicGoroutineCount--
		mh.dynamicMutex.Unlock()

		// 通知WaitGroup当前协程已退出
		mh.dynamicWg.Done()
	}()

	// 处理请求
	mh.DoMsgHandler(request)

	// 动态协程空闲超时后退出
	timeout := time.NewTimer(time.Duration(utils.GlobalObject.DynamicGoroutineTimeout) * time.Second)
	defer timeout.Stop()

	for {
		select {
		case req, ok := <-taskChan:
			if !ok {
				fmt.Println("Dynamic goroutine exited")
				return
			}
			mh.DoMsgHandler(req)
		case <-timeout.C:
			// 超时后退出
			fmt.Println("Dynamic goroutine exited due to timeout")
			return
		case <-mh.stopChan:
			// 收到停止信号后退出
			fmt.Println("Dynamic goroutine exited due to stop signal")
			return
		}
	}
}

// StopWorkerPool 停止工作池
func (mh *MsgHandler) StopWorkerPool() {
	// 1. 停止接收新请求
	close(mh.stopChan)

	// 2. 等待所有进行中的请求处理完成
	mh.wg.Wait()

	// 3. 关闭所有任务队列
	for i := range mh.TaskQueue {
		close(mh.TaskQueue[i])
	}

	// 4. 等待所有worker协程退出
	for i := range mh.queueWg {
		mh.queueWg[i].Wait()
	}

	// 5. 等待所有动态协程退出
	mh.dynamicWg.Wait()

	fmt.Println("Worker pool fully stopped")
}
