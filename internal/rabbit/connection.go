package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type Rabbit struct {
	Channel *amqp.Channel
	QueueName string
	URL string
	Concurrency int
	ReconAttempt int
	ReconInterval time.Duration
}

type Opts struct {
	QueueName string
	URL string
	Concurrency int
	ReconAttempt int
	ReconInterval time.Duration
}

func NewConn(opts Opts) (*Rabbit, error) {
	conn, err := amqp.Dial(opts.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()

	return &Rabbit{
		Channel: ch,
		QueueName: opts.QueueName,
		URL: opts.URL,
		Concurrency: opts.Concurrency,
		ReconAttempt: opts.ReconAttempt,
		ReconInterval: opts.ReconInterval,
	}, err
}
