package ziface

type IMsgHandler interface {
	// DoMsgHandler 执行对应Router的消息处理方法
	DoMsgHandler(request IRequest)

	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
}
