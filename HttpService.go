package main

import (
	"log"
	"net/http"
	"time"
)

// HttpService struct
type HttpService struct {
	addr string
	mux  *http.ServeMux
}

// NewHttpService initializes the service with an address
func NewHttpService(addrs ...string) *HttpService {
	addr := ":9000" // default value
	if len(addrs) > 0 && addrs[0] != "" {
		addr = addrs[0]
	}

	return &HttpService{
		addr: addr,
		mux:  http.NewServeMux(),
	}
}

// RegisterRoute lets you add custom routes/handlers
func (s *HttpService) RegisterRoute(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// Start launches the HTTP server
func (s *HttpService) Start() error {
	go func() {
		if err := http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	// Give it a moment to fail if port is in use
	time.Sleep(250 * time.Millisecond)
	return nil
}
