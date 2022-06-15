package lib

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type IDdb interface {
	GetUsers() ([]User, error)
}
type Ddb struct {
	client *dynamodb.Client
}

func (d Ddb) GetUsers() ([]User, error) {
	var key map[string]types.AttributeValue
	users := []User{}
	for loop := true; loop; loop = (key != nil) {
		params := &dynamodb.ScanInput{
			TableName: aws.String("SpotifyProfile"),
		}
		if key != nil {
			params.ExclusiveStartKey = key
		}
		r, err := d.client.Scan(context.TODO(), params)
		if err != nil {
			return nil, err
		}
		for _, v := range r.Items {
			u := User{}
			err = attributevalue.UnmarshalMap(v, &u)
			if err != nil {
				return nil, err
			}
			users = append(users, u)
		}
		key = r.LastEvaluatedKey
	}
	return users, nil
}

func NewDdb() Ddb {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
	if err != nil {
		log.Fatal(err)
	}

	client := dynamodb.NewFromConfig(cfg)
	d := Ddb{client}
	return d
}
