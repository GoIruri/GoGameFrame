package main

import "zinx/gnet"

/*
基于Zinx框架来开发的服务端应用程序
*/
func main() {
	//	1创建一个server句柄,使用zinx的api
	s := gnet.NewServer("[game V0.2]")
	//	2启动server
	s.Serve()
}
