package protocol

import (
	"encoding/json"
	"io"

	"github.com/c-f/talje/lib/ladung"
)

// Json protocol en/decoder struct
type Json struct {
	encoder *json.Encoder
	decoder *json.Decoder
}

// WritePacket encodes a gut and writes it to writer
func (jl *Json) WritePacket(g ladung.Packet, w io.Writer) (err error) {
	if jl.encoder == nil {
		jl.encoder = json.NewEncoder(w)
	}

	err = jl.encoder.Encode(g)

	return
}

// ReadPacket decodes a Packet from a reader
func (jl *Json) ReadPacket(r io.Reader) (g ladung.Packet, err error) {
	if jl.decoder == nil {
		jl.decoder = json.NewDecoder(r)
	}

	err = jl.decoder.Decode(&g)
	// retry if io.EOF
	if err == io.EOF {

		jl.decoder = json.NewDecoder(r)
		err = jl.decoder.Decode(&g)
	}

	return
}
