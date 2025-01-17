package api

import (
	"database/sql"
	"net/http"

	"github.com/Satishcg12/sati-vers/sso-mono/repository"
	"github.com/Satishcg12/sati-vers/sso-mono/types"
	"github.com/Satishcg12/sati-vers/sso-mono/utils"
	"github.com/google/uuid"
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
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=6,max=50"`
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

	// get the credentials
	credentials, err := h.repo.GetSaltAndCredentialsByUserId(c.Request().Context(), uuid.NullUUID{
		Valid: true, UUID: user.ID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success:   false,
			ErrorCode: "FAILED_TO_GET_CREDENTIALS",
			Message:   "Failed to get credentials" + err.Error(),
		})
	}

	if credentials.UserCredential.CredentialType != "password" {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "INVALID_CREDENTIAL_TYPE",
			Message:   "Invalid credential type",
		})
	}

	// verify the password
	if !utils.VerifyPassword(req.Password, credentials.Salt.SaltValue, credentials.UserCredential.CredentialValue.String) {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "INVALID_CREDENTIALS",
			Message:   "Invalid credentials",
		})
	}

	// request for a token from the token service with grpc

	return c.JSON(http.StatusOK, types.Response{
		Success: true,
		Message: "User logged in successfully",
	})
}
