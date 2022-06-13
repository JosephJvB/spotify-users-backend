package lib

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

type User struct {
	SpotifyId      string `json:"spotifyId"`
	DisplayName    string `json:"displayName"`
	DisplayPicture string `json:"displayPicture"`
}

type UsersResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Users   []User `json:"users"`
}
type BaseResponse struct {
	Message string `json:"message"`
}

func getCors() map[string]string {
	h := map[string]string{}
	h["Content-Type"] = "application/json"
	h["Allow"] = "*"
	h["Access-Control-Allow-Headers"] = "*"
	h["Access-Control-Allow-Methods"] = "*"
	h["Access-Control-Allow-Origin"] = "*"
	return h
}
func NewBasicResponse(code int, message string) events.APIGatewayProxyResponse {
	r := &BaseResponse{
		Message: message,
	}
	b, err := json.Marshal(r)
	if err != nil {
		log.Fatal("failed to Marshal json")
		panic(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        code,
		Headers:           getCors(),
		Body:              string(b),
		MultiValueHeaders: map[string][]string{},
		IsBase64Encoded:   false,
	}
}
func NewUserResponse(users []User, token string) events.APIGatewayProxyResponse {
	r := &UsersResponse{
		Message: "success",
		Users:   users,
		Token:   token,
	}
	b, err := json.Marshal(r)
	if err != nil {
		log.Fatal("failed to Marshal json")
		panic(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           getCors(),
		Body:              string(b),
		MultiValueHeaders: map[string][]string{},
		IsBase64Encoded:   false,
	}
}

type JWTData struct {
	Expires   int64  `json:"expires"`
	SpotifyId string `json:"spotifyId"`
}
type JWTClaims struct {
	Data JWTData `json:"data"`
	jwt.StandardClaims
}
