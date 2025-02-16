package znet

import (
	"fmt"
	"go_zinx/ziface"
	"sync"
)

type ConnManager struct {
	Connections map[uint32]ziface.IConnection // 管理连接
	ConnLock    sync.RWMutex                  // 读写锁
}

func NewConnManager() ziface.IConnectionManager {
	return &ConnManager{
		Connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	cm.Connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.ConnLock.Lock()

	defer cm.ConnLock.Unlock()

	delete(cm.Connections, conn.GetConnID())
	fmt.Println("connection delete from ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.ConnLock.RLock()
	defer cm.ConnLock.RUnlock()

	if conn, ok := cm.Connections[connID]; !ok {
		return nil, fmt.Errorf("connection not found")
	} else {
		return conn, nil
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.Connections)
}

func (cm *ConnManager) ClearConn() {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	for connID, conn := range cm.Connections {
		fmt.Println("ClearConn, connID = ", connID)
		conn.Stop()
		delete(cm.Connections, connID)
	}
	fmt.Println("clear all connections successfully: conn num = ", cm.Len())
}
