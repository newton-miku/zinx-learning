package znet

import "zinx/ziface"

// 实现Router接口时，先嵌入BaseRouter,然后根据需要重写
type BaseRouter struct {
}

// 处理方法默认留空，根据需要重写
func (br *BaseRouter) PreHandle(request ziface.IRequest)  {}
func (br *BaseRouter) Handle(request ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
