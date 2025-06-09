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
	IP        string `json:"Host"`
	Port      int
	sessionID string
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
	//使用协程读取服务器返回的消息，避免影响发送数据流程
	go func() {
		dp := znet.NewDataPack()
		header := make([]byte, dp.GetHeadLen())
		for {
			_, err := io.ReadFull(conn, header)
			if err != nil {
				slog.Error("Client", "read err", err)
				os.Exit(1)
			}
			msgHead, err := dp.Unpack(header)
			if err != nil {
				slog.Error("Client", "unpack err", err)
				os.Exit(1)
			}
			if msgHead.GetDataLen() > 0 {
				data := make([]byte, msgHead.GetDataLen())
				_, err := io.ReadFull(conn, data)
				if err != nil {
					slog.Error("Client", "read err", err)
					os.Exit(1)
				}
				slog.Info("Client", "msgID", msgHead.GetMsgID(), "msgInfo", data)
				if msgHead.GetMsgID() == 102 {
					ServAdd.sessionID = string(data)
				}
			}
			//当收到id==0的hello消息时，将msgType改为1
			if msgHead.GetMsgID() == 0 {
				msgType = 1
			} else if msgHead.GetMsgID() == 3 {
				// 收到验证请求，回传 sessionID
				dp := znet.NewDataPack()
				msg, _ := dp.Pack(znet.NewMessage(2, []byte(ServAdd.sessionID)))
				_, err = conn.Write(msg)
				if err != nil {
					slog.Error("Client", "write session verify err", err)
					os.Exit(1)
				}
			}
		}
	}()
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
		time.Sleep(1 * time.Second)
	}
}
