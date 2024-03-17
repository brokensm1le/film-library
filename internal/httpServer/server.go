package httpServer

import (
	"film_library/config"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(); err != nil {
		log.Fatalf("Cannot map handlers. Error: {%s}", err)
	}

	log.Printf("Start server on {host:port - %s:%s}", s.cfg.Server.Host, s.cfg.Server.Port)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port), nil); err != nil {
		log.Fatalf("Cannot listen. Error: {%s}", err)
	}
	return nil
}
