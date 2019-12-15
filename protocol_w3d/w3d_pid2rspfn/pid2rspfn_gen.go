// Code generated by "genprotocol -ver=f5b1d289172cf84ad5d01b91533408be6b17961cf28ddd6fe767224298a8aedd -basedir=. -prefix=w3d -statstype=int"

package w3d_pid2rspfn

import (
	"fmt"
	"sync"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
)

type HandleRspFn func(w3d_packet.Header, interface{}) error
type PID2RspFn struct {
	mutex      sync.Mutex
	pid2recvfn map[uint32]HandleRspFn
	pid        uint32
}

func New() *PID2RspFn {
	rtn := &PID2RspFn{
		pid2recvfn: make(map[uint32]HandleRspFn),
	}
	return rtn
}
func (p2r *PID2RspFn) NewPID(fn HandleRspFn) uint32 {
	p2r.mutex.Lock()
	defer p2r.mutex.Unlock()
	p2r.pid++
	p2r.pid2recvfn[p2r.pid] = fn
	return p2r.pid
}
func (p2r *PID2RspFn) HandleRsp(header w3d_packet.Header, body interface{}) error {
	p2r.mutex.Lock()
	if recvfn, exist := p2r.pid2recvfn[header.ID]; exist {
		delete(p2r.pid2recvfn, header.ID)
		p2r.mutex.Unlock()
		return recvfn(header, body)
	}
	p2r.mutex.Unlock()
	return fmt.Errorf("pid not found")
}
