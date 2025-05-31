package znet

import "zinx/ziface"

type Request struct {
	// 与客户端建立的连接
	Conn ziface.IConnection
	// 客户端发来的数据
	msg Message
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

func NewRequest(conn ziface.IConnection, message Message) *Request {
	return &Request{
		Conn: conn,
		msg:  message,
	}
}
