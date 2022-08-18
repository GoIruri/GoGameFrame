package giface

// Iserver 定义一个服务器接口
type Iserver interface {
	Start()

	Stop()

	Serve()
	// AddRouter 路由功能:给当前的服务注册一个路由方法, 供客户端的链接处理使用
	AddRouter(router IRouter)
}
