package api

import (
	"database/sql"
	"net/http"

	"github.com/Satishcg12/sati-vers/sso-mono/repository"
	"github.com/Satishcg12/sati-vers/sso-mono/types"
	"github.com/Satishcg12/sati-vers/sso-mono/utils"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db   *sql.DB
	repo *repository.Queries
}
type IHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

func NewHandler(repo *repository.Queries, db *sql.DB) IHandler {
	return &Handler{
		repo: repo,
		db:   db,
	}
}

type (
	RegisterRequest struct {
		Username        string `json:"username" validate:"required,min=3,max=50"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=6,max=50"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	}
	LoginRequest struct {
		Identifier   string `json:"identifier" validate:"required"`
		Password     string `json:"password" validate:"required,min=6,max=50"`
		ClientID     string `query:"client_id" validate:"required"`
		RedirectURI  string `query:"response_uri" validate:"required"`
		ResponseType string `query:"response_type" validate:"required"`
		Scopes       string `query:"scopes" validate:"required"`
		State        string `query:"state" validate:"required"`
	}
)

func (h *Handler) Register(c echo.Context) error {
	// retrive the request body
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "INVALID_REQUEST_BODY",
			Message:   "Invalid request body",
		})
	}

	// validate the request body
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponseWithData{
			Success:   false,
			ErrorCode: "INVALID_REQUEST_BODY",
			Message:   "Invalid request body",
			Data:      err.(*echo.HTTPError).Message,
		})
	}

	// does the user already exist?
	_, err := h.repo.GetUserByEmail(c.Request().Context(), req.Email)
	if err == nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "USER_ALREADY_EXISTS",
			Message:   "User already exists",
		})
	}

	_, err = h.repo.GetUserByUsername(c.Request().Context(), req.Username)
	if err == nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "USER_ALREADY_EXISTS",
			Message:   "User already exists",
		})
	}

	// hash the password
	hashedPassword, err := utils.HashPasswordWithArgon2(req.Password, &utils.Algon2Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  32,
		KeyLength:   32,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_HASH_PASSWORD",
			Message:   "Failed to hash password",
		})
	}

	// insert the user
	_, err = h.repo.CreateUser(c.Request().Context(), repository.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_CREATE_USER",
			Message:   "Failed to create user",
		})
	}

	// publish notificaiton message

	return c.JSON(http.StatusOK, types.Response{
		Success: true,
		Message: "User Registerd Succesfully ",
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "INVALID_REQUEST_BODY",
			Message:   "Invalid request body",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponseWithData{
			Success:   false,
			ErrorCode: "INVALID_REQUEST_BODY",
			Message:   "Invalid request body",
			Data:      err.(*echo.HTTPError).Message,
		})
	}

	// check if the user exists by username or email
	user, err := h.repo.GetUserByUsernameOrEmail(c.Request().Context(), req.Identifier)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "USER_NOT_FOUND",
			Message:   "User not found",
		})
	}

	// verify the password
	isValid, err := utils.VerifyPasswordWithArgon2(req.Password, user.PasswordHash)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_VERIFY_PASSWORD",
			Message:   "Failed to verify password",
		})
	}

	if !isValid {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "INVALID_CREDENTIALS",
			Message:   "Invalid credentials",
		})
	}

	// generate auth code

	return c.JSON(http.StatusOK, types.Response{
		Success: true,
		Message: "User logged in successfully",
	})
}
