package znet

import "zinx/ziface"

type Request struct {
	// 与客户端建立的连接
	Conn ziface.IConnection
	// 客户端发来的数据
	Data []byte
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Data
}

func NewRequest(conn ziface.IConnection, data []byte) *Request {
	return &Request{
		Conn: conn,
		Data: data,
	}
}
