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

// DoConnBegin 创建连接之后执行的钩子函数
func DoConnBegin(conn giface.IConnection) {
	fmt.Println("====> DoConnBegin is Called")
	err := conn.SendMsg(202, []byte("DOCONN BEGIN"))
	if err != nil {
		fmt.Println(err)
	}

	// 给当前的链接设置一些属性
	fmt.Println("set conn name, iruri...")
	conn.SetProperty("name", "GoIruri")
}

// DoConnLost 连接断开前需要执行的钩子函数
func DoConnLost(conn giface.IConnection) {
	fmt.Println("====> DoConnLost is Called")
	fmt.Println("conn ID = ", conn.GetConnID(), "is lost")

	// 获取链接属性
	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("name = ", name)
	}
}

func main() {
	//	1创建一个server句柄,使用api
	s := gnet.NewServer("[game V0.5]")
	//  2注册链接Hook函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)
	//	3给当前框架添加一个自定义的Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &Hello{})
	//	4启动server
	s.Serve()
}
