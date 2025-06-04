package ziface

type IMsgHandler interface {
	// 路由调度函数，用于处理请求
	DoMsgHandler(request IRequest)
	//  添加路由
	AddRouter(msgID uint32, router IRouter)
	// 创建Worker工作池
	CreatWorkerPool()
	// 将请求发送到Worker工作池中
	SendReqToWorker(req IRequest)
}
