package pubsub

import (
	"context"
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publish[T any](
	ch *amqp.Channel,
	exchange,
	key string,
	val T,
	marshaller func(T) ([]byte, string, error),
) error {
	data, contentType, err := marshaller(val)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		context.Background(),
		exchange, key,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        data,
		},
	)
}

func PublishGameLog(ch *amqp.Channel, log routing.GameLog) error {
	return PublishGob(
		ch,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%s.%s", routing.GameLogSlug, log.Username),
		log,
	)
}
