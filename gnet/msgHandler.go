package gnet

import (
	"fmt"
	"strconv"
	"zinx/giface"
	"zinx/utils"
)

// MsgHandle 消息处理模块的实现
type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]giface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan giface.IRequest
	// 业务工作Worker池的Worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]giface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan giface.IRequest, utils.GlobalObject.MaxPackageSize),
	}
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

// StartWorkerPool 启动一个Worker工作池（只能发生一次，框架只能有一个Worker工作池）
func (mh *MsgHandle) StartWorkerPool() {
	// 根据WorkerPoolSize 分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 给当前的worker对应的channel消息队列开辟空间
		mh.TaskQueue[i] = make(chan giface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前的worker，阻塞的等待消息从channel中传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerId int, taskQueue chan giface.IRequest) {
	fmt.Println("WorkerId = ", workerId, "is started...")

	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request giface.IRequest) {
	// 将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), " request MsgID = ",
		request.GetMsgID(),
		"to workerID = ", workerID)
	// 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
