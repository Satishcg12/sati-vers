package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/satishcg12/sati-vers/auth/authorization-service/repository"
	"github.com/satishcg12/sati-vers/auth/authorization-service/utils"
)

func (s *Server) LoadRoutes() {

	// a good logger message
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${status} ${method} ${uri} ${latency_human}` + "\n",
	})

	s.router.Use(logger)
	// Add routes here
	s.router.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	authRoutes := s.router.Group("/api/v1/auth")

	repo := repository.New(s.db)
	handler := NewHandler(repo, s.db)
	s.router.Validator = utils.NewValidator()

	authRoutes.POST("/register", handler.Register)
	authRoutes.POST("/login", handler.Login)
	authRoutes.POST("/logout", func(c echo.Context) error {
		return c.String(http.StatusOK, "Logout")
	})

}
