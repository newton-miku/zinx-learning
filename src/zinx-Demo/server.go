package main

import (
	"log"
	"zinx/ziface"
	"zinx/znet"
)

type EchoRouter struct {
	znet.BaseRouter
}

func (r *EchoRouter) PreHandle(request ziface.IRequest) {
	log.Printf("[Echo] PreHandle %s\n", request.GetData())
	err := request.GetConnection().Send([]byte("Welcome to Zinx v0.3\n"))
	if err != nil {
		log.Printf("[Echo] PreHandle Send Error %s\n", err)
	}
}
func (r *EchoRouter) Handle(request ziface.IRequest) {
	log.Printf("[Echo] Handle %s\n", request.GetData())
	msg := []byte("You says:")
	msg = append(msg, request.GetData()...)
	err := request.GetConnection().Send(msg)
	if err != nil {
		log.Printf("[Echo] Handle Send Error %s\n", err)
	}
}

func (r *EchoRouter) PostHandle(request ziface.IRequest) {
	log.Println("[Echo] PostHandle")
	err := request.GetConnection().Send([]byte("See You Next Time!\n"))
	if err != nil {
		log.Printf("[Echo] PostHandle Send Error %s\n", err)
	}
}

func runServer() {
	//创建一个server服务
	s := znet.NewServer("[zinx v0.3]", "127.0.0.1", 8999)
	s.AddRouter(&EchoRouter{})
	//启动server
	s.Start()
}
