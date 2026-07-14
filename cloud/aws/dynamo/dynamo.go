package dynamo

import (
	"context"
	"fmt"
	"github.com/Direct-Debit/go-commons/cloud"
	"github.com/Direct-Debit/go-commons/concurrency"
	"github.com/Direct-Debit/go-commons/stdext"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
)

var connection *dynamodb.Client

type Item map[string]dynamotypes.AttributeValue

func Connect() *dynamodb.Client {
	if connection != nil {
		return connection
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	connection = dynamodb.NewFromConfig(cfg)
	return connection
}

func TableExists(tableName *string) bool {
	db := Connect()

	descOutput, err := db.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{TableName: tableName})
	if err != nil {
		return false
	}
	return strings.ToLower(string(descOutput.Table.TableStatus)) == "active"
}

func PutItems(items []Item, table *string, delete bool) {
	var wg sync.WaitGroup
	for i := 0; i < len(items); i += 25 {
		max := i + 25
		if max > len(items) {
			max = len(items)
		}
		submitItems := items[i:max]

		writeRequests := make([]dynamotypes.WriteRequest, len(submitItems))
		for idx, item := range submitItems {
			if delete {
				writeRequests[idx] = dynamotypes.WriteRequest{
					DeleteRequest: &dynamotypes.DeleteRequest{Key: item},
				}
			} else {
				writeRequests[idx] = dynamotypes.WriteRequest{
					PutRequest: &dynamotypes.PutRequest{Item: item},
				}
			}
		}

		log.Debugf("Writing %d items to dynamo table %s", len(submitItems), *table)
		wg.Add(1)
		go func() {
			defer wg.Done()
			putItemsWithBackoff(
				map[string][]dynamotypes.WriteRequest{*table: writeRequests},
				1)
		}()
	}
	wg.Wait()
}

func putItemsWithBackoff(items map[string][]dynamotypes.WriteRequest, backoff int) {
	db := Connect()

	if backoff < 1 {
		backoff = 1
	}

	out, err := db.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: items,
	})
	errlib.PanicError(err, "Couldn't write batch")
	if len(out.UnprocessedItems) > 0 {
		log.Warn(fmt.Sprintf(
			"Failed to write %d items, retrying after %d milliseconds",
			len(out.UnprocessedItems),
			backoff,
		))
		time.AfterFunc(time.Duration(backoff)*time.Millisecond, func() {
			backoff += int(math.Ceil(rand.Float64() * float64(backoff)))
			putItemsWithBackoff(out.UnprocessedItems, backoff)
		})
	}
}

func GetItemsLambda(keys []Item, table *string, batchSize int) []Item {
	res := make([]Item, 0, len(keys))
	output := make(chan []Item)

	callCount := 0
	for i := 0; i < len(keys); i += batchSize {
		max := i + batchSize
		if max > len(keys) {
			max = len(keys)
		}
		go invokeGetItemsLambda(keys[i:max], table, output)
		callCount++
	}

	for i := 1; i <= callCount; i++ {
		res = append(res, <-output...)
		log.Debugf("GetItemsLambda has returned %d items out of %d", len(res), len(keys))
		log.Debugf("GetItemsLambda received %d results from %d calls", i, callCount)
	}

	return res
}

// toV1AttributeValue converts a V2 dynamotypes.AttributeValue into the JSON
// shape produced by the old aws-sdk-go v1 dynamodb.AttributeValue struct
// (e.g. {"S": "foo"}, {"N": "123"}, {"M": {...}}), since the cmn-dynamo-get-items
// lambda still expects that wire format and cannot be changed.
func toV1AttributeValue(av dynamotypes.AttributeValue) map[string]interface{} {
	switch v := av.(type) {
	case *dynamotypes.AttributeValueMemberS:
		return map[string]interface{}{"S": v.Value}
	case *dynamotypes.AttributeValueMemberN:
		return map[string]interface{}{"N": v.Value}
	case *dynamotypes.AttributeValueMemberB:
		return map[string]interface{}{"B": v.Value}
	case *dynamotypes.AttributeValueMemberSS:
		return map[string]interface{}{"SS": v.Value}
	case *dynamotypes.AttributeValueMemberNS:
		return map[string]interface{}{"NS": v.Value}
	case *dynamotypes.AttributeValueMemberBS:
		return map[string]interface{}{"BS": v.Value}
	case *dynamotypes.AttributeValueMemberBOOL:
		return map[string]interface{}{"BOOL": v.Value}
	case *dynamotypes.AttributeValueMemberNULL:
		return map[string]interface{}{"NULL": v.Value}
	case *dynamotypes.AttributeValueMemberL:
		list := make([]interface{}, len(v.Value))
		for i, elem := range v.Value {
			list[i] = toV1AttributeValue(elem)
		}
		return map[string]interface{}{"L": list}
	case *dynamotypes.AttributeValueMemberM:
		m := make(map[string]interface{}, len(v.Value))
		for key, elem := range v.Value {
			m[key] = toV1AttributeValue(elem)
		}
		return map[string]interface{}{"M": m}
	default:
		return nil
	}
}

func toV1Item(item Item) map[string]interface{} {
	v1Item := make(map[string]interface{}, len(item))
	for key, av := range item {
		v1Item[key] = toV1AttributeValue(av)
	}
	return v1Item
}

