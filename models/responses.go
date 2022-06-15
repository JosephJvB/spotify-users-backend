package models

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func getCors() map[string]string {
	h := map[string]string{}
	h["Content-Type"] = "application/json"
	h["Allow"] = "*"
	h["Access-Control-Allow-Headers"] = "*"
	h["Access-Control-Allow-Methods"] = "*"
	h["Access-Control-Allow-Origin"] = "*"
	return h
}

type UsersResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Users   []User `json:"users"`
}
type BaseResponse struct {
	Message string `json:"message"`
}

func NewBasicResponse(code int, message string) events.APIGatewayProxyResponse {
	r := &BaseResponse{
		Message: message,
	}
	return newProxyResponse(code, r)
}
func NewUserResponse(users []User, token string) events.APIGatewayProxyResponse {
	r := &UsersResponse{
		Message: "success",
		Users:   users,
		Token:   token,
	}
	return newProxyResponse(200, r)
}
func newProxyResponse(code int, body interface{}) events.APIGatewayProxyResponse {
	p := events.APIGatewayProxyResponse{
		StatusCode:        code,
		Headers:           getCors(),
		MultiValueHeaders: map[string][]string{},
		IsBase64Encoded:   false,
	}
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal("failed to Marshal body to json")
		panic(err)
	}
	p.Body = string(b)

	return p
}
