package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	sqsClient   *sqs.SQS
	Queues      map[string]*string
	Environment string // dev, stage, prod etc...
}

type Attributes map[string]*sqs.MessageAttributeValue

func NewClient(env string) Client {
	log.Trace("Setting up sqs client")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return Client{
		sqsClient:   sqs.New(sess),
		Queues:      make(map[string]*string),
		Environment: env,
	}
}

func (c Client) getQueueURL(queue string) (*string, error) {
	if queueUrl, ok := c.Queues[queue]; ok {
		return queueUrl, nil
	} else {
		result, err := c.sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
			QueueName: &queue,
		})
		if err != nil {
			return nil, err
		}
		c.Queues[queue] = result.QueueUrl
	}
	return c.Queues[queue], nil
}

// Attr accepts a Nil map if no additional attributes should be set
func (c Client) SendMessage(queue string, message string, attr Attributes) error {
	queueUrl, err := c.getQueueURL(queue)
	if err != nil {
		return errors.Wrapf(err, "failed to get SQS queue url for %v", queue)
	}

	if attr == nil {
		attr = make(Attributes)
	}

	attr["Env"] = &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(c.Environment),
	}

	_, err = c.sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: attr,
		MessageBody:       aws.String(message),
		QueueUrl:          queueUrl,
	})
	return errors.Wrapf(err, "failed to send message to SQS queue %v", queue)
}

func (c Client) DeleteMessage(queue string, receiptHandle string) error {
	queueUrl, err := c.getQueueURL(queue)
	if err != nil {
		return errors.Wrapf(err, "failed to get SQS queue url for %v", queue)
	}

	_, err = c.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		ReceiptHandle: &receiptHandle,
		QueueUrl:      queueUrl,
	})
	return errors.Wrapf(err, "failed to delete message from %v", queue)
}
