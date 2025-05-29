package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type ServerAddr struct {
	IP   string `json:"Host"`
	Port int
}

var ServAdd ServerAddr

func init() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		log.Fatalf("[%s]ReadFile err: %v\n", "Client", err)
	}
	err = json.Unmarshal(data, &ServAdd)
	if err != nil {
		log.Fatalf("[%s]Unmarshal err: %v\n", "Client", err)
	}
	log.Printf("[%s]ServerAddr: %s:%d\n", "Client", ServAdd.IP, ServAdd.Port)
}

func runClient() {
	// 等待2秒，避免服务器未启动完成
	time.Sleep(time.Second * 2)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServAdd.IP, ServAdd.Port))
	if err != nil {
		log.Fatalf("[%s]Client dial err: %v\n", "Client", err)
	}
	log.Printf("[%s]Client start\n", "Client")
	defer conn.Close()
	for {
		sft := time.Now().Format("2006-01-02 15:04:05")
		_, err := conn.Write([]byte(string(sft) + "\tHello Zinx V0.1 , learning by newton_miku\n"))
		if err != nil {
			log.Fatalf("[%s]Client write err: %v\n", "Client", err)
		}
		go io.Copy(os.Stdout, conn)
		time.Sleep(1 * time.Second)
	}
}
