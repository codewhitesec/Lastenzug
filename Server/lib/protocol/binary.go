package protocol

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/c-f/talje/lib/kommando"
	"github.com/c-f/talje/lib/ladung"
)

// Error Messages for handling with byte streams
var (
	ErrMessageToSmall   ErrParsing = errors.New("message to small")
	ErrMessageBadLength ErrParsing = errors.New("bad length")
)

// ByteLadung creates bytestream Ladung
type Binary struct{}

// WritePacket encodes a Packet and writes it to a writer
func (bl *Binary) WritePacket(g ladung.Packet, w io.Writer) (err error) {
	var bts []byte

	session := []byte(g.SessionID)[0:8]

	command := make([]byte, 8)
	binary.LittleEndian.PutUint64(command, uint64(g.Kommando))

	l := make([]byte, 8)
	binary.LittleEndian.PutUint64(l, uint64(len(g.Message)))

	// Create Protocol
	// [ sess | komando | len(msg) | msg ]
	bts = append(bts, session...)
	bts = append(bts, command...)
	bts = append(bts, l...)
	bts = append(bts, g.Message...)

	_, err = w.Write(bts)
	return
}

// ReadPacket decodes a Packet from a reader
func (bl *Binary) ReadPacket(r io.Reader) (g ladung.Packet, err error) {
	var n int

	session := make([]byte, 8)
	command := make([]byte, 8)
	l := make([]byte, 8)

	readVar := func(val []byte) error {
		n, err = r.Read(val)
		if err == nil && n != len(val) {
			err = ErrMessageToSmall
		}
		return err
	}

	// Read Variables
	if err = readVar(session); err != nil {
		return
	}
	if err = readVar(command); err != nil {
		return
	}
	if err = readVar(l); err != nil {
		return
	}

	length := int(binary.LittleEndian.Uint64(l))
	if length < 0 || length >= 100000 {
		err = ErrMessageBadLength
		return
	}
	msg := make([]byte, length)
	if err = readVar(msg); err != nil {
		return
	}

	// Stitch together
	g.SessionID = string(session)
	g.Kommando = kommando.Code(BytesToInt(command))
	g.Message = msg

	return
}

// BytesToInt converts byte to int
func BytesToInt(bts []byte) int {
	return int(binary.LittleEndian.Uint64(bts))
}
