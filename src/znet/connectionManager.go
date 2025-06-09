package znet

import (
	"errors"
	"log/slog"
	"sync"
	"zinx/ziface"
)

type ConnectionManager struct {
	// 连接map，key为ConnID
	connections map[uint32]ziface.IConnection
	// 连接读写锁
	connLock sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnectionManager) ForEach(callback func(ziface.IConnection)) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	for _, conn := range cm.connections {
		callback(conn)
	}
	return nil, nil
}

// 新增连接
func (cm *ConnectionManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()] = conn
	slog.Debug("ConnMgr",
		"func", "AddConn",
		"ConnID", conn.GetConnID(),
		"msg", "add conn successfully")

}

// 删除连接
func (cm *ConnectionManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
	slog.Debug("ConnMgr",
		"func", "RemoveConn",
		"ConnID", conn.GetConnID(),
		"msg", "remove conn successfully")
}

// 获取连接（通过ConnID）
func (cm *ConnectionManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("Connection not found")
	}
}

// 获取当前连接总数
func (cm *ConnectionManager) Len() int {
	return len(cm.connections)
}

// 清空连接
func (cm *ConnectionManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
}
