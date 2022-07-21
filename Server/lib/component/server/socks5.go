package server

import (
	"log"
	"net"
)

//
type Socks5 struct {
	addr      string
	HandleCon ConnHandleFunc

	l net.Listener
}

// NewSocksServer returns a new SocksServer
func NewSocks5(addr string) *Socks5 {
	return &Socks5{
		HandleCon: DefaultConnHandler,
		addr:      addr,
	}
}

// Bind binds the handleCon
func (s *Socks5) Bind(handler ConnHandleFunc) {
	s.HandleCon = handler
}

// Start starts the socks server
func (s *Socks5) Start() error {
	log.Println("Start Socks server", s.addr)

	return s.ListenAndServe()
}

// ListenAndServe starts a new tcp listener and waiting for connections
func (s *Socks5) ListenAndServe() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.ServeForever(l)
}

// ServeForever is used to serve connections from a listener to the handler
func (s *Socks5) ServeForever(l net.Listener) (err error) {
	s.l = l

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		log.Println("Incoming Connection")
		go s.HandleCon(conn)
	}
}

// Close closes the listener
func (s *Socks5) Close() error {
	return s.l.Close()
}
