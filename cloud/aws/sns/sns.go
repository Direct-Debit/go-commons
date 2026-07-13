package sns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	snsClient   *sns.Client
	Topic       string
	Environment string // dev, stage, prod etc...
}

func NewClient(topic string, env string) Client {
	log.Trace("Setting up sns client")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return Client{
		snsClient:   sns.NewFromConfig(cfg),
		Topic:       topic,
		Environment: env,
	}
}

func (s Client) attributes() map[string]snstypes.MessageAttributeValue {
	return map[string]snstypes.MessageAttributeValue{
		"env": {
			DataType:    aws.String("String"),
			StringValue: aws.String(s.Environment),
		},
	}
}

func (s Client) Publish(subject string, message string) error {
	log.Debugf("Publishing %s to %s", subject, s.Topic)

	_, err := s.snsClient.Publish(context.TODO(), &sns.PublishInput{
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
