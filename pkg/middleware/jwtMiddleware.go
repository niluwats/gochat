package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * 60 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(w http.ResponseWriter, r *http.Request) error {
	if r.Header["Token"] == nil {
		return errors.New("cannot find token in header")
	}

	token, err := jwt.ParseWithClaims(r.Header["Token"][0], &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return errors.New("token parsing error")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return errors.New("token not valid")
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return errors.New("token has expired")
	}
	return nil
}
