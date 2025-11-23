package httpserver

import (
	"net/http"
	"time"
)

type Server struct {
	http *http.Server
}

func New(port string, h http.Handler) *Server {
	return &Server{
		http: &http.Server{
			Addr:              ":" + port,
			Handler:           h,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}
