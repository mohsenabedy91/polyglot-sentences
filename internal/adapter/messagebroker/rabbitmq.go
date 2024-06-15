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

func (r *Queue) SetupRabbitMQ(url string, log logger.Logger) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	r.Driver = &RabbitMQ{
		conn: conn,
		log:  log,
	}

	return nil
}

func (r *RabbitMQ) Close() {
	err := r.conn.Close()
	if err != nil {
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

	ch, err := r.conn.Channel()
	if err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error create channel: %v", err), extra)
		return err
	}

	defer func(ch *amqp.Channel) {
		err = ch.Close()
		if err != nil {
			r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error closing channel: %v", err), extra)
		}
	}(ch)

	err = ch.ExchangeDeclare(
		"delayed_exchange",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{"x-delayed-type": "direct"},
	)
	if err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error ExchangeDeclare: %v", err), extra)
		return err
	}

	q, err := ch.QueueDeclare(
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

	err = ch.QueueBind(
		q.Name,
		q.Name,
		"delayed_exchange",
		false,
		nil,
	)
	if err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error QueueBind: %v", err), extra)
		return err
	}

	headers := amqp.Table{"x-delay": delaySeconds * 1000}
	err = ch.Publish(
		"delayed_exchange",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
			Headers:     headers,
		},
	)
	if err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error Publish message: %v", err), extra)
		return err
	}

	return nil
}

func (r *RabbitMQ) RegisterConsumer(name string, callback func(message []byte) error) error {

	extra := map[logger.ExtraKey]interface{}{
		logger.QueueName: name,
	}

	ch, err := r.conn.Channel()
	if err != nil {
		r.log.Error(logger.Queue, logger.RabbitMQProduce, fmt.Sprintf("Error create channel: %v", err), extra)
		return err
	}

	q, err := ch.QueueDeclare(
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

	deliveries, err := ch.Consume(
		q.Name,
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
