package utils

import (
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
func (j *JWT) GenerateJWTWtihHS256(data map[string]interface{}, regClaim jwt.RegisteredClaims) (string, error) {

	claims := jwt.MapClaims{
		"exp": regClaim.ExpiresAt,
		"iat": regClaim.IssuedAt,
		"iss": regClaim.Issuer,
		"aud": regClaim.Audience,
		"sub": regClaim.Subject,
		"nbf": regClaim.NotBefore,
		"jti": regClaim.ID,
	}

	for key, value := range data {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWT) VerifyJWTWtihHS256(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
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
