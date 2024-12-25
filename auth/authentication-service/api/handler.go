package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/satishcg12/sati-vers/authentication-service/repository"
	"github.com/satishcg12/sati-vers/authentication-service/types"
	"github.com/satishcg12/sati-vers/authentication-service/utils"
)

type Handler struct {
	db   *sql.DB
	repo *repository.Queries
}
type IHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

func NewHandler(repo *repository.Queries, db *sql.DB) *Handler {
	return &Handler{
		repo: repo,
		db:   db,
	}
}

type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6,max=50"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

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
	fmt.Println(err)
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

	salt, err := utils.GenerateSalt(32)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_GENERATE_SALT",
			Message:   "Failed to generate salt",
		})
	}

	// hash the password
	hashedPassword, err := utils.HashPassword(req.Password, salt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_HASH_PASSWORD",
			Message:   "Failed to hash password",
		})
	}

	// create a new user
	// start a transaction
	tx, err := h.db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_START_TRANSACTION",
			Message:   "Failed to start transaction",
		})
	}
	defer tx.Rollback()

	txStart := h.repo.WithTx(tx)

	// create a new user
	userData, err := txStart.CreateUser(c.Request().Context(), repository.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_CREATE_USER",
			Message:   "Failed to create user",
		})
	}
	_, err = txStart.CreateCredentials(c.Request().Context(), repository.CreateCredentialsParams{
		UserID:          uuid.NullUUID{UUID: userData.ID, Valid: true},
		CredentialType:  "password",
		CredentialValue: sql.NullString{String: hashedPassword, Valid: true},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_CREATE_CREDENTIALS",
			Message:   "Failed to create credentials" + err.Error(),
		})
	}
	_, err = txStart.CreateSalt(c.Request().Context(), repository.CreateSaltParams{
		UserID:    userData.ID,
		SaltValue: salt,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_CREATE_SALT",
			Message:   "Failed to create salt " + err.Error(),
		})
	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_COMMIT_TRANSACTION",
			Message:   "Failed to commit transaction",
		})
	}

	// publish notificaiton message

	return c.JSON(http.StatusOK, types.Response{
		Success: true,
		Message: "User Registerd Succesfully ",
	})
}

func (h *Handler) Login(c echo.Context) error {
	return c.String(200, "Login")
}
