package api

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/Satishcg12/sati-vers/sso-mono/config"
	"github.com/labstack/echo/v4"
)

type Server struct {
	cfg    config.HTTPServer
	router *echo.Echo
	db     *sql.DB
}

func NewServer(cfg config.HTTPServer, db *sql.DB) *Server {
	srv := &Server{
		cfg:    cfg,
		router: echo.New(),
		db:     db,
	}
	srv.LoadRoutes()

	return srv
}

func (s *Server) Start(ctx context.Context) {

	go func() {
		if err := s.router.Start(":" + strconv.Itoa(s.cfg.Port)); err != nil {
			s.router.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.router.Shutdown(ctx); err != nil {
		s.router.Logger.Fatal(err)
	}
}