func invokeGetItemsLambda(items []Item, table *string, output chan []Item) {
	v1Items := make([]map[string]interface{}, len(items))
	for i, item := range items {
		v1Items[i] = toV1Item(item)
	}

	input := map[string]interface{}{
		"items": v1Items,
		"table": *table,
	}
	result, err := cloud.CallFunc("cmn-dynamo-get-items", input)
	errlib.ErrorError(err, "failed to invoke dynamo-get-items lambda")

	var resp []Item
	for _, item := range result["items"].([]interface{}) {
		dynamoItem, err := attributevalue.MarshalMap(item)
		if errlib.ErrorError(err, "failed to marshal attribute values") {
			continue
		}
		resp = append(resp, dynamoItem)
	}
	output <- resp
}

func GetItems(items []Item, table *string) []Item {
	res := make([]Item, 0, len(items))
	for i := 0; i < len(items); i += 100 {
		max := i + 100
		if max > len(items) {
			max = len(items)
		}

		getItems := make([]map[string]dynamotypes.AttributeValue, 0, max-i)
		for item := i; item < max; item++ {
			getItems = append(getItems, items[item])
		}

		keysAndAttr := dynamotypes.KeysAndAttributes{
			Keys: getItems,
		}

		log.Debugf("Getting %d items from dynamo table %s", len(getItems), *table)
		subRes := getItemsWithBackoff(
			map[string]dynamotypes.KeysAndAttributes{*table: keysAndAttr}, 1,
		)

		res = append(res, subRes[*table]...)
	}
	return res
}

func GetItemsConcurrent(workers int, items []Item, table *string) []Item {
	itemChunks := stdext.ChunkifyBySize(items, 100)
	returnedChunks := concurrency.Workers(workers, itemChunks,
		func(chunk []Item) (ret []Item, success bool) {
			defer func() {
				if recover() != nil {
					ret, success = nil, false
				}
			}()

			keys := stdext.Map(chunk,
				func(item Item) map[string]dynamotypes.AttributeValue {
					return item
				},
			)
			keysAndAttr := dynamotypes.KeysAndAttributes{
				Keys: keys,
			}

			log.Debugf("Getting %d items from dynamo table %s", len(chunk), *table)
			itemsReturned := getItemsWithBackoff(
				map[string]dynamotypes.KeysAndAttributes{*table: keysAndAttr}, 1,
			)
			ret, success = itemsReturned[*table], true
			return
		},
	)
	return stdext.Flatten(returnedChunks)
}

func getItemsWithBackoff(items map[string]dynamotypes.KeysAndAttributes, backoff int) map[string][]Item {
	db := Connect()

	if backoff < 1 {
		backoff = 1
	}

	out, err := db.BatchGetItem(context.TODO(), &dynamodb.BatchGetItemInput{
		RequestItems: items,
	})
	errlib.PanicError(err, "Couldn't get batch")

	var retriesResult map[string][]Item
	if len(out.UnprocessedKeys) > 0 {
		log.Warn(fmt.Sprintf(
			"Failed to get %d items, retrying after %d milliseconds",
			len(out.UnprocessedKeys),
			backoff,
		))
		time.Sleep(time.Duration(backoff) * time.Millisecond)
		backoff += int(math.Ceil(rand.Float64() * float64(backoff)))
		retriesResult = getItemsWithBackoff(out.UnprocessedKeys, backoff)
	}
	res := make(map[string][]Item)
	for key, response := range out.Responses {
		items := make([]Item, 0, len(response))
		for _, item := range response {
			items = append(items, item)
		}
		res[key] = append(items, retriesResult[key]...)
	}

	return res
}

func QueryAll(initialInput *dynamodb.QueryInput) ([]Item, error) {
	res := make([]Item, 0)

	queryDone := false
	db := Connect()
	for !queryDone {
		result, err := db.Query(context.TODO(), initialInput)
		if err != nil {
			return res, err
		}
		for _, item := range result.Items {
			res = append(res, item)
		}

		queryDone = len(result.LastEvaluatedKey) == 0
		initialInput.ExclusiveStartKey = result.LastEvaluatedKey
	}

	return res, nil
}

type CountOutput struct {
	Count        int64
	ScannedCount int64
}

func CountAll(initialInput *dynamodb.QueryInput) (CountOutput, error) {
	initialInput.Select = dynamotypes.SelectCount
	res := CountOutput{}

	queryDone := false
	db := Connect()
	for !queryDone {
		result, err := db.Query(context.TODO(), initialInput)
		if err != nil {
			return res, err
		}
		res.Count += int64(result.Count)
		res.ScannedCount += int64(result.ScannedCount)

		queryDone = len(result.LastEvaluatedKey) == 0
		initialInput.ExclusiveStartKey = result.LastEvaluatedKey
	}

	return res, nil
}

func ScanAll(initialInput *dynamodb.ScanInput) ([]Item, error) {
	res := make([]Item, 0)

	scanDone := false
	db := Connect()
	for !scanDone {
		result, err := db.Scan(context.TODO(), initialInput)
		if err != nil {
			return res, err
		}
		for _, item := range result.Items {
			res = append(res, item)
		}

		scanDone = len(result.LastEvaluatedKey) == 0
		initialInput.ExclusiveStartKey = result.LastEvaluatedKey
	}

	return res, nil
}
