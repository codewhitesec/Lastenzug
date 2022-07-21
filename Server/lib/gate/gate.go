package gate

import (
	"io"

	"github.com/c-f/talje/lib/ladung"
	"github.com/c-f/talje/lib/protocol"
)

var (
	counter = 1
	// Poolsize = 500
)

var (
	UUIDGen = func() protocol.ID {
		return protocol.NewUUID()
	}
	DefaultTranslator protocol.Translator = &protocol.Json{}

	IncrementGen = func() protocol.ID {
		counter++
		return protocol.NewInt64ID(counter)
	}

	// DefaultLadung ladung.Ladung = &ladung.ByteLadung{}

	// DefaultKueper
	DefaultGate = Binary()
)

// Gate send/receive packets to/from other Agents
type Gate struct {
	Translator protocol.Translator

	IDGen func() protocol.ID
}

/// Send sends a packet to the writer and returns an error
func (n *Gate) Send(gut ladung.Packet, w io.Writer) (err error) {

	// possibility to Wrap/UnWrap the packet
	// gut := n.Container.Einpacken(gut)

	return n.Translator.WritePacket(gut, w)
}

// Receive returns a packet, which were read from the reader using the translator
func (n *Gate) Receive(r io.Reader) (gut ladung.Packet, err error) {

	gut, err = n.Translator.ReadPacket(r)
	if err != nil {
		return
	}

	// possibility to Wrap/UnWrap the packet
	// gut = n.Container.Auspacken(gut)

	return
}

// Returns the binary gate
func Binary() *Gate {
	return &Gate{
		Translator: &protocol.Binary{},
		IDGen:      IncrementGen,
	}
}

// Returns a JSON Gate
func Json() *Gate {
	return &Gate{
		Translator: &protocol.Json{},
		IDGen:      UUIDGen,
	}
}
