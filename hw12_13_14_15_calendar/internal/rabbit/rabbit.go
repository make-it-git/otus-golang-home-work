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
	queue      string
	done       chan struct{}
}

func NewProducer(cnf config.RabbitmqConnectionConf) (*Rabbit, error) {
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
		done:       make(chan struct{}),
	}, nil
}

func NewConsumer(cnf config.RabbitmqConsumerConf) (*Rabbit, error) {
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

	_, err = ch.QueueDeclare(
		cnf.Queue,
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

	err = ch.QueueBind(
		cnf.Queue,
		"",
		cnf.Exchange,
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
		queue:      cnf.Queue,
		done:       make(chan struct{}),
	}, nil
}

func (r *Rabbit) Stop(ctx context.Context) {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.connection != nil {
		_ = r.connection.Close()
	}
	close(r.done)
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

func (r *Rabbit) Consume() (<-chan []byte, error) {
	msgs, err := r.channel.Consume(
		r.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte)

	go func() {
		c := true
		for c {
			select {
			case m := <-msgs:
				ch <- m.Body
			case <-r.done:
				c = false
				break
			}
		}
		close(ch)
	}()

	return ch, nil
}
