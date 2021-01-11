package cloud

import (
	"github.com/Direct-Debit/go-commons/cloud/aws"
	"sync"
)

var setup sync.Once

type FunctionProvider interface {
	Validate(number string, branch string, accountType string) (map[string]string, error)
}

var provider FunctionProvider

func ValidateCDV(number string, branch string, accountType string) (map[string]string, error) {
	setup.Do(func() {
		provider = aws.NewLambdaClient()
	})
	return provider.Validate(number, branch, accountType)
}
