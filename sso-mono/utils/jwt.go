package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SessionToken struct {
	UserID string
}

type AuthRequestToken struct {
	ClientID     string
	RedirectURI  string
	ResponseType string

	Scopes string
	State  string
}
type AccessToken struct {
	UserID string
}
type JWT struct {
	SecretKey string
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

func (j *JWT) GenerateAuthRequestToken(claim AuthRequestToken) (string, error) {

	claims := jwt.MapClaims{
		"client_id":     claim.ClientID,
		"redirect_uri":  claim.RedirectURI,
		"response_type": claim.ResponseType,
		"scopes":        claim.Scopes,
		"state":         claim.State,
		"exp":           time.Now().Add(time.Hour * 1).Unix(),
		"iat":           time.Now().Unix(),
		"iss":           "sso-mono",
		"aud":           "sso-mono",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}
func (j *JWT) GenerateAccessToken(claim AccessToken) (string, error) {

	claims := jwt.MapClaims{
		"user_id": claim.UserID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.SecretKey))
}

// func (j *JWT) VerifyAccessToken(tokenString string) (*jwt.Token, error) {

func (j *JWT) VerifyToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
