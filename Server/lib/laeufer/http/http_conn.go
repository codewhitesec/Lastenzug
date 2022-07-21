package http

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type HTTPContent []byte

type HTTPConn struct {
	msgOut chan []byte

	//outqueue []HTTPContent

	rw io.ReadWriter

	addr string

	cond *sync.Cond
}

func NewHTTPConn() *HTTPConn {
	m := sync.Mutex{}
	return &HTTPConn{
		msgOut: make(chan []byte),
		rw:     bytes.NewBuffer([]byte{}),
		addr:   "empty",
		cond:   sync.NewCond(&m),
	}
}

func (c *HTTPConn) Polling(addr string) {
	c.addr = addr
	// ticker := time.NewTicker(time.Second * 10)

	//queue := make([]byte, 0)
	for {
		var msg []byte

		log.Println("Polling ....")

		// wait either for new or when the timer ticks
		select {
		case <-time.After(1 * time.Second):
			msg = []byte("")
			// ping msg

			log.Println("PingPoll ....")

		case msg = <-c.msgOut:
			log.Println("Message ....")

		}

		//
		res, err := c.Send(msg)

		if err != nil {
			log.Println("lel")
			continue
		}

		if len(res) > 0 {
			log.Println("Also Write to stuff")
			c.rw.Write(res)
			c.cond.Broadcast()
		}

	}
}

func (c *HTTPConn) AddMessage(bts []byte) {
	c.rw.Write(bts)
	c.cond.Broadcast()
}

// TODO(author): channel or queue ?
func (c *HTTPConn) GetMessage() []byte {
	select {
	case msg := <-c.msgOut:
		return msg
	default:
		return []byte{}
	}
}

func (c *HTTPConn) Read(bts []byte) (n int, err error) {
	n, err = c.rw.Read(bts)
	if err == io.EOF {
		c.cond.L.Lock()
		c.cond.Wait()
		c.cond.L.Unlock()
		n, err = c.rw.Read(bts)
	}
	log.Println("Reading ...", n, err)

	return
}

func (c *HTTPConn) Write(b []byte) (n int, err error) {
	c.msgOut <- b
	return len(b), nil
}

type HTTPAddr string

func (a HTTPAddr) Network() string {
	return "tcp"
}
func (a HTTPAddr) String() string {
	return string(a)
}

func (c *HTTPConn) Close() error {
	return nil
}

func (c *HTTPConn) LocalAddr() net.Addr {
	return HTTPAddr(c.addr)
}

func (c *HTTPConn) RemoteAddr() net.Addr {
	return HTTPAddr(c.addr)
}

//

func (c *HTTPConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *HTTPConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *HTTPConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *HTTPConn) Send(msg []byte) (bts []byte, err error) {
	buf := bytes.NewBuffer(msg)
	resp, err := http.Post(c.addr, "image/jpeg", buf)
	if err != nil {
		log.Println("err", err)
		return
	}
	defer resp.Body.Close()

	bts, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("errr", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("invalid message")
	}
	return
}
