package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func NewToken(userId string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 3).Unix(),
		Issuer: userId,
		IssuedAt: time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func parseJwtCallback(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf(("unexpected signing method: %v"), token.Header["alg"])
	}
	return []byte("secret"), nil
}

func ExtractToken(r *http.Request) (string, error) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return "", fmt.Errorf("no authorization header found")
	}
	if !strings.HasPrefix(header, "Bearer") {
		return "", fmt.Errorf("invalid authorization header")
	}

	return strings.TrimSpace(strings.Split(header, " ")[1]), nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, parseJwtCallback)
}
