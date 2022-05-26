package server

import (
	"github.com/alpaka-alpachino/job-engine/config"
	"github.com/alpaka-alpachino/job-engine/internal/data"
	"html/template"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(c *config.EngineConfig, t *template.Template, categories map[string]data.Category) (*Server, error) {
	router, err := newRouter(t, categories)
	if err != nil {
		return nil, err
	}

	s := Server{
		server: &http.Server{
			Addr:    c.Port,
			Handler: router,
		},
	}

	return &s, nil
}

func (s *Server) RunServer() error {
	return s.server.ListenAndServe()
}
