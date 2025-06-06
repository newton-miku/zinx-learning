package main

import (
	"fmt"
	"log"
	"log/slog"
	"zinx/ziface"
	"zinx/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

func (r *HelloRouter) Handle(request ziface.IRequest) {
	slog.Debug("Router",
		"Router Name", "Welcome",
		"ClientAddr", request.GetConnection().RemoteAddr(),
		"msgID", request.GetMsgID(),
		"msg", request.GetData())
	err := request.GetConnection().SendMsg(0, []byte("Welcome to Zinx Welcom Server"))
	if err != nil {
		log.Printf("[Welcome Server] Handle Send Error %s\n", err)
	}
}

type EchoRouter struct {
	znet.BaseRouter
}

func (r *EchoRouter) Handle(request ziface.IRequest) {
	slog.Debug("Router",
		"RouterName", "Echo",
		"ClientAddr", request.GetConnection().RemoteAddr(),
		"msgID", request.GetMsgID(),
		"msg", request.GetData())
	msg := fmt.Sprintf("Echo Server: %s", request.GetData())
	err := request.GetConnection().SendMsg(1, []byte(msg))
	if err != nil {
		slog.Error("[Router]",
			"RouterName", "Echo",
			"msg", "Handle Send Error",
			"err", err)
	}
}
func ConnStart(connection ziface.IConnection) {
	slog.Debug("ConnStart",
		"msg", "Conn Start")
	connection.SendMsg(101, fmt.Appendf(nil, "[ConnStart]Hello,You are %s", connection.RemoteAddr()))
}
func ConnStop(connection ziface.IConnection) {
	slog.Debug("ConnStop",
		"msg", "Conn Stop")
	//  此处通过打印模拟在数据库中设置相关客户端的状态
	fmt.Println("Conn lost,ID:", connection.GetConnID())
}
func runServer() {
	//创建一个server服务
	s := znet.NewServer()
	s.SetOnConnStart(ConnStart)
	s.SetOnConnStop(ConnStop)
	s.AddRouter(0, &HelloRouter{})
	s.AddRouter(1, &EchoRouter{})
	//启动server
	s.Start()
}
