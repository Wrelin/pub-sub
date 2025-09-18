package pubsub

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	err = ch.Qos(10, 0, false)
	if err != nil {
		return err
	}

	messages, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for message := range messages {
			obj, err := unmarshaller(message.Body)
			if err != nil {
				fmt.Printf("Error while try numarshall body: %+v\n", err)
				continue
			}

			switch handler(obj) {
			case Ack:
				message.Ack(false)
			case NackRequeue:
				message.Nack(false, true)
			case NackDiscard:
				fallthrough
			default:
				message.Nack(false, false)
			}
		}
	}()

	return nil
}
