package main

import (
	"log"
	"zinx/ziface"
	"zinx/znet"
)

type WelcomRouter struct {
	znet.BaseRouter
}

func (r *WelcomRouter) Handle(request ziface.IRequest) {
	log.Printf("[Welcome Server] Handle msg from client Addr=%s msg=%s\n", request.GetConnection().RemoteAddr(), request.GetData())
	err := request.GetConnection().SendMsg(0, []byte("Welcome to Zinx Welcom Server"))
	if err != nil {
		log.Printf("[Welcome Server] Handle Send Error %s\n", err)
	}
}

func runServer() {
	//创建一个server服务
	s := znet.NewServer()
	s.AddRouter(&WelcomRouter{})
	//启动server
	s.Start()
}
