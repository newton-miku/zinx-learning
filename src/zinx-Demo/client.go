package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
	"zinx/znet"
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
}

func runClient() {
	// 等待2秒，避免服务器未启动完成
	time.Sleep(time.Second * 2)
	log.Printf("[%s]ServerAddr: %s:%d\n", "Client", ServAdd.IP, ServAdd.Port)
	conn, err := net.Dial("tcp", net.JoinHostPort(ServAdd.IP, fmt.Sprint(ServAdd.Port)))
	if err != nil {
		log.Fatalf("[%s]Client dial err: %v\n", "Client", err)
	}
	log.Printf("[%s]Client start\n", "Client")
	defer conn.Close()
	for {
		sft := time.Now().Format("2006-01-02 15:04:05")
		dp := znet.NewDataPack()
		msg, err := dp.Pack(znet.NewMessage(1, []byte(sft)))
		if err != nil {
			log.Fatalf("[%s]Client pack err: %v\n", "Client", err)
		}
		_, err = conn.Write(msg)
		if err != nil {
			log.Fatalf("[%s]Client write err: %v\n", "Client", err)
		}
		go func() {
			header := make([]byte, dp.GetHeadLen())
			_, err := io.ReadFull(conn, header)
			if err != nil {
				log.Fatalf("[%s]Client read err: %v\n", "Client", err)
			}
			msgHead, err := dp.Unpack(header)
			if err != nil {
				log.Fatalf("[%s]Client unpack err: %v\n", "Client", err)
			}
			if msgHead.GetDataLen() > 0 {
				data := make([]byte, msgHead.GetDataLen())
				_, err := io.ReadFull(conn, data)
				if err != nil {
					log.Fatalf("[%s]Client read err: %v\n", "Client", err)
				}
				log.Printf("[%s]Client recv msg ID=%d msgInfo=%s\n", "Client", msgHead.GetMsgID(), string(data))
			}
		}()
		time.Sleep(1 * time.Second)
	}
}
