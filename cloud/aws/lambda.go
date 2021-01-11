package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type LambdaClient struct {
	client *lambda.Lambda
}

func NewLambdaClient() LambdaClient {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return LambdaClient{
		client: lambda.New(sess),
	}
}

func (l LambdaClient) Validate(number string, branch string, accountType string) (map[string]string, error) {
	data := map[string]string{
		"number":       number,
		"branch":       branch,
		"account_type": accountType,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	response, err := l.client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("cloudwabbit"), Payload: payload})
	if err != nil {
		return nil, err
	}
	var result map[string]map[string]string
	err = json.Unmarshal(response.Payload, &result)
	return result["errors"], err
}
