package znet

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//  对应的服务器
	fatherServer ziface.IServer
	conn         *net.TCPConn
	connID       uint32
	isClosed     bool
	//  消息处理器
	MsgHandler ziface.IMsgHandler
	//  发送给客户端的消息的channel
	msgChan chan []byte
	//  退出信号的channel
	ExitChan chan struct{} //  退出信号，使用struct不占内存，效率更高
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		fatherServer: server,
		conn:         conn,
		connID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		msgChan:      make(chan []byte),
		ExitChan:     make(chan struct{}),
	}
	//  添加到连接管理器中
	c.fatherServer.GetConnectionManager().Add(c)
	return c
}
func (c *Connection) StartReader() {
	defer c.Stop()
	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			slog.Error("Conn", "ConnID", c.connID, "msg", "read head err", "err", err)
			break
		}
		msg, err := dp.Unpack(headData)
		if err != nil && err != io.EOF {
			slog.Error("Conn", "ConnID", c.connID, "msg", "unpack err", "err", err)
			break
		}
		if msg.GetDataLen() > 0 {
			data := make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				slog.Error("Conn", "ConnID", c.connID, "msg", "read data err", "err", err)
				break
			}
			msg.SetData(data)
		}
		req := NewRequest(c, *msg.(*Message))
		//  如果启用了工作池机制
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//  将请求发送到工作池中
			c.MsgHandler.SendReqToWorker(req)
		} else {
			//  没有启用工作池，直接处理
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Connection) StartWriter() {
	for {
		select {
		case data := <-c.msgChan:
			_, err := c.conn.Write(data)
			if err != nil {
				slog.Error("Client", "write err", err)
				c.ExitChan <- struct{}{}
				return
			}
		case <-c.ExitChan:
			//  Reader goroutine 已退出，Writer goroutine 退出
			return
		}
	}

}
func (c *Connection) Start() {
	slog.Debug("Conn", "ConnID", c.connID, "msg", "start working")
	//  启动读Goroutine
	go c.StartReader()
	//  启动写Goroutine
	go c.StartWriter()
}
func (c *Connection) Stop() {
	if !c.isClosed {
		slog.Debug("Conn", "ConnID", c.connID, "msg", "connection stoped")
		c.isClosed = true
		c.ExitChan <- struct{}{}
		c.conn.Close()
		//  从连接管理器中移除
		c.fatherServer.GetConnectionManager().Remove(c)
		close(c.ExitChan)
		close(c.msgChan)
	} else {
		slog.Warn("Conn", "ConnID", c.connID, "msg", "connection already stoped")
	}
}
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}
func (c *Connection) GetConnID() uint32 {
	return c.connID
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	//  如果连接已关闭
	if c.isClosed {
		return errors.New("connection is closed")
	}
	//  创建一个数据包
	dp := NewDataPack()
	// 将数据封包
	msgPack, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		return fmt.Errorf("pack msg err:%v", err)
	}
	//  发送数据封包
	c.msgChan <- msgPack
	return nil
}
