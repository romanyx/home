package httprouter

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
)

// GzipOn enables gzip compression for Server.
func GzipOn(s *Server) {
	s.gzip = true
}

// Letsencrypt sets TLS config
func Letsencrypt(s *Server) {
	m := &autocert.Manager{
		Cache:      autocert.DirCache("/acme"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("romanyx.info"),
	}

	go http.ListenAndServe(":http", m.HTTPHandler(nil))

	s.server.TLSConfig = &tls.Config{
		GetCertificate: m.GetCertificate,
	}
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

// ListenAndServeLetsencrypt starts tls server
func (s *Server) ListenAndServeLetsencrypt() error {
	return s.server.ListenAndServeTLS("", "")
}

// Close closes the server.
func (s *Server) Close() error {
	return s.server.Close()
}

func redirectToInfo(fn httprouter.Handle) httprouter.Handle {
	f := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		url := r.URL
		if url.Hostname() == "romanyx.ru" {
			url.Host = "romanyx.info"
			http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
			return
		}

		fn(w, r, p)
	}

	return f
}

func (h *Handler) recover(fn httprouter.Handle) httprouter.Handle {
	f := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if r := recover(); r != nil {
				h.logFunc(fmt.Errorf("recover: %v", r))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "internal server error")
			}
		}()

		fn(w, r, p)
	}

	return f
}

// NewServer returns initialized Server.
func NewServer(addr string, h *Handler, options ...func(*Server)) *Server {
	r := httprouter.New()

	r.GET("/", h.recover(redirectToInfo(h.GetIndex)))
	r.GET("/cv", h.recover(redirectToInfo(h.GetCV)))

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
