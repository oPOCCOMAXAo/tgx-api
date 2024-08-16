package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Server struct {
	config Config
	engine *gin.Engine
	server *http.Server
}

type Config struct {
	Port string `env:"PORT" envDefault:"8080"`
}

func New(config Config) *Server {
	res := &Server{
		config: config,
	}

	res.initEngine()
	res.initServer()

	return res
}

func (s *Server) initEngine() {
	s.engine = gin.New()
	_ = s.engine.SetTrustedProxies(nil)
	s.engine.Use(gin.Recovery())
}

func (s *Server) initServer() {
	s.server = &http.Server{
		Addr:              ":" + s.config.Port,
		Handler:           s.engine,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      time.Minute,
	}
}

func (s *Server) Serve(context.Context) error {
	err := s.server.ListenAndServe()

	return errors.WithStack(err)
}

func (s *Server) Shutdown() error {
	err := s.server.Shutdown(context.Background())

	return errors.WithStack(err)
}
