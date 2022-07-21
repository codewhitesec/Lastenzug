package protocol

import (
	"encoding/binary"

	"github.com/google/uuid"
)

// Int64 defines an 8 byte numerical identifier
type Int64ID int

func NewInt64ID(i int) ID {
	return Int64ID(i)
}

func (id Int64ID) ID() string {
	bts := make([]byte, 8)
	binary.LittleEndian.PutUint64(bts, uint64(id))
	return string(bts)
}

// UUID defines a unique generated ID
type UUID struct {
	uuid.UUID
}

func NewUUID() ID {
	return UUID{uuid.New()}
}

func (id UUID) ID() string {
	return id.UUID.String()
}

// Static ID
type StaticID string

func (id StaticID) ID() string {
	return string(id)
}

func NewStaticID(str string) ID {
	return StaticID(str)
}
