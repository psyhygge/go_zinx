package znet

import (
	"fmt"
	"go_zinx/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个MsgID对应的处理方法
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
		Apis: make(map[uint32]ziface.IRouter),
	}
}
