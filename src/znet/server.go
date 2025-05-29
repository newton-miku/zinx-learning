package znet

import (
	"errors"
	"fmt"
	"log"
	"net"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IP        string
	Port      int
	IPVersion string
	Router    ziface.IRouter
}

func EchoHandlerFunc(conn *net.TCPConn, data []byte, datalen int) error {
	log.Printf("[%s] recv from client: %s\n", conn.RemoteAddr().String(), string(data[:datalen]))
	_, err := conn.Write(data[:datalen])
	if err != nil {
		log.Printf("[%s]send data err: %v\n", conn.RemoteAddr().String(), err)
		return errors.New("EchoHandlerFunc send data err")
	}
	return nil
}
func (s *Server) Start() {
	log.Printf("[%s] Listening and accepting at %s Port %d\n", s.Name, s.IP, s.Port)
	Addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
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
	log.Printf("[%s]Zinx server start successfully, now listening at %s Port %d\n", s.Name, s.IP, s.Port)

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
	log.Printf("[%s] add router successfully\n", s.Name)
}
func NewServer(name string, ip string, port int) ziface.IServer {
	return &Server{
		Name:      name,
		IP:        ip,
		Port:      port,
		IPVersion: "tcp",
		Router:    nil,
	}
}
