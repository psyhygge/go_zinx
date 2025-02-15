package ziface

type IMsgHandler interface {
	// DoMsgHandler 执行对应Router的消息处理方法
	DoMsgHandler(request IRequest)

	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)

	// StartWorkerPool 启动一个Worker工作池
	StartWorkerPool()

	// StartOneWorker 启动一个Worker工作流程
	StartOneWorker(workerId int, taskChan chan IRequest)

	// SendMsgToTaskQueue 将消息交给TaskQueue, 由worker进行处理
	SendMsgToTaskQueue(request IRequest)

	// StopWorkerPool 停止worker工作池
	StopWorkerPool()
}
