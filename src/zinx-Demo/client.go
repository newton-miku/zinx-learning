package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func runClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
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
