package cmd

import (
	"log"

	"github.com/c-f/talje/lib/gate"
)

func setProtocol(protocol string) {
	switch protocol {
	case "json":
		gate.DefaultGate = gate.Json()
	case "binary":
		gate.DefaultGate = gate.Binary()
	default:
		log.Fatal("Unknown Protocol", protocol)
	}
}
