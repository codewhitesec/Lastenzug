package config

// Direction defines if the laeufer should wait for incomming connections
// or if it should activly connect to an endpoint to establish it.
type Direction int

const (
	// Receiving indicate that the laeufer will wait for connections
	Receiving Direction = iota + 1

	// Sending indicate that the laeufer will activly create the connection
	Sending
)

// mapping
var directions = map[Direction]string{
	Receiving: "pull",
	Sending:   "push",
}

// String stringify the direction
func (r *Direction) String() string {
	return directions[*r]
}
