package laeufer

import (
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/armon/go-socks5"
	"github.com/c-f/talje/lib/config"
	"github.com/c-f/talje/lib/gate"
	"github.com/c-f/talje/lib/kommando"
	"github.com/c-f/talje/lib/ladung"
	"github.com/c-f/talje/lib/protocol"
)

var (
	ErrNoConn = errors.New("No connection available")
)

// Laeufer is the glue code, which handles the internal connection management
// Reading and Writing to the Packet data streams based on the configuration
type Laeufer struct {
	// conns contains all Connections
	conns *ConnPool

	// Data Streams
	inC  chan ladung.Packet
	outC chan ladung.Packet

	// Type/Direction of communication
	commType      config.Communication
	commDirection config.Direction

	log *log.Logger
}

// NewFromConfig creates a new Laeufer
func NewFromConfig(l *log.Logger, or config.Config) *Laeufer {
	return New(l, or.InC, or.OutC, or.CommType, or.Direction)
}

// New creates a new Laeufer
func New(logger *log.Logger, inc chan ladung.Packet, outc chan ladung.Packet, commType config.Communication, direction config.Direction) *Laeufer {
	lfr := &Laeufer{
		log:  logger,
		inC:  inc,
		outC: outc,

		commType:      commType,
		commDirection: direction,

		conns: NewConnPool(),
	}

	return lfr
}

// HandleDial can be called from sending agents to establish a connection to an other agent (receiving)
func (lfr *Laeufer) HandleDial(id string, network, addr string, dial DialFunc) (c net.Conn, err error) {

	return lfr.handleDial(id, network, addr, dial)
}

func (lfr *Laeufer) handleDial(id string, network, addr string, dial DialFunc) (c net.Conn, err error) {

	c, err = lfr.GetConn(id, false)
	if err != nil {
		// if no connection is known then connect

		if err == ErrNoConn {
			c, err = dial(network, addr)
			if err != nil {
				log.Println("dialing error", err)
				return
			}
			lfr.Register(id, c)

		} else {
			lfr.log.Println("handling Error", err)
			return
		}
	}

	return
}

// HandleRead is responsible to read Packets from
//
func (lfr *Laeufer) HandleRead(conn net.Conn, id string, label string) (err error) {

	// Temporary buffer for direct reads
	tmp := make([]byte, 100000) // using small tmo buffer for demonstrating

	var isEoF = false

	var n int

	// Handle Network Connection
	for {

		var g ladung.Packet

		// Read Packet
		// based on commType decoding is necessary thus using the Gate
		switch lfr.commType {
		case config.Agent2Agent:
			g, err = gate.DefaultGate.Receive(conn)

			if err != nil {
				lfr.log.Println("[lfr] entladen problems :/", err)
				return
			}
			lfr.log.Println("received", g.SessionID, len(g.Message))

		case config.Forward:
			panic("not implemented yet")

		case config.Direct:
			lfr.log.Println("received", g.SessionID, len(g.Message))

			// read Directly from connection
			n, err = conn.Read(tmp)

			payload := make([]byte, n)
			copy(payload, tmp)

			// Create Packages
			if err != nil {
				n = 0
				if err == io.EOF {
					lfr.log.Println("finished")
				} else {
					lfr.log.Println("reading problem", err)
				}
				g = ladung.StopPacket(id, ladung.EmptyMessage)

				isEoF = true

			} else {
				// create DataPacket
				g = ladung.DataPacket(id, payload)
			}
		}

		// Send Packet to channel
		lfr.outC <- g

		// Stopp if everything is ok
		if isEoF {
			lfr.log.Println("End of File - stopp reading ")
			return
		}
	}
}

