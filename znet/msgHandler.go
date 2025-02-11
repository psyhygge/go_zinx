package znet

import (
	"fmt"
	"go_zinx/utils"
	"go_zinx/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter // 存放每个MsgID对应的处理方法
	TaskQueue      []chan ziface.IRequest    // 负责Worker取任务的消息队列, TaskQueue[0]对应worker0, TaskQueue[1]对应worker1...
	WorkerPoolSize uint32                    // 业务工作Worker池的worker数量
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
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// StartWorkerPool 启动一个Worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 创建当前Worker对应的TaskQueue
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskNum)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workerId int, taskChan chan ziface.IRequest) {
	fmt.Println("Worker ", workerId, " is started.")

	for {
		select {
		case request := <-taskChan:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue, 由worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配算法
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgID(), " to workerID=", workerId)
	// 将消息发送给worker的taskChannel即可
	mh.TaskQueue[workerId] <- request
}
