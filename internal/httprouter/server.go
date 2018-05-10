package httprouter

import (
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
)

// GzipOn enables gzip compression for Server.
func GzipOn(s *Server) {
	s.gzip = true
}

// ReadTimeout sets Server response read timeout.
func ReadTimeout(timeout time.Duration) func(*Server) {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// Server is a http server.
type Server struct {
	server http.Server
	gzip   bool
}

// ListenAndServe starts server.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Close closes the server.
func (s *Server) Close() error {
	return s.server.Close()
}

// NewServer returns initialized Server.
func NewServer(addr string, h *Handler, options ...func(*Server)) *Server {
	r := httprouter.New()
	r.GET("/", h.GetIndex)

	s := Server{
		server: http.Server{
			Addr:    addr,
			Handler: r,
		},
	}

	for _, option := range options {
		option(&s)
	}

	if s.gzip {
		s.server.Handler = gziphandler.GzipHandler(r)
	}

	return &s
}
