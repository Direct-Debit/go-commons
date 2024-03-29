package cloud

import (
	"github.com/Direct-Debit/go-commons/cloud/aws/lambda"
	"sync"
)

var setup sync.Once

type CdvValidator interface {
	Validate(number string, branch string, accountType string) (map[string]string, error)
}

type FunctionCaller interface {
	General(functionName string, in interface{}) (out map[string]interface{}, err error)
	GeneralAsync(functionName string, in interface{}) (err error)
}

// Deprecated: Rather use CdvValidator
type FunctionProvider interface {
	CdvValidator
}

var caller interface {
	CdvValidator
	FunctionCaller
}

func doSetup() {
	setup.Do(func() {
		caller = lambda.NewClient()
	})
}

func CallFunc(functionName string, in interface{}) (out map[string]interface{}, err error) {
	doSetup()
	return caller.General(functionName, in)
}

func CallAsync(functionName string, in interface{}) (err error) {
	doSetup()
	return caller.GeneralAsync(functionName, in)
}

func ValidateCDV(number string, branch string, accountType string) (map[string]string, error) {
	doSetup()
	return caller.Validate(number, branch, accountType)
}
