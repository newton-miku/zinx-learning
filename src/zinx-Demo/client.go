package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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
		slog.Error("Client", "ReadFile err", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &ServAdd)
	if err != nil {
		slog.Error("Client", "Unmarshal err", err)
		os.Exit(1)
	}
}

func runClient() {
	// 等待2秒，避免服务器未启动完成
	time.Sleep(time.Second * 2)
	slog.Info("[Client]", "ServerAddr", ServAdd.IP, "Port", ServAdd.Port)
	conn, err := net.Dial("tcp", net.JoinHostPort(ServAdd.IP, fmt.Sprint(ServAdd.Port)))
	if err != nil {
		slog.Error("Client", "dial err", err)
		return
	}
	slog.Info("Client", "msg", "started")
	defer conn.Close()
	var msgType uint32 = 0
	for {
		//获取当前时间并格式化
		sft := time.Now().Format("2006-01-02 15:04:05")
		//创建数据包工具
		dp := znet.NewDataPack()
		//信息封包，内容为当前时间
		msg, err := dp.Pack(znet.NewMessage(msgType, []byte(sft)))
		if err != nil {
			slog.Error("Client", "pack err", err)
			return
		}
		_, err = conn.Write(msg)
		if err != nil {
			slog.Error("Client", "write err", err)
			return
		}
		//使用协程读取服务器返回的消息，避免影响发送数据流程
		go func() {
			header := make([]byte, dp.GetHeadLen())
			_, err := io.ReadFull(conn, header)
			if err != nil {
				slog.Error("Client", "read err", err)
				return
			}
			msgHead, err := dp.Unpack(header)
			if err != nil {
				slog.Error("Client", "unpack err", err)
				return
			}
			if msgHead.GetDataLen() > 0 {
				data := make([]byte, msgHead.GetDataLen())
				_, err := io.ReadFull(conn, data)
				if err != nil {
					slog.Error("Client", "read err", err)
					return
				}
				slog.Info("Client", "msgID", msgHead.GetMsgID(), "msgInfo", string(data))
			}
			//当收到id==0的hello消息时，将msgType改为1
			if msgHead.GetMsgID() == 0 {
				msgType = 1
			}
		}()
		time.Sleep(1 * time.Second)
	}
}
