package main

import "zinx/znet"

func main() {
	//创建一个server服务
	s := znet.NewServer("[zinx v0.1]", "127.0.0.1", 8999)
	//启动server
	s.Start()
}
