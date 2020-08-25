// Code generated by "genprotocol.exe -ver=a99e26984c4d0465e81623a4767d3a2aa3cb4fcc9890904054c0c51f30e0b79f -basedir=protocol_w3d -prefix=w3d -statstype=int"

package w3d_connbytemanager

import (
	"fmt"
	"sync"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_serveconnbyte"
)

type Manager struct {
	mutex   sync.RWMutex
	id2Conn map[string]*w3d_serveconnbyte.ServeConnByte
}

func New() *Manager {
	rtn := &Manager{
		id2Conn: make(map[string]*w3d_serveconnbyte.ServeConnByte),
	}
	return rtn
}
func (cm *Manager) Add(id string, c2sc *w3d_serveconnbyte.ServeConnByte) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if cm.id2Conn[id] != nil {
		return fmt.Errorf("already exist %v", id)
	}
	cm.id2Conn[id] = c2sc
	return nil
}
func (cm *Manager) Del(id string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if cm.id2Conn[id] == nil {
		return fmt.Errorf("not exist %v", id)
	}
	delete(cm.id2Conn, id)
	return nil
}
func (cm *Manager) Get(id string) *w3d_serveconnbyte.ServeConnByte {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.id2Conn[id]
}
func (cm *Manager) Len() int {
	return len(cm.id2Conn)
}
func (cm *Manager) GetList() []*w3d_serveconnbyte.ServeConnByte {
	rtn := make([]*w3d_serveconnbyte.ServeConnByte, 0, len(cm.id2Conn))
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	for _, v := range cm.id2Conn {
		rtn = append(rtn, v)
	}
	return rtn
}
