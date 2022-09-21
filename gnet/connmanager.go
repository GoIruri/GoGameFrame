package gnet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/giface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	// 管理的连接集合
	connections map[uint32]giface.IConnection
	// 保护连接集合的读写锁
	connLock sync.RWMutex
}

func (c *ConnManager) Add(conn giface.IConnection) {
	// 保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 将conn加入到ConnManager中
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn giface.IConnection) {
	// 保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除连接信息
	delete(c.connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), " remove from ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Get(connID uint32) (giface.IConnection, error) {
	// 保护共享资源map，加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connID not found")
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	// 保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}

	fmt.Println("Clear All connections succ！ conn num = ", c.Len())
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]giface.IConnection),
	}
}
