package znet

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name string
	IP   string
	Port int
	// 监听的ip类型 "tcp", "tcp4", "tcp6"
	IPVersion   string
	MsgHandler  ziface.IMsgHandler
	ConnManager ziface.IConnectionManager
	OnConnStart func(connection ziface.IConnection)
	OnConnStop  func(connection ziface.IConnection)
}

func (s *Server) Start() {
	addStr := fmt.Sprintf("%s:%d", s.IP, s.Port)
	slog.Info("[Zinx Server]",
		slog.Group("server",
			slog.String("name", s.Name),
			slog.String("addr", addStr),
			slog.String("status", "listening"),
		))
	slog.Info("[Zinx Server]",
		slog.Group("config",
			slog.String("version", utils.GlobalObject.Version),
			slog.Int("maxConn", int(utils.GlobalObject.MaxConn)),
			slog.Int("maxPacketSize", int(utils.GlobalObject.MaxPacketSize)),
		),
	)
	Addr, err := net.ResolveTCPAddr(s.IPVersion, addStr)
	//检查Addr是否有误
	if err != nil {
		slog.Error("[Zinx Server]",
			"ip ver", s.IPVersion,
			slog.String("addr", addStr),
			"msg", "ResolveTCPAddr Error",
			"err", err)
		return
	}
	listener, err := net.ListenTCP(s.IPVersion, Addr)
	if err != nil {
		slog.Error("[Zinx Server]",
			"ip ver", s.IPVersion,
			"addr", addStr,
			"msg", "listen Error",
			"err", err)
		return
	}
	slog.Info(fmt.Sprintf("[%s]", s.Name), "msg", "Zinx Server is running")

	s.MsgHandler.CreatWorkerPool()
	connID := uint32(0)
	for {
		//服务器开始接受连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("[%s]accept err: %v\n", s.Name, err)
			continue
		}

		// 判断连接数是否已达最大
		if s.ConnManager.Len() >= int(utils.GlobalObject.MaxConn) {
			slog.Warn("Server", "MaxConnections", utils.GlobalObject.MaxConn, "msg", "too many connections")
			// TODO 返回服务器超连接信息
			conn.Close()
			continue
		}
		dealConn := NewConnection(s, conn, connID, s.MsgHandler)
		connID++
		dealConn.Start()
	}
}

func (s *Server) GetConnectionManager() ziface.IConnectionManager {
	return s.ConnManager
}

func (s *Server) Stop() {
	slog.Info(fmt.Sprintf("[%s]", s.Name), "msg", "Zinx Server is stopping")
	s.ConnManager.ClearConn()

}

func (s *Server) Init() {
	s.Start()
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

//  设置OnConnStart函数

func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//  设置OnConnStop函数

func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//  调用OnConnStart函数

func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(connection)
	}
}

// 调用OnConnStop函数
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(connection)
	}
}

func NewServer() ziface.IServer {
	return &Server{
		Name:        utils.GlobalObject.Name,
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.Port,
		IPVersion:   "tcp",
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnectionManager(),
	}
}
