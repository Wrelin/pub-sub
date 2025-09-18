package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	isDurableQueue := queueType == SimpleQueueDurable
	isTransientQueue := queueType == SimpleQueueTransient

	queue, err := ch.QueueDeclare(queueName, isDurableQueue, isTransientQueue, isTransientQueue, false, nil)
	if err != nil {
		return ch, queue, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)

	return ch, queue, err
}
