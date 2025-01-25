package api

import (
	"net/http"

	"github.com/Satishcg12/sati-vers/sso-mono/repository"
	"github.com/Satishcg12/sati-vers/sso-mono/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) LoadRoutes() {

	if s.cfg.Debug {
		// a good logger message
		logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `${time_rfc3339_nano} ${status} ${method} ${uri} ${latency_human}` + "\n",
		})

		s.router.Use(logger)
	}
	// Add routes here
	s.router.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	r := s.router.Group("/api/v1")

	repo := repository.New(s.db)
	jwt := utils.NewJWT(s.cfg.SecretKey)
	handler := NewHandler(repo, s.db, jwt)
	s.router.Validator = utils.NewValidator()

	authGroup := r.Group("/auth")
	authGroup.POST("/register", handler.Register)

	authGroup.GET("/authorize", handler.Authorize)
	authGroup.POST("/login", handler.Login)
	authGroup.POST("/logout", func(c echo.Context) error {
		return c.String(http.StatusOK, "Logout")
	})

}
