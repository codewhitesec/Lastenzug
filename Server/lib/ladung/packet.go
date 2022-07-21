package ladung

import (
	"github.com/c-f/talje/lib/kommando"
)

// Packet is the struct for transporting between Laeufer and Agent
type Packet struct {
	// SessionID defines messages to a specific session/request
	SessionID string

	// Message defines the payload
	Message []byte

	// Kommando defines what's the purpose of the packet
	Kommando kommando.Kommando
}

var (
	EmptyMessage = []byte{}
)

// NewDataPacket creates a new DataPacket
func DataPacket(id string, msg []byte) Packet {
	return new(id, msg, kommando.Data)
}

// StopPacket creates a new StopPacket
func StopPacket(id string, msg []byte) Packet {
	return new(id, msg, kommando.Stop)
}

// new creates a new Packet
func new(id string, msg []byte, cmd kommando.Kommando) Packet {
	return Packet{
		SessionID: id,
		Message:   msg,
		Kommando:  cmd,
	}
}
