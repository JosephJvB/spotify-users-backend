package lib

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
)

type IAuth interface {
	Decode(token string) (interface{}, error)
	Encode(interface{}) (string, error)
}
type Auth struct {
	JwtSecret []byte
}

func (a *Auth) Decode(tokenStr string) (JWTClaims, error) {
	claims := JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
		// 	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		// }

		return a.JwtSecret, nil
	})
	if err != nil {
		return claims, err
	}

	if !token.Valid {
		return claims, errors.New("Decode error: Token invalid")
	}
	return claims, nil
}

func (a *Auth) Encode(claims JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.JwtSecret)
}

func NewAuth() *Auth {
	a := &Auth{
		JwtSecret: []byte(os.Getenv("JafJwtSecret")),
	}
	return a
}
