package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/satishcg12/sati-vers/authentication-service/repository"
)

type Handler struct {
	repo *repository.Queries
}
type IHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

func NewHandler(repo *repository.Queries) *Handler {
	return &Handler{
		repo: repo,
	}
}

type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6,max=50"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func (h *Handler) Register(c echo.Context) error {
	return c.JSON(http.StatusOK, "Register")
}

func (h *Handler) Login(c echo.Context) error {
	return c.String(200, "Login")
}
