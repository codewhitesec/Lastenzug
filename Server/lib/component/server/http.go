package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Http struct {
	Server http.Server // handles connections - implements ListenAndServe
	router *mux.Router

	// Tls options
	Cert string
	Key  string
}

// NewSocksServer returns a new SocksServer
func NewHttp(addr string) *Http {
	// HTTP Handler
	router := mux.NewRouter()

	return &Http{
		Server: http.Server{
			ReadTimeout: 10 * time.Second,
			Addr:        addr,
			Handler:     router,
		},
		router: router,
	}
}

// Bind attach a handlerFunc to the router
func (h *Http) Bind(pattern string, handler http.HandlerFunc) {
	h.router.HandleFunc(pattern, handler)
	h.Server.Handler = h.router
}

// Start starts the HTTP server
func (h *Http) Start() error {
	log.Println("Start HTTP server", h.Server.Addr)

	if h.Cert != "" && h.Key != "" {
		return h.ListenAndServeTLS(h.Cert, h.Key)
	} else {
		return h.ListenAndServe()
	}
}

// ListenAndServe starts the HTTP server
func (h *Http) ListenAndServe() error {
	return h.Server.ListenAndServe()
}

// ListenAndServeTls starts the HTTPS server
func (h *Http) ListenAndServeTLS(certFile, keyFile string) error {
	return h.Server.ListenAndServeTLS(certFile, keyFile)
}

// Close closes the listener
func (h *Http) Close() error {
	return h.Server.Shutdown(context.Background())
}
