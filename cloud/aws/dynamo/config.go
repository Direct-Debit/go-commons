package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

const tableName = "config"
const keyColumnName = "key"

type Config struct{}

func (c Config) Query(key string) (interface{}, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "could not load aws config")
	}
	connection := dynamodb.NewFromConfig(cfg)

	dbKey, err := attributevalue.MarshalMap(map[string]interface{}{keyColumnName: key})
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal key for dynamo config table")
	}

	table := tableName
	item, err := connection.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       dbKey,
		TableName: &table,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could query dynamo config table")
	}

	if len(item.Item) == 0 {
		return nil, nil
	}

	var value interface{}
	err = attributevalue.Unmarshal(item.Item["value"], &value)
	return value, errors.Wrap(err, "Could not unmarshal dynamo config value")
}

func (c Config) Reload() error {
	return nil
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
