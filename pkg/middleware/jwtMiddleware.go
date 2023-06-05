package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
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
	tokenString, err := extractToken(r)
	if err != nil {
		return err
	}

	token, err := parseToken(tokenString)
	if err != nil {
		return err
	}

	err = verifyClaims(token)
	if err != nil {
		return err
	}

	return nil
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("cannot find token in header")
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	return tokenString, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, errors.New("token parsing error")
	}

	return token, nil
}

func verifyClaims(token *jwt.Token) error {
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return errors.New("token not valid")
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return errors.New("token has expired")
	}
	return nil
}
