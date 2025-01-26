package api

import (
	"database/sql"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/Satishcg12/sati-vers/sso-mono/repository"
	"github.com/Satishcg12/sati-vers/sso-mono/types"
	"github.com/Satishcg12/sati-vers/sso-mono/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db   *sql.DB
	repo *repository.Queries
	jwt  *utils.JWT
}
type IHandler interface {
	Register(c echo.Context) error
	Authorize(c echo.Context) error
	Login(c echo.Context) error
}

func NewHandler(repo *repository.Queries, db *sql.DB, jwt *utils.JWT) IHandler {
	return &Handler{
		repo: repo,
		db:   db,
		jwt:  jwt,
	}
}

type (
	RegisterRequest struct {
		Username        string `json:"username" validate:"required,min=3,max=50"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=6,max=50"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
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

func (h *Handler) Authorize(c echo.Context) error {
	type AuthorizeRequest struct {
		ResponseType string `query:"response_type" validate:"required"`
		ClientID     string `query:"client_id" validate:"required"`
		RedirectURI  string `query:"redirect_uri" validate:"required"`
		Scopes       string `query:"scopes" validate:"required"`
		State        string `query:"state" validate:"required"`
	}

	type AuthorizeResponse struct {
		AuthCode string `json:"auth_code"`
		State    string `json:"state"`
	}

	// retrive the request body
	var req AuthorizeRequest
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

	// check if the client exists
	clientId, err := uuid.Parse(req.ClientID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "INVALID_CLIENT_ID",
			Message:   "Invalid client id",
		})
	}

	client, err := h.repo.GetClientById(c.Request().Context(), clientId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "CLIENT_NOT_FOUND",
			Message:   "Client not found",
		})
	}

	// check client redirect uri
	if !slices.Contains(client.ClientRedirectUris, req.RedirectURI) {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success:   false,
			ErrorCode: "INVALID_REDIRECT_URI",
			Message:   "Invalid redirect uri",
		})
	}
	// check if the response type is code
	if req.ResponseType != "code" {
		uri := url.Values{}
		uri.Add("error", "unsupported_response_type")
		uri.Add("state", req.State)
		return c.Redirect(http.StatusFound, req.RedirectURI+"?"+uri.Encode())
	}

	// check if the scopes are valid
	scopes := strings.Split(req.Scopes, " ")
	for _, scope := range scopes {
		if !slices.Contains(client.ClientScopes, scope) {
			return c.Redirect(http.StatusFound, req.RedirectURI+"?error=invalid_scope&state="+req.State)
		}
	}

	// check user session in cookie
	// if not redirect to login page

	session, err := c.Cookie("session_token")
	if err != nil || session.Value == "" {
		// encode the request params and redirect to login page
		token, err := h.jwt.GenerateJWTWtihHS256(map[string]interface{}{
			"client_id":     req.ClientID,
			"redirect_uri":  req.RedirectURI,
			"response_type": req.ResponseType,
			"scopes":        req.Scopes,
			"state":         req.State,
		}, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sso-mono",
			Subject:   "login_request",
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		})

		if err != nil {
			return c.Redirect(http.StatusFound, req.RedirectURI+"?error=server_error&state="+req.State)
		}

		uri := url.Values{}
		uri.Add("auth_request", token)
		return c.Redirect(http.StatusFound, "/login?"+uri.Encode())
	}

	// check if the user has already authorized the client
	// generate auth code
	authCode, err := utils.GenerateString(32)
	if err != nil {
		return c.Redirect(http.StatusFound, req.RedirectURI+"?error=server_error&state="+req.State)
	}

	// hash the auth code
	authCodeHash, err := utils.HashPasswordWithArgon2(authCode, &utils.Algon2Params{
		Memory:      64 * 1024, // 64MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  32,
		KeyLength:   32,
	})

	if err != nil {
		return c.Redirect(http.StatusFound, req.RedirectURI+"?error=server_error&state="+req.State)
	}

	// insert the auth code
	_, err = h.repo.CreateAuthCode(c.Request().Context(), repository.CreateAuthCodeParams{
		ClientID:     uuid.NullUUID{UUID: client.ClientID, Valid: true},
		UserID:       uuid.NullUUID{},
		AuthCodeHash: authCodeHash,
		RedirectUri:  req.RedirectURI,
		Scopes:       scopes,
		ExpiresAt:    time.Now().Add(time.Minute * 5),
	})

	if err != nil {
		return c.Redirect(http.StatusFound, req.RedirectURI+"?error=server_error&state="+req.State)
	}

	uri := url.Values{}
	uri.Add("code", authCode)
	uri.Add("state", req.State)
	return c.Redirect(http.StatusFound, req.RedirectURI+"?"+uri.Encode())
}

func (h *Handler) Login(c echo.Context) error {
	type LoginRequest struct {
		Identifier  string `json:"identifier" validate:"required"`
		Password    string `json:"password" validate:"required"`
		AuthRequest string `query:"auth_request" validate:"required"`
	}

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
