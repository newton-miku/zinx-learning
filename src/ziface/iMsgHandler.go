package ziface

type IMsgHandler interface {
	// 路由调度函数，用于处理请求
	DoMsgHandler(request IRequest)
	//  添加路由
	AddRouter(msgID uint32, router IRouter)
}
