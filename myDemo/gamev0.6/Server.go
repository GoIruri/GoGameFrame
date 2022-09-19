package main

import (
	"fmt"
	"zinx/giface"
	"zinx/gnet"
)

/*
基于Zinx框架来开发的服务端应用程序
*/

// PingRouter ping test自定义路由
type PingRouter struct {
	gnet.BaseRouter
}

func (pr *PingRouter) Handle(request giface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), "data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// Hello test自定义路由
type Hello struct {
	gnet.BaseRouter
}

// Handle Test
func (h *Hello) Handle(request giface.IRequest) {
	fmt.Println("Call Hello Handle...")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), "data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//	1创建一个server句柄,使用api
	s := gnet.NewServer("[game V0.5]")
	//2给当前框架添加一个自定义的Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &Hello{})
	//	3启动server
	s.Serve()
}
