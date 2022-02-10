package amqp

import (
	"fmt"
	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
)

// ConsumerFunc is used to apply logic to the incoming message body.
//
// If ConsumerFunc returns a non-empty byte string,
// those bytes will be published to the message's ReplyTo queue if it has one.
// If the returned error is not nil,
// the message will be n'acked and requeued without any further error logging or handling.
type ConsumerFunc func(body []byte) (response []byte, err error)

type QueueConsumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	queue         amqp.Queue
	closeListener chan *amqp.Error
}

func (qc *QueueConsumer) Connect(url string, queueName string, prefetch int) error {
	var err error

	qc.conn, err = amqp.Dial(url)
	if err != nil {
		return errors.Wrapf(err, "could not dial AMQP on %s", url)
	}
	log.Infof("Dialed to amqp on %s", url)

	qc.channel, err = qc.conn.Channel()
	if err != nil {
		return errors.Wrapf(err, "could not open a channel to AMQP on %s", url)
	}
	log.Info("Opened channel")

	qc.queue, err = qc.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrapf(err, "could not declare amqp queue %s on %s", queueName, url)
	}
	log.Infof("Using queue: %s", qc.queue.Name)

	err = qc.channel.Qos(prefetch, 0, false)
	if err != nil {
		return errors.Wrapf(err, "could not set channel QoS")
	}

	qc.closeListener = qc.channel.NotifyClose(make(chan *amqp.Error))
	return nil
}

func (qc *QueueConsumer) Disconnect() {
	err := qc.channel.Close()
	errlib.WarnError(err, "Failed to close amqp channel")
	err = qc.conn.Close()
	errlib.WarnError(err, "Failed to close amqp connection")
}

func (qc *QueueConsumer) Consume(consume ConsumerFunc) error {
	msgs, err := qc.channel.Consume(
		qc.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrapf(err, "could not create consumer for %s", qc.queue.Name)
	}
	log.Info("Waiting for messages")

	var wg sync.WaitGroup
	for msg := range msgs {
		log.Tracef("Message received: %s", msg.Body)

		wg.Add(1)
		d := msg
		go func() {
			log.Debugf("Starting goroutine for consumer on %s", qc.queue.Name)
			defer wg.Done()

			res, err := consume(d.Body)
			if err != nil {
				err = d.Nack(false, true)
				errlib.WarnError(err, "Failed to n'ack and requeue message on %s", qc.queue.Name)
				return
			}

			if len(res) > 0 && len(d.ReplyTo) > 0 {
				err = qc.channel.Publish(
					"",
					d.ReplyTo,
					false,
					false,
					amqp.Publishing{
						CorrelationId: d.CorrelationId,
						Body:          res,
					},
				)
				if err != nil {
					errlib.WarnError(err, "Failed publish reply from %s to %s", qc.queue.Name, d.ReplyTo)
					err = d.Nack(false, true)
					errlib.WarnError(err, "Failed to n'ack and requeue message on %s", qc.queue.Name)
					return
				}
			}

			err = d.Ack(false)
		}()
	}
	log.Info("Message queue closed, waiting for last message to process")
	wg.Wait()

	select {
	case e := <-qc.closeListener:
		return e
	default:
		return fmt.Errorf("queue %s closed unexpectedly", qc.queue.Name)
	}
}