// HandleWrite reads incoming Packets and writes them to the specified connection
func (lfr *Laeufer) HandleWrite() {
	for gut := range lfr.inC {
		id := gut.SessionID

		if lfr.commType == config.Direct {
			lfr.log.Println("send", gut.SessionID, len(gut.Message))
		}

		// Currently if the client is not reachable it will block the complete queue
		// might be possible for direct connections to send to channel instead with go routine

		// replace ID
		if lfr.commType == config.Agent2Agent {
			id = AgentConnID
		}

		// wait that client is active
		c, err := lfr.GetConn(id, true)
		if err != nil {

			// if no Connection is available and the laeufer is configured as sender
			// try to establish a connection
			if lfr.commDirection == config.Sending {

				c1, c2 := net.Pipe()
				c = c1

				lfr.conns.add(id, c)

				// Handle Reading of the connection
				go func() {
					err = lfr.HandleRead(c, id, "lel")
					lfr.log.Println("Error", err)
				}()

				// handle Read/Write to upstream
				go func() {
					// // Possible hooking of dialer Function possible
					// dialfunc := func(ctx xcontext.Context, network, addr string) (net.Conn, error) {
					// 	var d net.Dialer
					// 	log.Println(":)")
					// 	underlyingConn, err := d.DialContext(ctx, network, addr)
					// 	//return utils.NewBuffConn(underlyingConn, c), err
					// 	return underlyingConn, err
					// }
					log.Println("Try to establish connection")

					// --[Socks5 Client]--
					soccksSrv, _ := socks5.New(&socks5.Config{})
					err := soccksSrv.ServeConn(c2)

					// --[Other options]-- HTTP

					// --[Other options]-- TCP


					// Error handling
					if err != nil {
						// Server said something
						lfr.log.Println("[-] Socks err", id, err)
					}

					c2.Close()
					lfr.log.Println("[-] Session terminated", id)

					// remove conns
					lfr.conns.remove(id)
				}()

			} else {
				lfr.log.Println("connection stopped ")
				lfr.log.Println("ignore package")
				continue
			}

		}

		// Handle Packet Message
		switch gut.Kommando {
		case kommando.Stop:

			// close connection if it's not an agent2agent connection
			if lfr.commType != config.Agent2Agent {
				lfr.log.Println("cmd:stop", gut.SessionID, len(gut.Message))

				c.Close()
			} else {
				// forward message
				err = gate.DefaultGate.Send(gut, c)
				if err != nil {
					lfr.log.Println("gate encoding problems :/", err)
				}
			}

		case kommando.Data:
			lfr.log.Println("send", gut.SessionID, len(gut.Message))

			// Send to the connection
			switch lfr.commType {
			case config.Agent2Agent:

				err = gate.DefaultGate.Send(gut, c)
				if err != nil {
					lfr.log.Println("gate encoding problems :/", err)
				}

			case config.Forward:
				panic("not implemented yet")

			case config.Direct:

				_, err := c.Write(gut.Message)
				if err != nil {
					lfr.log.Println("[!]", "error writing to lose connection", err)
				}
			}
		default:
			lfr.log.Println("[!]", "cmd:unknown", gut.Kommando)
		}
	}
}

// --[Conn handling]--

// Register a new Connection
func (lfr *Laeufer) Register(id string, c net.Conn) {
	lfr.log.Println("[register]", id)
	lfr.conns.add(id, c)
}

// Delete a new Connection
func (lfr *Laeufer) Done(id string) {
	lfr.log.Println("[done]", id)
	lfr.conns.remove(id)
}

// Get a Connection, potentically wait
func (lfr *Laeufer) GetConn(oriID string, shouldWait bool) (c net.Conn, err error) {
	var id string

	switch lfr.commType {
	case config.Agent2Agent:
		id = AgentConnID
	case config.Forward:
		id = AgentConnID
	case config.Direct:
		id = oriID
	}

	// timeout
	attempts := 0
	for {
		c = lfr.conns.get(id)

		// connection was found return it
		if c != nil {
			return c, nil
		}

		// connection will never be reached since the richtung is fier / push
		if lfr.commType == config.Direct {
			err = ErrNoConn
			return
		}

		if !shouldWait {
			err = ErrNoConn
			return
		}

		// Newly added
		if lfr.commDirection == config.Sending {
			err = ErrNoConn
			return
		}

		lfr.log.Println("got message but socket is not available", oriID)
		time.Sleep(1 * time.Second)
		attempts++

		if attempts >= 10 {
			lfr.log.Println("Discard message. Client needs to resend it ", oriID)
      err = ErrNoConn
      return
		}
	}
}

// New Id
func (lfr *Laeufer) GetID() (id protocol.ID) {
	switch lfr.commDirection {
	case config.Sending: // push
		id = protocol.NewStaticID(AgentConnID)

	case config.Receiving: // pull
		id = gate.DefaultGate.IDGen()

	default:
		panic("unknown direction")
	}

	switch lfr.commType {
	case config.Agent2Agent:
		id = protocol.NewStaticID(AgentConnID)
	}

	return
}
