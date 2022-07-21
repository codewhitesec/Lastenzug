package tcp

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/c-f/talje/lib/config"
	"github.com/c-f/talje/lib/laeufer"
	"github.com/c-f/talje/lib/protocol"
)

// inspired by https://github.com/jpillora/chisel/blob/92d90be68a989377daf61294ef7458612d10da8c/share/tunnel/tunnel_out_ssh_udp.go#L33
//

type Laeufer struct {
	lock sync.Mutex

	pipe *laeufer.Laeufer

	or config.Config
}

func New(or config.Config) *Laeufer {
	prefix := fmt.Sprintf("%s:%s:%s\t", "tcp", or.CommType.String(), or.Direction.String())
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
func (tcp *Laeufer) Receive(c net.Conn) error {
	defer c.Close()

	// Receive is only allowed if Direction is hieven / pull
	if tcp.or.Direction != config.Receiving {
		log.Println("logic ")
		return laeufer.ErrWrongDirection
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

	return nil
}

func (tcp *Laeufer) Dial() {

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

	tcp.pipe.HandleRead(c, id, "send")
	log.Fatal("You should")

	return nil
}

func connectTo(ctx context.Context, network, addr string) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, network, addr)
}
