package server

// Server
type Server interface {
	Start() error
	Stop() error
}
