package dynamo

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

var connection *dynamodb.DynamoDB

type Item map[string]*dynamodb.AttributeValue

func Connect() *dynamodb.DynamoDB {
	if connection != (*dynamodb.DynamoDB)(nil) {
		return connection
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	connection = dynamodb.New(sess)
	return connection
}

func TableExists(tableName *string) bool {
	db := Connect()

	descOutput, err := db.DescribeTable(&dynamodb.DescribeTableInput{TableName: tableName})
	if err != nil {
		return false
	}
	return strings.ToLower(*descOutput.Table.TableStatus) == "active"
}

func PutItems(items []Item, table *string, delete bool) {
	var wg sync.WaitGroup
	for i := 0; i < len(items); i += 25 {
		max := i + 25
		if max > len(items) {
			max = len(items)
		}
		submitItems := items[i:max]

		writeRequests := make([]*dynamodb.WriteRequest, len(submitItems))
		for idx, item := range submitItems {
			if delete {
				writeRequests[idx] = &dynamodb.WriteRequest{
					DeleteRequest: &dynamodb.DeleteRequest{Key: item},
				}
			} else {
				writeRequests[idx] = &dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{Item: item},
				}
			}
		}

		log.Debugf("Writing %d items to dynamo table %s", len(submitItems), *table)
		wg.Add(1)
		go func() {
			defer wg.Done()
			putItemsWithBackoff(
				map[string][]*dynamodb.WriteRequest{*table: writeRequests},
				1)
		}()
	}
	wg.Wait()
}

func putItemsWithBackoff(items map[string][]*dynamodb.WriteRequest, backoff int) {
	db := Connect()

	if backoff < 1 {
		backoff = 1
	}

	out, err := db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
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

func GetItems(items []Item, table *string) []Item {
	res := make([]Item, 0, len(items))
	for i := 0; i < len(items); i += 100 {
		max := i + 100
		if max > len(items) {
			max = len(items)
		}

		getItems := make([]map[string]*dynamodb.AttributeValue, 0, max-i)
		for item := i; item < max; item++ {
			getItems = append(getItems, items[item])
		}

		keysAndAttr := &dynamodb.KeysAndAttributes{
			Keys: getItems,
		}

		log.Debugf("Getting %d items from dynamo table %s", len(getItems), *table)
		subRes := getItemsWithBackoff(
			map[string]*dynamodb.KeysAndAttributes{*table: keysAndAttr}, 1,
		)

		res = append(res, subRes[*table]...)
	}
	return res
}

func getItemsWithBackoff(items map[string]*dynamodb.KeysAndAttributes, backoff int) map[string][]Item {
	db := Connect()

	if backoff < 1 {
		backoff = 1
	}

	out, err := db.BatchGetItem(&dynamodb.BatchGetItemInput{
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
		result, err := db.Query(initialInput)
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

func ScanAll(initialInput *dynamodb.ScanInput) ([]Item, error) {
	res := make([]Item, 0)

	scanDone := false
	db := Connect()
	for !scanDone {
		result, err := db.Scan(initialInput)
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
