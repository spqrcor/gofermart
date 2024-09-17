package authenticate

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}
type ContextKey string

var ContextUserID ContextKey = "UserID"

type Authenticate struct {
	secretKey string
	tokenExp  time.Duration
}

func NewAuthenticateService(secretKey string, tokenExp time.Duration) *Authenticate {
	return &Authenticate{secretKey: secretKey, tokenExp: tokenExp}
}

func (a *Authenticate) createCookie(UserID uuid.UUID) (http.Cookie, error) {
	token, err := a.createToken(UserID)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{Name: "Authorization", Value: token, Expires: time.Now().Add(a.tokenExp), HttpOnly: true, Path: "/"}, nil
}

func (a *Authenticate) createToken(UserID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenExp)),
		},
		UserID: UserID,
	})

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}
	return "Bearer " + tokenString, nil
}

func (a *Authenticate) GetUserIDFromCookie(tokenString string) (uuid.UUID, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(strings.TrimPrefix(tokenString, "Bearer "), claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(a.secretKey), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	return claims.UserID, nil
}

func (a *Authenticate) SetCookie(rw http.ResponseWriter, UserID uuid.UUID) {
	cookie, err := a.createCookie(UserID)
	if err == nil {
		http.SetCookie(rw, &cookie)
	}
}
