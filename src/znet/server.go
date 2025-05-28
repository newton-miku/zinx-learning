package znet

import (
	"fmt"
	"io"
	"log"
	"net"
	"zinx/ziface"
)

type Server struct {
	Name      string
	Addr      string
	IP        string
	Port      int
	IPVersion string
}

func (s *Server) Start() {
	log.Printf("[%s] Listening and accepting at %s\n", s.Name, s.Addr)
	_, err := net.ResolveTCPAddr(s.IPVersion, s.Addr)
	//检查Addr是否有误
	if err != nil {
		log.Printf("[%s]start Zinx server failed, ip ver:%s\t addr err: %v\n", s.Name, s.IPVersion, err)
		return
	}
	listener, err := net.Listen(s.IPVersion, s.Addr)
	if err != nil {
		log.Printf("[%s]start Zinx server failed, ip ver:%s\t listen err: %v\n", s.Name, s.IPVersion, err)
		return
	}
	log.Printf("[%s]Zinx server start successfully, now listening at %s\n", s.Name, s.Addr)

	for {
		//服务器开始接受连接
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[%s]accept err: %v\n", s.Name, err)
			continue
		}
		go func() {
			defer conn.Close()
			str := make([]byte, 1024)
			for {
				n, err := conn.Read(str)
				if err != nil && err != io.EOF {
					log.Printf("[%s]Client %s read err: %v\n", s.Name, conn.LocalAddr().String(), err)
					return
				}
				if n == 0 {
					log.Printf("[%s]Client %s read 0 bytes,May Dead\n", s.Name, conn.LocalAddr().String())
					return
				}
				if _, err := conn.Write(str[:n]); err != nil {
					log.Printf("[%s]Client %s write err: %v\n", s.Name, conn.LocalAddr().String(), err)
					return
				}
			}
		}()
	}
}

func (s *Server) Stop() {

}

func (s *Server) Init() {
	s.Start()
}

func NewServer(name string, ip string, port int) ziface.IServer {
	return &Server{
		Name:      name,
		Addr:      net.JoinHostPort(ip, fmt.Sprint(port)),
		IP:        ip,
		Port:      port,
		IPVersion: "tcp",
	}
}
