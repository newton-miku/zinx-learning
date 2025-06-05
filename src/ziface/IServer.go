package ziface

type IServer interface {
	//  启动服务器
	Start()
	//  停止服务器
	Stop()
	//  服务器初始化
	Init()
	//  注册路由
	AddRouter(msgID uint32, router IRouter)
	//  获取对应的连接管理器
	GetConnectionManager() IConnectionManager
}
