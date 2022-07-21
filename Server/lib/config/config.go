package config

import "github.com/c-f/talje/lib/ladung"

// Config is the configuration for a new laeufer
type Config struct {

	// InC is a channel of Packets which should be send
	InC chan ladung.Packet
	// OutC is a channel of Packets which are received
	OutC chan ladung.Packet

	// Befestigung defines if the communication is between agents
	// or if additional encoding/decoding is required
	CommType  Communication
	Direction Direction
}

// NewAgent creates a new Agent2Agent configuration
func NewAgent(dir Direction) *Config {
	return new(dir, Agent2Agent)
}

// NewServer creates a Direct Receiving config
func NewServer() *Config {
	return new(Receiving, Direct)
}

// NewServer creates a Direct Sending config
func NewClient() *Config {
	return new(Sending, Direct)
}

// new creates a new Config
func new(direction Direction, commType Communication) *Config {
	return &Config{
		CommType:  commType,
		Direction: direction,

		InC:  make(chan ladung.Packet),
		OutC: make(chan ladung.Packet),
	}
}

// Bind binds the packets channels together
func (c *Config) Bind(other *Config) {
	other.InC = c.OutC
	other.OutC = c.InC
}
