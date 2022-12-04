package server

import (
	"github.com/alpaka-alpachino/job-engine/config"
	"github.com/alpaka-alpachino/job-engine/internal/service"
	"html/template"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(c *config.EngineConfig, s *service.Service, t *template.Template) (*Server, error) {
	router, err := newRouter(s, t)
	if err != nil {
		return nil, err
	}

	srv := Server{
		server: &http.Server{
			Addr:    c.Port,
			Handler: router,
		},
	}

	return &srv, nil
}

func (s *Server) RunServer() error {
	return s.server.ListenAndServe()
}
