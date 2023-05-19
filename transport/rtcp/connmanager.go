package rtcp

import (
	"errors"
	"sync"

	"github.com/cr-mao/lori/log"
	"github.com/cr-mao/lori/transport/rtcp/iface"
)

type ConnManager struct {
	connections map[uint64]iface.IConnection
	connLock    sync.RWMutex
}

func newConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint64]iface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn iface.IConnection) {

	connMgr.connLock.Lock()
	connMgr.connections[conn.GetConnID()] = conn //将conn连接添加到ConnMananger中
	connMgr.connLock.Unlock()

	log.Infof("connection add to ConnManager successfully: conn num = %d", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn iface.IConnection) {

	connMgr.connLock.Lock()
	delete(connMgr.connections, conn.GetConnID()) //删除连接信息
	connMgr.connLock.Unlock()

	log.Infof("connection Remove ConnID=%d successfully: conn num = %d", conn.GetConnID(), connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint64) (iface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

func (connMgr *ConnManager) Len() int {

	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()

	return length
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	// Stop and delete all connection information
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	connMgr.connLock.Unlock()

	log.Infof("Clear All Connections successfully: conn num = %d", connMgr.Len())
}

func (connMgr *ConnManager) GetAllConnID() []uint64 {
	ids := make([]uint64, 0, len(connMgr.connections))

	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	for id := range connMgr.connections {
		ids = append(ids, id)
	}

	return ids
}

func (connMgr *ConnManager) Range(cb func(uint64, iface.IConnection, interface{}) error, args interface{}) (err error) {

	for connID, conn := range connMgr.connections {
		err = cb(connID, conn, args)
	}

	return err
}
