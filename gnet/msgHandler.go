package gnet

import (
	"fmt"
	"strconv"
	"zinx/giface"
)

// MsgHandle 消息处理模块的实现
type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]giface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{Apis: make(map[uint32]giface.IRouter)}
}

func (mh *MsgHandle) DoMsgHandler(request giface.IRequest) {
	// 从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is NOT FOUND! Need Register!")
	}
	// 根据msgID调度对应Router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router giface.IRouter) {
	// 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id已经注册了
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	// 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, "succ!")
}
