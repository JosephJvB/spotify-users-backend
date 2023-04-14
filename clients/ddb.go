package clients

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"jaf-unwrapped.com/users/models"
)

type IDdb interface {
	GetUsers() ([]models.User, error)
}
type Ddb struct {
	client *dynamodb.Client
}

func (d Ddb) GetUsers() ([]models.User, error) {
	var key map[string]types.AttributeValue
	users := []models.User{}
	for loop := true; loop; loop = (key != nil) {
		params := &dynamodb.ScanInput{
			TableName: aws.String("SpotifyProfile"),
		}
		if key != nil {
			params.ExclusiveStartKey = key
		}
		r, err := d.client.Scan(context.TODO(), params)
		if err != nil {
			return users, err
		}
		for _, v := range r.Items {
			u := models.User{}
			err = attributevalue.UnmarshalMap(v, &u)
			if err != nil {
				return users, err
			}
			users = append(users, u)
		}
		key = r.LastEvaluatedKey
	}
	return users, nil
}

func NewDdb() Ddb {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatal(err)
	}

	client := dynamodb.NewFromConfig(cfg)
	d := Ddb{client}
	return d
}
