package pubsub

import (
	"bytes"
	"encoding/gob"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	marshaller := func(obj T) ([]byte, string, error) {
		var b bytes.Buffer
		enc := gob.NewEncoder(&b)
		err := enc.Encode(val)
		if err != nil {
			return []byte{}, "", err
		}
		return b.Bytes(), "application/gob", nil
	}

	return publish(ch, exchange, key, val, marshaller)
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	unmarshaller := func(data []byte) (T, error) {
		b := bytes.NewBuffer(data)
		var obj T

		dec := gob.NewDecoder(b)
		err := dec.Decode(&obj)
		return obj, err
	}

	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		unmarshaller,
	)
}
