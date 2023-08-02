package sqs

import (
	"fmt"

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
func (c Client) SendMessage(queue string, message string, delay int, attr Attributes) error {
	if delay > 900 {
		delay = 900
	}

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
		DelaySeconds:      aws.Int64(int64(delay)),
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

func (c Client) Listen(queue string, waitTime int, msgs chan *sqs.Message) error {
	queueUrl, err := c.getQueueURL(queue)
	if err != nil {
		return errors.Wrapf(err, "failed to get SQS queue url for %v", queue)
	}

	if waitTime <= 0 || waitTime > 20 {
		return fmt.Errorf("waitTime must be between 1 and 20 seconds")
	}

	for {
		output, err := c.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			MaxNumberOfMessages: aws.Int64(int64(10)),
			QueueUrl:            queueUrl,
			WaitTimeSeconds:     aws.Int64(int64(waitTime)),
		})
		if err != nil {
			return errors.Wrap(err, "failed to receive sqs messages")
		}

		for _, m := range output.Messages {
			msgs <- m
		}
	}
}
