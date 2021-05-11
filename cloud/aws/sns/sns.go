package sns

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	snsClient   *sns.SNS
	Topic       string
	Environment string // dev, stage, prod etc...
}

func NewClient(topic string, env string) Client {
	log.Trace("Setting up lambda client")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return Client{
		snsClient:   sns.New(sess),
		Topic:       topic,
		Environment: env,
	}
}

func (s Client) attributes() map[string]*sns.MessageAttributeValue {
	var env sns.MessageAttributeValue
	env.SetStringValue(s.Environment)
	env.SetDataType("String")

	return map[string]*sns.MessageAttributeValue{"env": &env}
}

func (s Client) Publish(subject string, message string) error {

	_, err := s.snsClient.Publish(&sns.PublishInput{
		Message:           &message,
		MessageAttributes: s.attributes(),
		Subject:           &subject,
		TopicArn:          &s.Topic,
	})
	if err != nil {
		log.Errorf("Could not publish to %s: %s", s.Topic, err)
	}
	return err
}
