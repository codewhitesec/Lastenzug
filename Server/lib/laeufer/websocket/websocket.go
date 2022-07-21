package websocket

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/c-f/talje/lib/config"
	"github.com/c-f/talje/lib/laeufer"
	"github.com/c-f/talje/lib/protocol"
	"github.com/gorilla/websocket"
)

var (
	// for websockets
	DefaultUpgrader = websocket.Upgrader{
		ReadBufferSize:  1000000, 
		WriteBufferSize: 1000000,
	}
	DefaultWebDialer = websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
)

type Laeufer struct {
	lock sync.Mutex

	pipe *laeufer.Laeufer

	or config.Config
}

func New(or config.Config) *Laeufer {
	prefix := fmt.Sprintf("%s:%s:%s\t", "websocket", or.CommType.String(), or.Direction.String())
	logger := log.New(log.Default().Writer(), prefix, log.LstdFlags)
	l := &Laeufer{
		pipe: laeufer.NewFromConfig(logger, or),
		or:   or,
	}

	go l.pipe.HandleWrite()
	return l
}

// Receive will be invoked by some other handler
// will attach the new connection and forward them
func (tcp *Laeufer) Receive(w http.ResponseWriter, r *http.Request) {
	ws, err := DefaultUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error: %s", err)
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	defer ws.Close()

	// Since Receive is called multiple times we need to generate new ids

	// id := ladung.NewInt64ID(1337)
	// TODO(author): possibility to release ID

	c := NewWebSocketConn(ws)

	// Receive is only allowed if Direction is hieven / pull
	if tcp.or.Direction != config.Receiving {
		log.Println("logic ")
		// return utils.ErrWrongDirection
		return
	}

	// Since Receive is called multiple times we need to generate new ids
	// id := kueper.DefaultKueper.IDGen()
	id := string(tcp.pipe.GetID().ID())
	log.Println("Register", id)

	// TODO(author): possibility to release ID

	tcp.pipe.Register(id, c)
	// remove connections if it cannot be re
	defer tcp.pipe.Done(id)

	tcp.pipe.HandleRead(c, id, "receive")

	return

}

func (tcp *Laeufer) Send(addr string) (err error) {
	// Send is only allowed if Direction is fieren / push
	if tcp.or.Direction != config.Sending {
		return laeufer.ErrWrongDirection
	}

	id := string(protocol.NewStaticID(laeufer.AgentConnID).ID())
	log.Println("Sending for id", id)
	for {
		tcp.send(addr, id)
	}
}

func (tcp *Laeufer) send(addr string, id string) (err error) {
	c, err := tcp.pipe.HandleDial(id, "tcp", addr, laeufer.DialWithRetry(connectTo))
	if err != nil {
		return err
	}
	defer c.Close()
	defer tcp.pipe.Done(id)

	log.Println("lelllele", id)
	tcp.pipe.HandleRead(c, id, "send")

	log.Fatal("You should")

	return nil
}

func connectTo(ctx context.Context, network, addr string) (net.Conn, error) {
	ws, res, err := DefaultWebDialer.Dial(addr, nil)
	if err != nil {
		log.Println("[websocket] Cannot establish Websocket", err)
		if res != nil {
			log.Println("[websocket] ", res.Status)
		}
		return nil, err
	}

	c := NewWebSocketConn(ws)

	return c, err
}
