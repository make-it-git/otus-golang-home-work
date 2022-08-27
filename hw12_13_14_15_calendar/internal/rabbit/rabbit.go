package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	_ "github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
}

func New(cnf config.RabbitmqConnectionConf) (*Rabbit, error) {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%d/", cnf.User, cnf.Password, cnf.Host, cnf.Port)
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	err = ch.ExchangeDeclare(
		cnf.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Rabbit{
		connection: conn,
		channel:    ch,
		exchange:   cnf.Exchange,
	}, nil
}

func (r *Rabbit) Stop(ctx context.Context) {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.connection != nil {
		_ = r.connection.Close()
	}
}

func (r *Rabbit) Publish(message json.RawMessage) error {
	return r.channel.PublishWithContext(
		context.Background(),
		r.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			Body: message,
		})
}
