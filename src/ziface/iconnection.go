package ziface

import "net"

type IConnection interface {
	//  启动连接
	Start()
	//  停止连接
	Stop()
	// 获取连接绑定的TCP socket
	GetTCPConnection() *net.TCPConn
	// 获取连接ID
	GetConnID() uint32
	// 获取远端地址
	RemoteAddr() net.Addr
	// 发送消息
	SendMsg(uint32, []byte) error

	// 获取属性
	GetProperty(string) (interface{}, error)
	// 设置属性
	SetProperty(string, interface{})
	// 移除属性
	RemoveProperty(string)
}

type ConnectionHandler func(*net.TCPConn, []byte, int) error
