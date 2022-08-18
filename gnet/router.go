package gnet

import (
	"zinx/giface"
)

// BaseRouter 实现Router时, 先嵌入BaseRouter结构体, 然后根据需要对BaseRouter的方法进行重写就好了
type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request giface.IRequest) {
}

func (br *BaseRouter) Handle(request giface.IRequest) {
}

func (br *BaseRouter) PostHandle(request giface.IRequest) {
}
