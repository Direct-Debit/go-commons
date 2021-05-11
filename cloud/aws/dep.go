package aws

import (
	"github.com/Direct-Debit/go-commons/cloud/aws/lambda"
)

// Deprecated: rather use the Client in cloud/aws/lambda
type LambdaClient struct {
	Client lambda.Client
}

// Deprecated: will be removed in v1 of this API
//goland:noinspection GoUnusedExportedFunction
func NewLambdaClient() LambdaClient {
	return LambdaClient{lambda.NewClient()}
}

// Deprecated: rather use the Client in cloud/aws/lambda
func (l LambdaClient) Validate(number string, branch string, accountType string) (map[string]string, error) {
	return l.Client.Validate(number, branch, accountType)
}
