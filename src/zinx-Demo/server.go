package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"log/slog"
	"time"
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

type SessionVerifyRouter struct {
	znet.BaseRouter
}

func (r *SessionVerifyRouter) Handle(request ziface.IRequest) {
	slog.Debug("SessionVerifyRouter", "msgID", request.GetMsgID(), "ClientAddr", request.GetConnection().RemoteAddr())

	// 获取客户端发来的 sessionID
	clientSessionID := request.GetData()

	// 获取服务端存储的 sessionID
	serverSessionID, err := request.GetConnection().GetProperty("sessionID")
	if err != nil {
		slog.Error("SessionVerify", "msg", "GetProperty error", "err", err)
		return
	}
	// 断言为 string 类型
	serverSessionIDStr, ok := serverSessionID.(string)
	if !ok {
		slog.Error("SessionVerify", "msg", "sessionID type assertion to string failed")
		return
	}
	// 比较 sessionID
	if string(clientSessionID) == serverSessionIDStr {
		slog.Info("SessionVerify", "msg", "Session valid", "sessionID", string(clientSessionID))
	} else {
		slog.Warn("SessionVerify", "msg", "Session invalid", "ConnID", request.GetConnection().GetConnID())
		//如果不正确则踢出连接
		request.GetConnection().Stop()
	}
}
func startSessionVerifyTicker(connMgr ziface.IConnectionManager) {
	ticker := time.NewTicker(10 * time.Second) // 每隔10秒验证一次
	go func() {
		if connMgr.Len() < 1 {
			time.Sleep(1 * time.Second)
		}
		for range ticker.C {
			connMgr.ForEach(func(conn ziface.IConnection) {
				err := conn.SendMsg(3, []byte("Please return your sessionID"))
				if err != nil {
					slog.Error("SessionVerify", "Send verify msg error", "conn", conn.RemoteAddr().String(), "err", err)
				}
			})
		}
	}()
}
func ConnStart(connection ziface.IConnection) {
	slog.Debug("ConnStart",
		"msg", "Conn Start")
	err := connection.SendMsg(101, fmt.Appendf(nil, "[ConnStart]Hello,You are %s", connection.RemoteAddr()))
	if err != nil {
		slog.Error("[ConnStart]",
			"msg", "Conn Start Send Error",
			"err", err)
	}

	//  property用法示例
	//  这里使用sha256对连接的相关信息进行加密，作为sessionID返回给客户端
	startTime := time.Now().Unix()
	connection.SetProperty("startTime", startTime)
	connID := connection.GetConnID()
	hasher := sha256.New()
	sessionStr := fmt.Sprint(startTime, connID, connection.RemoteAddr())
	hasher.Write([]byte(sessionStr))
	sessionID := fmt.Sprintf("%x", hasher.Sum(nil))
	connection.SetProperty("sessionID", sessionID)
	err = connection.SendMsg(102, []byte(sessionID))
	if err != nil {
		slog.Error("ConnStart",
			"msg", "Session SendMsg err",
			"error", err)
	}
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
	s.AddRouter(2, &SessionVerifyRouter{})
	startSessionVerifyTicker(s.GetConnectionManager())
	//启动server
	s.Start()
}
