package aws

import (
	"github.com/Direct-Debit/go-commons/cloud/aws/lambda"
)

// Deprecated: will be removed in v1 of this API
//goland:noinspection GoUnusedExportedFunction
func NewLambdaClient() lambda.Client {
	return lambda.NewClient()
}
