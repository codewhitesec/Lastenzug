package protocol

import (
	"io"

	"github.com/c-f/talje/lib/ladung"
)

// Translater defines how to read/write Packets between agents
type Translator interface {
	// Einpacken might also be used with PacketHandle
	WritePacket(g ladung.Packet, w io.Writer) (err error)
	ReadPacket(r io.Reader) (g ladung.Packet, err error)
}

// defines the ID interface
type ID interface {
	ID() string
}

// ErrParsing defines a protocol-specific Errors
type ErrParsing error
