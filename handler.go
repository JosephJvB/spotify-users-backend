package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"jaf-unwrapped.com/users/lib"
)

var (
	auth           *lib.Auth
	adminSpotifyId string
	ddb            *lib.Ddb
)

func init() {
	log.SetPrefix("LoadUsers:")
	log.SetFlags(0)
	auth = lib.NewAuth()
	adminSpotifyId = os.Getenv("AdminSpotifyId")
	ddb = lib.NewDdb()
}

func HandleLambdaEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println(request)
	log.Println("request.HTTPMethod:", request.HTTPMethod)
	log.Println("request.Body:", request.Body)
	if request.HTTPMethod == "OPTIONS" {
		return lib.NewBasicResponse(200, ""), nil
	}

	authHeader, ok := request.Headers["Authorization"]
	if !ok || authHeader == "" {
		msg := "Invalid request, missing Authorization header"
		return lib.NewBasicResponse(400, msg), nil
	}
	var token string
	s := strings.Split(authHeader, " ")
	if len(s) == 2 {
		token = s[1]
	}
	if token == "" {
		msg := "Invalid request, invalid Authorization header"
		return lib.NewBasicResponse(400, msg), nil
	}

	claims, err := auth.Decode(token)
	if err != nil || claims == nil {
		msg := "Invalid request, failed to decode bearer token"
		return lib.NewBasicResponse(400, msg), err
	}

	if claims.Data.SpotifyId != adminSpotifyId {
		msg := "Invalid request, Unauthorized user, not joe!"
		return lib.NewBasicResponse(400, msg), err
	}
	// https://stackoverflow.com/questions/36051177/date-now-equivalent-in-go
	now := time.Now().UTC().UnixNano() / 1e6
	if claims.Data.Expires < now {
		msg := "Invalid request, token expired"
		return lib.NewBasicResponse(400, msg), err
	}

	users, err := ddb.GetUsers()
	if err != nil {
		msg := "Failed to get users from ddb " + err.Error()
		return lib.NewBasicResponse(400, msg), err
	}

	nextClaims := lib.JWTClaims{
		Data: lib.JWTData{
			Expires:   now * 1000,
			SpotifyId: claims.Data.SpotifyId,
		},
	}
	token, err = auth.Encode(nextClaims)
	if err != nil {
		msg := "Failed to encode token " + err.Error()
		return lib.NewBasicResponse(500, msg), err
	}

	return lib.NewUserResponse(
		users,
		token,
	), nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
