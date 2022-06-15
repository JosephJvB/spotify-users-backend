package models

import (
	"github.com/golang-jwt/jwt"
)

type User struct {
	SpotifyId      string `json:"spotifyId"`
	DisplayName    string `json:"displayName"`
	DisplayPicture string `json:"displayPicture"`
}

type JWTData struct {
	Expires   int64  `json:"expires"`
	SpotifyId string `json:"spotifyId"`
}
type JWTClaims struct {
	Data JWTData `json:"data"`
	jwt.StandardClaims
}
