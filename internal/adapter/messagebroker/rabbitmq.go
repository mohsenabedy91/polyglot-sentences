package messagebroker

import (
	"encoding/json"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	log  logger.Logger
}

func NewRabbitMQ(url string, log logger.Logger) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		log:  log,
	}, nil
}

func (r *RabbitMQ) Close() {
	if err := r.conn.Close(); err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQ, err.Error(), nil)
	}
	return
}

func (r *RabbitMQ) Produce(name string, msg interface{}, delaySeconds int64) error {
	message, err := json.Marshal(msg)
	if err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error marshalling value: %v", err), nil)
		return err
	}

	extra := map[logger.ExtraKey]interface{}{
		logger.QueueName: name,
		logger.Body:      message,
	}

	channel, err := r.conn.Channel()
	if err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error create channel: %v", err), extra)
		return err
	}

	defer func(ch *amqp.Channel) {
		if err = ch.Close(); err != nil {
			r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error closing channel: %v", err), extra)
		}
	}(channel)

	if err = channel.ExchangeDeclare(
		"delayed_exchange",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{"x-delayed-type": "direct"},
	); err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error ExchangeDeclare: %v", err), extra)
		return err
	}

	queueDeclare, err := channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error QueueDeclare: %v", err), extra)
		return err
	}

	if err = channel.QueueBind(
		queueDeclare.Name,
		queueDeclare.Name,
		"delayed_exchange",
		false,
		nil,
	); err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error QueueBind: %v", err), extra)
		return err
	}

	if err = channel.Publish(
		"delayed_exchange",
		queueDeclare.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
			Headers:     amqp.Table{"x-delay": delaySeconds * 1000},
		},
	); err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error Publish message: %v", err), extra)
		return err
	}

	return nil
}

func (r *RabbitMQ) RegisterConsumer(name string, callback func(message []byte) error) error {

	extra := map[logger.ExtraKey]interface{}{
		logger.QueueName: name,
	}

	channel, err := r.conn.Channel()
	if err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error create channel: %v", err), extra)
		return err
	}

	queueDeclare, err := channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.log.Error(
			logger.Queue,
			logger.RabbitMQRegisterConsumer,
			fmt.Sprintf("Error QueueDeclare channel: %v", err),
			extra,
		)
		return err
	}

	deliveries, err := channel.Consume(
		queueDeclare.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.log.Error(
			logger.Queue,
			logger.RabbitMQRegisterConsumer,
			fmt.Sprintf("Error Consume channel: %v", err),
			extra,
		)
		return err
	}

	go func() {
		for delivery := range deliveries {
			extra = map[logger.ExtraKey]interface{}{
				logger.QueueName: name,
				logger.Body:      string(delivery.Body),
			}
			if err = callback(delivery.Body); err != nil {
				r.log.Error(
					logger.Queue,
					logger.RabbitMQRegisterConsumer,
					fmt.Sprintf("Error Consume message: %v", err),
					extra,
				)
			} else {
				if err = delivery.Ack(false); err != nil {
					r.log.Error(
						logger.Queue,
						logger.RabbitMQRegisterConsumer,
						fmt.Sprintf("Error Ack Consume message: %v", err),
						extra,
					)
				}
			}
		}
	}()

	return nil
}
