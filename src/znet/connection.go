package znet

import (
	"io"
	"log"
	"net"
	"zinx/ziface"
)

type Connection struct {
	conn     *net.TCPConn
	connID   uint
	isClosed bool
	//当前连接处理的方法Router
	Router ziface.IRouter
	//退出信号的channel
	ExitChan chan struct{} //退出信号，使用struct不占内存，效率更高
}

func NewConnection(conn *net.TCPConn, connID uint, router ziface.IRouter) *Connection {
	return &Connection{
		conn:     conn,
		connID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan struct{}),
	}
}
func (c *Connection) StartReader() {
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		n, err := c.conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("[Conn %d]\tread err: %v\n", c.connID, err)
			return
		}

		req := NewRequest(c, buf[:n])
		c.Router.PreHandle(req)
		c.Router.Handle(req)
		c.Router.PostHandle(req)
	}
}
func (c *Connection) StartWriter() {

}
func (c *Connection) Start() {
	log.Printf("[Conn %d]\tconnection start...\n", c.connID)
	// 启动读Goroutine
	go c.StartReader()
	// 启动写Goroutine
	go c.StartWriter()
}
func (c *Connection) Stop() {
	if !c.isClosed {
		log.Printf("[Conn %d]\tconnection stop...\n", c.connID)
		c.isClosed = true
		c.conn.Close()
		close(c.ExitChan)
	} else {
		log.Printf("[Conn %d]!!!\tconnection already stop...\n", c.connID)
	}
}
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}
func (c *Connection) GetConnID() uint {
	return c.connID
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
func (c *Connection) Send(data []byte) error {
	if c.isClosed {
		return nil
	}
	_, err := c.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}
