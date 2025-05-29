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
	Name      string
	IP        string
	Port      int
	IPVersion string
	Router    ziface.IRouter
}

func (s *Server) Start() {
	log.Printf("[%s] Listening and accepting at %s Port %d\n", s.Name, s.IP, s.Port)
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
		log.Printf("[%s]start Zinx server failed, ip ver:%s\t addr err: %v\n", s.Name, s.IPVersion, err)
		return
	}
	listener, err := net.ListenTCP(s.IPVersion, Addr)
	if err != nil {
		log.Printf("[%s]start Zinx server failed, ip ver:%s\t listen err: %v\n", s.Name, s.IPVersion, err)
		return
	}
	slog.Info(fmt.Sprintf("[%s]", s.Name), "msg", "Zinx Server is running")

	connID := 0
	for {
		//服务器开始接受连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("[%s]accept err: %v\n", s.Name, err)
			continue
		}
		dealConn := NewConnection(conn, uint(connID), s.Router)
		connID++
		dealConn.Start()
	}
}

func (s *Server) Stop() {

}

func (s *Server) Init() {
	s.Start()
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	// log.Printf("[%s] add router successfully\n", s.Name)
	slog.Debug("Router", "msg", "add router successfully")
}
func NewServer() ziface.IServer {
	return &Server{
		Name:      utils.GlobalObject.Name,
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.Port,
		IPVersion: "tcp",
		Router:    nil,
	}
}
