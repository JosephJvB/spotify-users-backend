package clients

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	"jaf-unwrapped.com/users/models"
)

type IAuth interface {
	Decode(token string) (models.JWTClaims, error)
	Encode(claims models.JWTClaims) (string, error)
}
type Auth struct {
	JwtSecret []byte
}

func (a Auth) Decode(tokenStr string) (models.JWTClaims, error) {
	claims := models.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return a.JwtSecret, nil
	})
	if err != nil {
		return claims, err
	}
	if token == nil || !token.Valid {
		return claims, errors.New("Invalid token, token not valid at ParseWithClaims")
	}
	return claims, nil
}

func (a Auth) Encode(claims models.JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.JwtSecret)
}

func NewAuth() Auth {
	a := Auth{
		JwtSecret: []byte(os.Getenv("JwtSecret")),
	}
	return a
}
