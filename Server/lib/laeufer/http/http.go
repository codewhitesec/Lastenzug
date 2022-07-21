package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/c-f/talje/lib/config"
	"github.com/c-f/talje/lib/laeufer"
	"github.com/c-f/talje/lib/protocol"
)

//TODO(author): implement http handler

type Laeufer struct {
	pipe *laeufer.Laeufer

	or config.Config

	fakeHTTP *HTTPConn
	once     sync.Once
}

func New(or config.Config) *Laeufer {
	prefix := fmt.Sprintf("%s:%s:%s\t", "websocket", or.CommType.String(), or.Direction.String())
	logger := log.New(log.Default().Writer(), prefix, log.LstdFlags)
	l := &Laeufer{
		pipe: laeufer.NewFromConfig(logger, or),
		or:   or,
	}

	l.fakeHTTP = NewHTTPConn()

	go l.pipe.HandleWrite()
	return l
}

// Receive will be invoked by some other handler
// will attach the new connection and forward them
func (tcp *Laeufer) Receive(w http.ResponseWriter, r *http.Request) {

	// Receive is only allowed if Richtung is hieven / pull
	if tcp.or.Direction != config.Receiving {
		log.Println("logic ")
		// return utils.ErrWrongRichtung
		return
	}

	// Since Receive is called multiple times we need to generate new ids
	// id := kueper.DefaultKueper.IDGen()
	id := string(tcp.pipe.GetID().ID())
	log.Println("Register", id)
	// TODO(author): possibility to release ID

	// Start Routine once
	go tcp.once.Do(func() {
		tcp.pipe.Register(id, tcp.fakeHTTP)
		// remove connections if it cannot be re
		defer tcp.pipe.Done(id)

		tcp.pipe.HandleRead(tcp.fakeHTTP, id, "receive")
	})
	log.Println("lelel")

	// now get the infos
	defer r.Body.Close()
	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("luls")
	}

	if len(bts) > 0 {
		// do something with it
		// get pipe write to pipe

		tcp.fakeHTTP.AddMessage(bts)
	}

	// send message themself
	if sendMsg := tcp.fakeHTTP.GetMessage(); len(sendMsg) > 0 {
		_, err := w.Write(sendMsg)
		if err != nil {
			log.Println("Could not write to http response ")
		}
	}

	return

}

func (tcp *Laeufer) Send(addr string) (err error) {
	// Send is only allowed if Richtung is fieren / push
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

	go tcp.fakeHTTP.Polling(addr)

	tcp.pipe.Register(id, tcp.fakeHTTP)
	defer tcp.pipe.Done(id)

	// c, err := tcp.pipe.HandleDial(id, "tcp", addr, laeufer.DialWithRetry(func(ctx context.Context, network, addr string) (net.Conn, error) {
	// 	return tcp.fakeHTTP, nil
	// }))
	if err != nil {
		return err
	}

	log.Println("lelllele", id)
	tcp.pipe.HandleRead(tcp.fakeHTTP, id, "send")

	log.Fatal("You should")

	return nil
}
