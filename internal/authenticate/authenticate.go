package authenticate

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spqrcor/gofermart/internal/logger"
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

const tokenExp = time.Hour * 3
const secretKey = "KLJ-fo3Fksd3fl!="

func CreateCookie(UserID uuid.UUID) (http.Cookie, error) {
	token, err := createToken(UserID)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{Name: "Authorization", Value: token, Expires: time.Now().Add(tokenExp), HttpOnly: true, Path: "/"}, nil
}

func createToken(UserID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: UserID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return "Bearer " + tokenString, nil
}

func GetUserIDFromCookie(tokenString string) (uuid.UUID, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(strings.TrimPrefix(tokenString, "Bearer "), claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	return claims.UserID, nil
}

func SetCookie(rw http.ResponseWriter, UserID uuid.UUID) {
	cookie, err := CreateCookie(UserID)
	if err != nil {
		logger.Log.Error(err.Error())
	} else {
		http.SetCookie(rw, &cookie)
	}
}
