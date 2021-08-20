package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

const tableName = "config"
const keyColumnName = "key"

type Config struct{}

func (c Config) Query(key string) (interface{}, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	connection := dynamodb.New(sess)

	dbKey, err := dynamodbattribute.MarshalMap(map[string]interface{}{keyColumnName: key})
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal key for dynamo config table")
	}

	item, err := connection.GetItem(&dynamodb.GetItemInput{
		Key:       dbKey,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "could query dynamo config table")
	}

	if len(item.Item) == 0 {
		return nil, nil
	}

	var value interface{}
	err = dynamodbattribute.Unmarshal(item.Item["value"], &value)
	return value, errors.Wrap(err, "Could not unmarshal dynamo config value")
}

func (c Config) GetDef(key string, def interface{}) (interface{}, error) {
	v, err := c.Query(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return def, nil
	}
	return v, nil
}

func (c Config) Get(key string) (interface{}, error) {
	v, err := c.Query(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, fmt.Errorf("%s not configured in dynamo", key)
	}
	return v, nil
}
