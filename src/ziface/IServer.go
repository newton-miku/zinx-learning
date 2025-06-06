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
	//  设置OnConnStart函数
	SetOnConnStart(func(connection IConnection))
	//  设置OnConnStop函数
	SetOnConnStop(func(connection IConnection))
	//  调用OnConnStart函数
	CallOnConnStart(connection IConnection)
	//  调用OnConnStop函数
	CallOnConnStop(connection IConnection)
}
