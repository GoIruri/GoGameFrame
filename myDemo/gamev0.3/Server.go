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

func (pr *PingRouter) PreHandle(request giface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	if _, err := request.GetConnection().GetTcpConnection().Write([]byte("before ping...\n")); err != nil {
		fmt.Println("call back ping error")
	}
}

func (pr *PingRouter) Handle(request giface.IRequest) {
	fmt.Println("Call Router Handle...")
	if _, err := request.GetConnection().GetTcpConnection().Write([]byte("ping...ping...ping\n")); err != nil {
		fmt.Println("ping...ping...ping error")
	}
}

func (pr *PingRouter) PostHandle(request giface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	if _, err := request.GetConnection().GetTcpConnection().Write([]byte("after ping\n")); err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	//	1创建一个server句柄,使用api
	s := gnet.NewServer("[game V0.2]")
	//2给当前框架添加一个自定义的Router
	s.AddRouter(&PingRouter{})
	//	3启动server
	s.Serve()
}
