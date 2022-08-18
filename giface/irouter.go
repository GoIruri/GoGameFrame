package giface

// IRouter 路由抽象接口, 路由里的数据都是IRequest
type IRouter interface {
	PreHandle(request IRequest)

	Handle(request IRequest)

	PostHandle(request IRequest)
}
