package lambda

import (
	"encoding/json"

	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	client *lambda.Lambda
}

func NewClient() Client {
	log.Trace("Setting up lambda client")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return Client{
		client: lambda.New(sess),
	}
}

func (l Client) General(fName string, in interface{}) (out map[string]interface{}, err error) {
	payload, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	log.Infof("Invoking %s", fName)
	response, err := l.client.Invoke(&lambda.InvokeInput{
		FunctionName: &fName,
		Payload:      payload,
	})
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	log.Tracef("%s response: %s", fName, string(response.Payload))
	err = json.Unmarshal(response.Payload, &result)
	return result, err
}

func (l Client) GeneralAsync(fName string, in interface{}) (err error) {
	payload, err := json.Marshal(in)
	if err != nil {
		return err
	}

	log.Infof("Invoking %s asynchronously", fName)
	_, err = l.client.Invoke(&lambda.InvokeInput{
		FunctionName:   &fName,
		Payload:        payload,
		InvocationType: aws.String(lambda.InvocationTypeEvent),
	})
	return
}

func (l Client) Validate(number string, branch string, accountType string) (map[string]string, error) {
	data := map[string]interface{}{
		"number":       number,
		"branch":       branch,
		"account_type": accountType,
	}
	out, err := l.General("cloudwabbit", data)
	if errlib.ErrorError(err, "Could not get cloudwabbit result") {
		return nil, err
	}

	errs := out["errors"].(map[string]any)
	// casting out["errors"] directly to map[string]string panics, that is why we have this two-step conversion
	result := make(map[string]string, len(errs))
	for k, v := range errs {
		result[k] = v.(string)
	}

	return result, nil
}
