package config

// Communication defines how to read/write to the endpoint
type Communication int

const (

	// Agent2Agent - Güter needs to be en/decoded first
	Agent2Agent Communication = iota + 1

	// Direct - Güter will be send/received directly
	Direct

	// Forward - Stream will be forwarded regardless of content
	// TODO(author): not implemented yet
	Forward
)

// mapping
var communications = map[Communication]string{
	Agent2Agent: "a2a",
	Direct:      "direct",
	Forward:     "forward",
}

// Return the name of a pro
func (b *Communication) String() string {
	return communications[*b]
}
