package ziface

type IConnectionManager interface {
	//  新增连接
	Add(conn IConnection)
	//  删除连接
	Remove(conn IConnection)
	//  获取连接（通过ConnID）
	Get(connID uint32) (IConnection, error)
	//  获取当前连接总数
	Len() int
	//  清空连接
	ClearConn()
	//  遍历所有连接,并调用回调函数
	ForEach(func(IConnection)) (IConnection, error)
}
