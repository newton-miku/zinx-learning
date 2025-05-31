package znet

import (
	"errors"
	"fmt"
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
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			log.Printf("[Conn %d] read head err: %v\n", c.connID, err)
			break
		}
		msg, err := dp.Unpack(headData)
		if err != nil && err != io.EOF {
			log.Printf("[Conn %d] unpack err: %v\n", c.connID, err)
			break
		}
		if msg.GetDataLen() > 0 {
			data := make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				log.Printf("[Conn %d] read data err: %v\n", c.connID, err)
				break
			}
			msg.SetData(data)
		}
		req := NewRequest(c, *msg.(*Message))
		if c.Router != nil {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		} else {
			log.Printf("[Conn %d] warning: router is nil", c.connID)
		}
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
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	// 如果连接已关闭
	if c.isClosed {
		return errors.New("connection is closed")
	}
	// 创建一个数据包
	dp := NewDataPack()
	// 将数据封包
	msgPack, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		return fmt.Errorf("pack msg err:%v", err)
	}
	// 发送数据封包
	if _, err = c.conn.Write(msgPack); err != nil {
		return fmt.Errorf("write msg err:%v", err)
	}
	return nil
}
