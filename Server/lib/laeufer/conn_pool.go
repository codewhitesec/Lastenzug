package laeufer

import (
	"net"
	"sync"
)

var (
	AgentConnID string = "Agent2Agent"
)

type ConnPool struct {
	sync.Mutex
	m map[string]net.Conn // contains faser
}

// NewConnPool
func NewConnPool() *ConnPool {
	return &ConnPool{
		m: make(map[string]net.Conn),
	}
}

func (cs *ConnPool) len() int {
	cs.Lock()
	l := len(cs.m)
	cs.Unlock()
	return l
}

func (cs *ConnPool) add(id string, c net.Conn) {
	cs.Lock()
	defer cs.Unlock()

	cs.m[id] = c
}

func (cs *ConnPool) get(id string) (conn net.Conn) {
	cs.Lock()
	defer cs.Unlock()
	if conn, ok := cs.m[id]; ok {
		return conn
	}
	return nil
}

func (cs *ConnPool) remove(id string) {
	cs.Lock()
	delete(cs.m, id)
	cs.Unlock()
}

func (cs *ConnPool) closeAll() {
	cs.Lock()
	for id, conn := range cs.m {
		conn.Close()
		delete(cs.m, id)
	}
	cs.Unlock()
}
