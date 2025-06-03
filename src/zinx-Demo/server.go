package main

import (
	"fmt"
	"log"
	"zinx/ziface"
	"zinx/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

func (r *HelloRouter) Handle(request ziface.IRequest) {
	log.Printf("[Welcome Server] Handle msg from client Addr=%s msg=%s\n", request.GetConnection().RemoteAddr(), request.GetData())
	err := request.GetConnection().SendMsg(0, []byte("Welcome to Zinx Welcom Server"))
	if err != nil {
		log.Printf("[Welcome Server] Handle Send Error %s\n", err)
	}
}

type EchoRouter struct {
	znet.BaseRouter
}

func (r *EchoRouter) Handle(request ziface.IRequest) {
	log.Printf("[Echo Server] Handle msg from client Addr=%s msg=%s\n", request.GetConnection().RemoteAddr(), request.GetData())
	msg := fmt.Sprintf("Echo Server: %s", request.GetData())
	err := request.GetConnection().SendMsg(1, []byte(msg))
	if err != nil {
		log.Printf("[Echo Server] Handle Send Error %s\n", err)
	}
}

func runServer() {
	//创建一个server服务
	s := znet.NewServer()
	s.AddRouter(0, &HelloRouter{})
	s.AddRouter(1, &EchoRouter{})
	//启动server
	s.Start()
}
