package server

import "net"

var (
	emptyHandler = func(c net.Conn) error {
		return c.Close()
	}

	DefaultConnHandler = emptyHandler
)

// ConnHandleFunc
type ConnHandleFunc func(net.Conn) error
