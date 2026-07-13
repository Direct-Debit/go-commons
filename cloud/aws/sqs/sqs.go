package sqs

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	sqsClient   *sqs.Client
	Queues      map[string]*string
	Environment string // dev, stage, prod etc...
}

type Attributes map[string]sqstypes.MessageAttributeValue

func NewClient(env string) Client {
	log.Trace("Setting up sqs client")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return Client{
		sqsClient:   sqs.NewFromConfig(cfg),
		Queues:      make(map[string]*string),
		Environment: env,
	}
}

func (c Client) getQueueURL(queue string) (*string, error) {
	if queueUrl, ok := c.Queues[queue]; ok {
		return queueUrl, nil
	} else {
		result, err := c.sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
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

	attr["Env"] = sqstypes.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(c.Environment),
	}

	_, err = c.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageAttributes: attr,
		MessageBody:       aws.String(message),
		DelaySeconds:      int32(delay),
		QueueUrl:          queueUrl,
	})
	return errors.Wrapf(err, "failed to send message to SQS queue %v", queue)
}

func (c Client) DeleteMessage(queue string, receiptHandle string) error {
	queueUrl, err := c.getQueueURL(queue)
	if err != nil {
		return errors.Wrapf(err, "failed to get SQS queue url for %v", queue)
	}

	_, err = c.sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		ReceiptHandle: &receiptHandle,
		QueueUrl:      queueUrl,
	})
	return errors.Wrapf(err, "failed to delete message from %v", queue)
}

func (c Client) Listen(queue string, waitTime int, msgs chan sqstypes.Message) error {
	queueUrl, err := c.getQueueURL(queue)
	if err != nil {
		return errors.Wrapf(err, "failed to get SQS queue url for %v", queue)
	}

	log.Infof("Listening for messages on queue %v...", queue)

	if waitTime <= 0 || waitTime > 20 {
		return fmt.Errorf("waitTime must be between 1 and 20 seconds")
	}

	for {
		output, err := c.sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			MaxNumberOfMessages: int32(10),
			QueueUrl:            queueUrl,
			WaitTimeSeconds:     int32(waitTime),
		})

		if err != nil {
			err = errors.Wrap(err, "failed to receive sqs messages")
			if strings.Contains(err.Error(), "connection reset by peer") {
				log.Warnf("Connection reset by peer, retrying: %v", err)
			} else {
				return err
			}
		}

		if output != nil {
			log.Tracef("Received %d messages from SQS queue %v", len(output.Messages), queue)
			for _, m := range output.Messages {
				msgs <- m
			}
		}
	}
}
